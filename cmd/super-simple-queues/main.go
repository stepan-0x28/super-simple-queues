package main

import (
	"log"
	"super-simple-queues/internal/app"
)

func main() {
	newApp := app.NewApp()

	if err := newApp.Run(); err != nil {
		log.Fatal(err)
	}
}
