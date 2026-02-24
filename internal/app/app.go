package app

import (
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/tcp"
	"super-simple-queues/internal/server/web"
)

type App struct {
	queueManager *queue.Manager
}

func NewApp() *App {
	return &App{queue.NewManager()}
}

func (a *App) Run() error {
	errChan := make(chan error)

	server.RunGo(tcp.NewTcp(a.queueManager), 8888, errChan)
	server.RunGo(web.NewWeb(a.queueManager), 8080, errChan)

	return <-errChan
}
