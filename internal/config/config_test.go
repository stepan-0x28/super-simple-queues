package config

import (
	"log/slog"
	"reflect"
	"strconv"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg := loadConfig(t)

	defaultCfg := Config{
		LoggingLevel:      slog.LevelInfo,
		TCPPort:           8888,
		HTTPPort:          8080,
		QueueChunkSize:    1024,
		TCPConnBufferSize: 256,
	}

	checkConfigs(t, cfg, defaultCfg)
}

func TestLoadConfig_SetValues(t *testing.T) {
	const (
		tcpPort           = 4444
		httpPort          = 4040
		queueChunkSize    = 512
		tcpConnBufferSize = 128
	)

	t.Setenv("LOGGING_LEVEL", "Debug")
	t.Setenv("TCP_PORT", strconv.Itoa(tcpPort))
	t.Setenv("HTTP_PORT", strconv.Itoa(httpPort))
	t.Setenv("QUEUE_CHUNK_SIZE", strconv.Itoa(queueChunkSize))
	t.Setenv("TCP_CONN_BUFFER_SIZE", strconv.Itoa(tcpConnBufferSize))

	cfg := loadConfig(t)

	expectedCfg := Config{
		LoggingLevel:      slog.LevelDebug,
		TCPPort:           tcpPort,
		HTTPPort:          httpPort,
		QueueChunkSize:    queueChunkSize,
		TCPConnBufferSize: tcpConnBufferSize,
	}

	checkConfigs(t, cfg, expectedCfg)
}

func TestLoadConfig_incorrectLoggingLevel(t *testing.T) {
	checkIncorrectValue(t, "LOGGING_LEVEL", "Emergency")
}

func TestLoadConfig_incorrectTCPPort(t *testing.T) {
	checkIncorrectValue(t, "TCP_PORT", "Open")
}

func TestLoadConfig_incorrectHTTPPort(t *testing.T) {
	checkIncorrectValue(t, "HTTP_PORT", "Open")
}

func TestLoadConfig_incorrectQueueChunkSize(t *testing.T) {
	checkIncorrectValue(t, "QUEUE_CHUNK_SIZE", "Medium")
}

func TestLoadConfig_incorrectTCPConnBufferSize(t *testing.T) {
	checkIncorrectValue(t, "TCP_CONN_BUFFER_SIZE", "Medium")
}

func loadConfig(t *testing.T) Config {
	t.Helper()

	cfg, err := LoadConfig()

	if err != nil {
		t.Fatalf("config loading error, %v", err)
	}

	return cfg
}

func checkConfigs(t *testing.T, cfg Config, expectedCfg Config) {
	t.Helper()

	if !reflect.DeepEqual(cfg, expectedCfg) {
		t.Fatalf("expected config %v, received config %v", expectedCfg, cfg)
	}
}

func checkIncorrectValue(t *testing.T, key string, value string) {
	t.Helper()

	t.Setenv(key, value)

	_, err := LoadConfig()

	if err == nil {
		t.Fatalf("config loading error expected")
	}
}
