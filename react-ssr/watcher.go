package main

import (
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

var Files map[string]bool

// Simple implementation to watch a directory for changes
func RunWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	Files = make(map[string]bool)
	dir := "./" // directory to watch

	// List all files in directory
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileInfo {
		Files[file.Name()] = true
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					Files[event.Name] = true // Fix: Remove the parentheses after event.Name
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file:", event.Name)
					delete(Files, event.Name)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					Files[event.Name] = true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan bool) // block until a signal is received
}
