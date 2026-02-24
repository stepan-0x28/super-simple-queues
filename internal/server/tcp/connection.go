package tcp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"super-simple-queues/internal/queue"
)

type Connection struct {
	conn     net.Conn
	status   string
	queueKey string
	manager  *queue.Manager
}

func NewConnection(conn net.Conn, manager *queue.Manager) *Connection {
	connection := &Connection{conn, "new", "", manager}

	return connection
}

func (c *Connection) Run() {
	for {
		message := c.read()

		if c.status == "new" {
			if message["mode"] == "sending" {
				c.status = "receiving"
				c.queueKey = fmt.Sprintf("%v", message["queue_key"])
			}
		} else if c.status == "receiving" {
			workingQueue := c.manager.GetQueue(c.queueKey)

			if workingQueue != nil {
				workingQueue.AddElement(message)
			}
		}
	}
}

// TODO change the return value
func (c *Connection) read() map[string]interface{} {
	b := make([]byte, 4)

	_, err := c.conn.Read(b)

	if err != nil {
		return nil
	}

	length := binary.BigEndian.Uint32(b)

	b = make([]byte, length)

	_, err = c.conn.Read(b)

	if err != nil {
		return nil
	}

	var message map[string]interface{}

	err = json.Unmarshal(b, &message)

	if err != nil {
		return nil
	}

	return message
}
