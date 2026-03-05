package tcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"super-simple-queues/internal/queue"
	"super-simple-queues/internal/server/tcp/message"
)

var confirmMessage = message.NewConfirm()

type connection struct {
	codec *codec
}

func newConnection(conn net.Conn) *connection {
	return &connection{
		codec: newCodec(bufio.NewReaderSize(conn, 1024), bufio.NewWriterSize(conn, 1024)),
	}
}

func (c *connection) run(queueManager *queue.Manager) error {
	msg, err := c.codec.readMessage()

	if err != nil {
		return err
	}

	initMessage, ok := msg.(*message.Init)

	if !ok {
		return errors.New("expected message type \"Init\"")
	}

	q, ok := queueManager.Get(initMessage.QueueKey)

	if !ok {
		return fmt.Errorf("queue with key \"%v\" does not exist", initMessage.QueueKey)
	}

	err = c.codec.writeMessage(confirmMessage)

	if err != nil {
		return err
	}

	if initMessage.Mode {
		err = c.readMessages(q)
	} else {
		err = c.writeMessages(q)
	}

	return err
}

func (c *connection) readMessages(q *queue.Queue) error {
	for {
		msg, err := c.codec.readMessage()

		if err != nil {
			return err
		}

		payloadMessage, ok := msg.(*message.Payload)

		if !ok {
			return errors.New("expected message type \"Payload\"")
		}

		if !json.Valid(payloadMessage.Data) {
			return errors.New("the message data is invalid json")
		}

		q.Add(payloadMessage.Data)

		err = c.codec.writeMessage(confirmMessage)

		if err != nil {
			return err
		}
	}
}

func (c *connection) writeMessages(q *queue.Queue) error {
	for {
		item := q.Take()

		err := c.codec.writeMessage(message.NewPayloadWithData(item))

		if err != nil {
			q.PutBack(item)

			return err
		}

		msg, err := c.codec.readMessage()

		_, ok := msg.(*message.Confirm)

		if !ok {
			q.PutBack(item)

			return errors.New("expected message type \"Confirm\"")
		}
	}
}
