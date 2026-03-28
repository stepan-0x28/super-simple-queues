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
	expectedCfg := Config{
		LoggingLevel:      slog.LevelDebug,
		TCPPort:           4444,
		HTTPPort:          4040,
		QueueChunkSize:    512,
		TCPConnBufferSize: 128,
	}

	t.Setenv("LOGGING_LEVEL", "Debug")
	t.Setenv("TCP_PORT", strconv.Itoa(expectedCfg.TCPPort))
	t.Setenv("HTTP_PORT", strconv.Itoa(expectedCfg.HTTPPort))
	t.Setenv("QUEUE_CHUNK_SIZE", strconv.Itoa(expectedCfg.QueueChunkSize))
	t.Setenv("TCP_CONN_BUFFER_SIZE", strconv.Itoa(expectedCfg.TCPConnBufferSize))

	cfg := loadConfig(t)

	checkConfigs(t, cfg, expectedCfg)
}

func TestLoadConfig_IncorrectValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"incorrect logging level", "LOGGING_LEVEL", "Emergency"},
		{"incorrect TCP port", "TCP_PORT", "Open"},
		{"incorrect HTTP port", "HTTP_PORT", "Open"},
		{"incorrect queue chunk size", "QUEUE_CHUNK_SIZE", "Medium"},
		{"incorrect TCP connection buffer size", "TCP_CONN_BUFFER_SIZE", "Medium"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(test.key, test.value)

			_, err := LoadConfig()

			if err == nil {
				t.Fatalf("config loading error expected")
			}
		})
	}
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
