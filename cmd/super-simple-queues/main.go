package main

import (
	"log"
	"os"
	"strconv"
	"super-simple-queues/internal/app"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	tcpPort, err := strconv.Atoi(getEnv("TCP_PORT", "8888"))

	if err != nil {
		log.Fatal(err)
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))

	if err != nil {
		log.Fatal(err)
	}

	queueChunkSize, err := strconv.Atoi(getEnv("QUEUE_CHUNK_SIZE", "1024"))

	if err != nil {
		log.Fatal(err)
	}

	a := app.New()

	if err = a.Run(tcpPort, httpPort, queueChunkSize); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
