package main

import (
	"github.com/sornick01/http_chat/server"
	"log"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Run("1234")
}
