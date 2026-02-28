package tcp

import (
	"fmt"
	"log"
	"net"
	"super-simple-queues/internal/queue"
)

type PortListener struct {
	queueManager *queue.Manager
}

func NewPortListener(queueManager *queue.Manager) *PortListener {
	return &PortListener{queueManager: queueManager}
}

func (p *PortListener) Run(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err)

			continue
		}

		connection := NewConnection(conn)

		go func(c *Connection) {
			if err := c.Run(p.queueManager); err != nil {
				log.Println(err)
			}
		}(connection)
	}
}
