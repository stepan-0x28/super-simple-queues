package main

import (
	"log/slog"
	"os"
	"super-simple-queues/internal/app"
	"super-simple-queues/internal/config"
)

func main() {
	var loggingLevel slog.LevelVar

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     &loggingLevel,
	})

	slog.SetDefault(slog.New(handler))

	cfg, err := config.LoadConfig()

	if err != nil {
		slog.Error("configuration loading failed", slog.Any("err", err))

		os.Exit(1)
	}

	loggingLevel.Set(cfg.LoggingLevel)

	a := app.New()

	if err = a.Run(cfg); err != nil {
		slog.Error("application error", slog.Any("err", err))

		os.Exit(1)
	}
}
