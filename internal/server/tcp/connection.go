package tcp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"super-simple-queues/internal/queue"
)

type Connection struct {
	conn  net.Conn
	queue *queue.Queue
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

func (c *Connection) Run(queueManager *queue.Manager) error {
	initMessage, err := c.readMessage()

	if err != nil {
		return err
	}

	queueKey, err := getStringValue(initMessage, "queue_key")

	if err != nil {
		return err
	}

	mode, err := getStringValue(initMessage, "mode")

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
	}
}

func (c *Connection) sendMessages() error {
	for {
		if err := c.writeMessage(c.queue.Take()); err != nil {
			return err
		}
	}
}

func (c *Connection) readMessage() (map[string]interface{}, error) {
	lengthBuffer := make([]byte, 4)

	_, err := io.ReadFull(c.conn, lengthBuffer)

	if err != nil {
		return nil, err
	}

	jsonBuffer := make([]byte, binary.BigEndian.Uint32(lengthBuffer))

	_, err = io.ReadFull(c.conn, jsonBuffer)

	if err != nil {
		return nil, err
	}

	var message map[string]interface{}

	if err = json.Unmarshal(jsonBuffer, &message); err != nil {
		return nil, err
	}

	return message, nil
}

func (c *Connection) writeMessage(message map[string]interface{}) error {
	jsonBytes, err := json.Marshal(message)

	if err != nil {
		return err
	}

	lengthBuffer := make([]byte, 4)

	binary.BigEndian.PutUint32(lengthBuffer, uint32(len(jsonBytes)))

	err = c.writeFull(lengthBuffer)

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

func getStringValue(message map[string]interface{}, key string) (string, error) {
	value, ok := message[key]

	if !ok {
		return "", fmt.Errorf("the \"%v\" key is missing", key)
	}

	stringValue, ok := value.(string)

	if !ok {
		return "", fmt.Errorf("the \"%v\" key is incorrect", key)
	}

	return stringValue, nil
}
