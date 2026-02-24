package tcp

import (
	"fmt"
	"net"
	"super-simple-queues/internal/queue"
)

type PortListener struct {
	manager *queue.Manager
}

func NewTcp(manager *queue.Manager) *PortListener {
	tcp := &PortListener{manager}

	return tcp
}

func (p *PortListener) Run(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			return err
		}

		connection := NewConnection(conn, p.manager)

		go connection.Run()
	}
}
