package hacker

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	hn "ytsruh.com/basics/hacker/client"
)

func RunCache() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./hacker/index.html"))
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/", handlerCache(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=")
}

func RunCacheWithMutex() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./hacker/index.html"))
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/", handlerCacheWithMutex(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handlerCache(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getCachedStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func handlerCacheWithMutex(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getCachedStoriesWithMutex(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

var (
	cache           []item
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
)

func getCachedStories(numStories int) ([]item, error) {
	if time.Since(cacheExpiration) < 0 {
		log.Print("Using cache")
		return cache, nil
	}
	log.Print("Not using cache")
	stories, err := getTopStoriesCache(numStories)
	if err != nil {
		return nil, err
	}
	cache = stories
	cacheExpiration = time.Now().Add(100 * time.Second)
	return cache, nil
}

func getCachedStoriesWithMutex(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Since(cacheExpiration) < 0 {
		log.Print("Using cache")
		return cache, nil
	}
	log.Print("Not using cache")
	stories, err := getTopStoriesCache(numStories)
	if err != nil {
		return nil, err
	}
	cache = stories
	cacheExpiration = time.Now().Add(100 * time.Second)
	return cache, nil
}

func getTopStoriesCache(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStoriesCache(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

func getStoriesCache(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			var client hn.Client // This line was previously outside the go routine but caused a race condition
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}
