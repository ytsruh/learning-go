package hacker

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	hn "learning/hacker/client"
)

func RunFirstAttempt() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./hacker/index.html"))

	http.HandleFunc("/", handlerFirst(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handlerFirst(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getTopStoriesFirst(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getTopStoriesFirst(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	var stories []item
	for _, id := range ids {
		type result struct {
			item item
			err  error
		}
		resultCh := make(chan result)
		go func(id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{err: err}
			}
			resultCh <- result{item: parseHNItem(hnItem)}
		}(id)

		res := <-resultCh
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
			if len(stories) >= numStories {
				break
			}
		}
	}
	return stories, nil
}
