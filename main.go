package main

import (
	"learning/utility"
	"log"
)

func main() {
	ciphered := utility.CaesarCipher("testing-hello", 2)
	log.Println(ciphered)
}
