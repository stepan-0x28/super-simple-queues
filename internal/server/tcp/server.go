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
	queueManager   *queue.Manager
	connBufferSize int
	listener       net.Listener
}

func NewServer(queueManager *queue.Manager, connBufferSize int) *Server {
	return &Server{queueManager: queueManager, connBufferSize: connBufferSize}
}

func (s *Server) Run(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	s.listener = l

	for {
		conn, err := l.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}

			slog.Warn("failed to accept connection", slog.Any("err", err))

			continue
		}

		go s.processConnection(conn)
	}
}

func (s *Server) processConnection(conn net.Conn) {
	addrAttr := slog.Any("addr", conn.RemoteAddr())

	defer func() {
		if err := conn.Close(); err != nil {
			slog.Warn("failed to close connection", slog.Any("err", err), addrAttr)
		} else {
			slog.Info("connection closed", addrAttr)
		}
	}()

	c := newConnection(conn, s.connBufferSize)

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
