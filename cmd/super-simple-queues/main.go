package main

import (
	"log/slog"
	"os"
	"strconv"
	"super-simple-queues/internal/app"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     convertLoggingLevel(getEnv("LOGGING_LEVEL", "Info")),
	})

	slog.SetDefault(slog.New(handler))

	tcpPort, err := getEnvInt("TCP_PORT", 8888)

	if err != nil {
		slog.Error("invalid TCP port", slog.Any("err", err))

		os.Exit(1)
	}

	httpPort, err := getEnvInt("HTTP_PORT", 8080)

	if err != nil {
		slog.Error("invalid HTTP port", slog.Any("err", err))

		os.Exit(1)
	}

	queueChunkSize, err := getEnvInt("QUEUE_CHUNK_SIZE", 1024)

	if err != nil {
		slog.Error("invalid queue chunk size", slog.Any("err", err))

		os.Exit(1)
	}

	a := app.New()

	if err = a.Run(tcpPort, httpPort, queueChunkSize); err != nil {
		slog.Error("application error", slog.Any("err", err))

		os.Exit(1)
	}
}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	return strconv.Atoi(getEnv(key, strconv.Itoa(defaultValue)))
}

func convertLoggingLevel(l string) slog.Level {
	switch l {
	case "Debug":
		return slog.LevelDebug
	case "Warn":
		return slog.LevelWarn
	case "Error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
