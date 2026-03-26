package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	LoggingLevel      slog.Level
	TCPPort           int
	HTTPPort          int
	QueueChunkSize    int
	TCPConnBufferSize int
}

func LoadConfig() (Config, error) {
	var loggingLevel slog.Level

	switch getEnv("LOGGING_LEVEL", "") {
	case "Debug":
		loggingLevel = slog.LevelDebug
	case "Warn":
		loggingLevel = slog.LevelWarn
	case "Error":
		loggingLevel = slog.LevelError
	default:
		loggingLevel = slog.LevelInfo
	}

	tcpPort, err := getEnvInt("TCP_PORT", 8888)

	if err != nil {
		return Config{}, err
	}

	httpPort, err := getEnvInt("HTTP_PORT", 8080)

	if err != nil {
		return Config{}, err
	}

	queueChunkSize, err := getEnvInt("QUEUE_CHUNK_SIZE", 1024)

	if err != nil {
		return Config{}, err
	}

	tcpConnBufferSize, err := getEnvInt("TCP_CONN_BUFFER_SIZE", 256)

	if err != nil {
		return Config{}, err
	}

	return Config{
		LoggingLevel:      loggingLevel,
		TCPPort:           tcpPort,
		HTTPPort:          httpPort,
		QueueChunkSize:    queueChunkSize,
		TCPConnBufferSize: tcpConnBufferSize,
	}, nil
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
