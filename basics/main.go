package main

import (
	"fmt"

	utility "ytsruh.com/basics/utility"
)

func main() {
	cipher := utility.CaesarCipher("hello world", 3)
	fmt.Println(cipher)
}
