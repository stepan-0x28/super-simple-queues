package app

import (
	"super-simple-queues/config"
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/http"
	"super-simple-queues/internal/server/tcp"
)

type App struct {
	config       *config.Config
	queueManager *queue.Manager
}

func NewApp(config *config.Config) *App {
	return &App{config, queue.NewManager()}
}

func (a *App) Run() error {
	errChan := make(chan error)

	server.RunGo(tcp.NewPortListener(a.queueManager), a.config.TCPPort, errChan)
	server.RunGo(http.NewHttp(a.queueManager), a.config.HTTPPort, errChan)

	return <-errChan
}
