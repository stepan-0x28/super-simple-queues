package main

import (
	"log"
	"os"
	"strconv"
	"super-simple-queues/internal/app"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	tcpPort, err := strconv.Atoi(os.Getenv("TCP_PORT"))

	if err != nil {
		log.Fatal(err)
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))

	if err != nil {
		log.Fatal(err)
	}

	a := app.New()

	if err = a.Run(tcpPort, httpPort); err != nil {
		log.Fatal(err)
	}
}
