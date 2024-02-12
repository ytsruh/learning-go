package main

import (
	"ytsruh.com/goly/model"
	"ytsruh.com/goly/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	model.Setup()
	server.SetupAndListen()
}
