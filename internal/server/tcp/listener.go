package tcp

import (
	"fmt"
	"log"
	"net"
	"super-simple-queues/internal/queue"
)

type Listener struct {
	queueManager *queue.Manager
}

func NewListener(queueManager *queue.Manager) *Listener {
	return &Listener{queueManager: queueManager}
}

func (l *Listener) Run(port int) error {
	tl, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		return err
	}

	for {
		conn, err := tl.Accept()

		if err != nil {
			log.Println(err)

			continue
		}

		go func(c *connection) {
			if err := c.run(l.queueManager); err != nil {
				log.Println(err)
			}
		}(newConnection(conn))
	}
}
