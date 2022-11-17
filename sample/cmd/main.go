package main

import (
	"log"
	"sample/server"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Run("1234")
}
