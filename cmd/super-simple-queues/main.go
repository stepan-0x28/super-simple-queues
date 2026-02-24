package main

import (
	"log"
	"super-simple-queues/config"
	"super-simple-queues/internal/app"
)

func main() {
	appConfig, err := config.LoadConfig("config/config.ini")

	if err != nil {
		return
	}

	newApp := app.NewApp(appConfig)

	if err = newApp.Run(); err != nil {
		log.Fatal(err)
	}
}
