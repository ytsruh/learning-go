package main

import (
	"learning/normalise"
	"log"
)

func main() {
	result := normalise.PhoneRegex("123 456 7891")
	log.Println(result)
}
