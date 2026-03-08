package app

import (
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/http"
	"super-simple-queues/internal/server/tcp"
)

type App struct {
	queueManager *queue.Manager
}

func New() *App {
	return &App{queue.NewManager()}
}

func (a *App) Run(tcpPort int, httpPort int) error {
	errChan := make(chan error)

	server.RunGo(tcp.NewServer(a.queueManager), tcpPort, errChan)
	server.RunGo(http.NewServer(a.queueManager), httpPort, errChan)

	return <-errChan
}
