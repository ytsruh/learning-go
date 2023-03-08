package main

import (
	"learning/sitemap"
	"log"
)

func main() {
	success, err := sitemap.Generate("https://www.ytsruh.com", 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Success: %t", success)
}
