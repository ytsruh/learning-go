package lessons

import (
	"fmt"
	"io/ioutil"
	"log"
)

func Files() {
	fileBytes, err := ioutil.ReadFile("example01.txt")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fileBytes)
	fileString := string(fileBytes)
	fmt.Println(fileString)
}

func MakeFile() {
	s := "Hello from Go World"
	err := ioutil.WriteFile("example02.txt", []byte(s), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
