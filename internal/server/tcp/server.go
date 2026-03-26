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
	bufferSize   int
}

func NewServer(queueManager *queue.Manager, bufferSize int) *Server {
	return &Server{queueManager: queueManager, bufferSize: bufferSize}
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
	addrAttr := slog.Any("addr", conn.RemoteAddr())

	defer func() {
		if err := conn.Close(); err == nil {
			slog.Info("connection closed", addrAttr)
		} else {
			slog.Warn("failed to close connection", slog.Any("err", err), addrAttr)
		}
	}()

	c := newConnection(conn, s.bufferSize)

	slog.Info("new connection", addrAttr)

	operatingMode, q, err := c.init(s.queueManager)

	if err != nil {
		slog.Warn("failed to initialize connection", slog.Any("err", err), addrAttr)

		return
	}

	if err = c.run(operatingMode, q); err != nil {
		switch {
		case errors.Is(err, io.EOF):
			// normal termination of connection processing
		case errors.Is(err, queue.ErrQueueDeleted):
			slog.Info("connection processing ended due to queue deletion", addrAttr)
		default:
			slog.Warn("failed to process connection", slog.Any("err", err), addrAttr)
		}
	}
}
