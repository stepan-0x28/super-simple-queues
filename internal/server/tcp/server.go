package tcp

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"super-simple-queues/internal/queue"
)

type Server struct {
	queueManager *queue.Manager
}

func NewServer(queueManager *queue.Manager) *Server {
	return &Server{queueManager: queueManager}
}

func (s *Server) Run(port int) error {
	tl, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	for {
		conn, err := tl.Accept()

		if err != nil {
			slog.Warn("failed to accept connection", slog.Any("err", err))

			continue
		}

		go s.processConnection(conn)
	}
}

func (s *Server) processConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Warn("failed to close connection", slog.Any("err", err))
		}
	}()

	c := newConnection(conn)

	operatingMode, q, err := c.init(s.queueManager)

	if err != nil {
		slog.Warn("failed to initialize connection", slog.Any("err", err))

		return
	}

	if err = c.run(operatingMode, q); err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Warn("failed to process connection", slog.Any("err", err))
		}
	}
}
