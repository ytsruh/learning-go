package main

import (
	"net/http"

	"learning/test-with-go/app"
)

func main() {
	http.ListenAndServe(":3000", &app.Server{})
}
