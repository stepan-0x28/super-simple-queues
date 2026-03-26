package app

import (
	"log/slog"
	"super-simple-queues/internal/config"
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/http"
	"super-simple-queues/internal/server/tcp"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run(cfg config.Config) error {
	m := queue.NewManager(cfg.QueueChunkSize)

	errChan := make(chan error)

	server.RunGo(tcp.NewServer(m, cfg.TCPBufferSize), cfg.TCPPort, errChan)
	server.RunGo(http.NewServer(m), cfg.HTTPPort, errChan)

	slog.Info("application started", slog.Any("config", cfg))

	return <-errChan
}
