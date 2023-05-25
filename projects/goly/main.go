package goly

import (
	"learning/projects/goly/model"
	"learning/projects/goly/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	model.Setup()
	server.SetupAndListen()
}
