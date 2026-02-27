package main

import (
	"log"
	"super-simple-queues/config"
	"super-simple-queues/internal/app"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	appConfig, err := config.LoadConfig("config/config.ini")

	if err != nil {
		log.Fatal(err)
	}

	newApp := app.NewApp(appConfig)

	if err = newApp.Run(); err != nil {
		log.Fatal(err)
	}
}
