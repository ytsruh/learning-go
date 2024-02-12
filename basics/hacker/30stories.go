package hacker

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	hn "ytsruh.com/basics/hacker/client"
)

func Run30Stories() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./hacker/index.html"))

	http.HandleFunc("/", handler30Stories(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler30Stories(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getTopStories30(numStories)
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

func getTopStories30(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4 // Multiply by 5/4 so its equal to 1.25. But using integers means no issues with floats
		stories = append(stories, getStories30(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

func getStories30(ids []int) []item {
	var client hn.Client
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
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
