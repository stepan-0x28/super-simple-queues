package tcp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/utils"
)

type Connection struct {
	conn         net.Conn
	lengthArray  [4]byte
	lengthBuffer []byte
	queue        *queue.Queue
}

func NewConnection(conn net.Conn) *Connection {
	c := &Connection{conn: conn}

	c.lengthBuffer = c.lengthArray[:]

	return c
}

func (c *Connection) Run(queueManager *queue.Manager) error {
	initMessage, err := c.readMessage()

	if err != nil {
		return err
	}

	queueKey, err := utils.GetStringValue(initMessage, "queue_key")

	if err != nil {
		return err
	}

	mode, err := utils.GetStringValue(initMessage, "mode")

	if err != nil {
		return err
	}

	workQueue, ok := queueManager.Get(queueKey)

	if !ok {
		return fmt.Errorf("queue with key \"%v\" does not exist", queueKey)
	}

	c.queue = workQueue

	switch mode {
	case "sending":
		err = c.receiveMessages()
	case "receiving":
		err = c.sendMessages()
	default:
		err = errors.New("unknown mode")
	}

	return err
}

func (c *Connection) receiveMessages() error {
	for {
		message, err := c.readMessage()

		if err != nil {
			return err
		}

		c.queue.Add(message)

		err = c.writeMessage(map[string]any{"confirmation": "1"})

		if err != nil {
			return err
		}
	}
}

func (c *Connection) sendMessages() error {
	for {
		queueMessage := c.queue.Take()

		err := c.writeMessage(queueMessage)

		if err != nil {
			c.queue.PutBack(queueMessage)

			return err
		}

		tcpMessage, err := c.readMessage()

		if err != nil {
			c.queue.PutBack(queueMessage)

			return err
		}

		value, err := utils.GetStringValue(tcpMessage, "confirmation")

		if err != nil || value != "1" {
			c.queue.PutBack(queueMessage)

			return errors.New("the \"confirmation\" key is missing or invalid")
		}
	}
}

func (c *Connection) readMessage() (map[string]any, error) {
	_, err := io.ReadFull(c.conn, c.lengthBuffer)

	if err != nil {
		return nil, err
	}

	jsonBuffer := make([]byte, binary.BigEndian.Uint32(c.lengthBuffer))

	_, err = io.ReadFull(c.conn, jsonBuffer)

	if err != nil {
		return nil, err
	}

	var message map[string]any

	if err = json.Unmarshal(jsonBuffer, &message); err != nil {
		return nil, err
	}

	return message, nil
}

func (c *Connection) writeMessage(message map[string]any) error {
	jsonBytes, err := json.Marshal(message)

	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(c.lengthBuffer, uint32(len(jsonBytes)))

	err = c.writeFull(c.lengthBuffer)

	if err != nil {
		return err
	}

	err = c.writeFull(jsonBytes)

	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) writeFull(data []byte) error {
	total := 0

	for total < len(data) {
		n, err := c.conn.Write(data[total:])

		if err != nil {
			return err
		}

		total += n
	}

	return nil
}
