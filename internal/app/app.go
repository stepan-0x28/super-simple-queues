package app

import (
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server"
	"super-simple-queues/internal/server/tcp"
)

type App struct {
	queueManager *queue.Manager
}

func NewApp() *App {
	return &App{queue.NewManager()}
}

func (a *App) Run() error {
	webServer := server.NewWeb(a.queueManager)
	newTcp := tcp.NewTcp(a.queueManager)

	errChan := make(chan error)

	//TODO remove code duplication
	go func() {
		err := newTcp.Run(8888)

		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		err := webServer.Run(8080)

		if err != nil {
			errChan <- err
		}
	}()

	return <-errChan
}
