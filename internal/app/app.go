package app

import (
	"log/slog"
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/http"
	"super-simple-queues/internal/server/tcp"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run(tcpPort int, httpPort int, queueChunkSize int) error {
	m := queue.NewManager(queueChunkSize)

	errChan := make(chan error)

	server.RunGo(tcp.NewServer(m), tcpPort, errChan)
	server.RunGo(http.NewServer(m), httpPort, errChan)

	slog.Info("application started", slog.Any("tcp_port", tcpPort), slog.Any("http_port", httpPort),
		slog.Any("queue_chunk_size", queueChunkSize))

	return <-errChan
}
