package tcp

import (
	"bufio"
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

func newConnection(conn net.Conn, bufferSize int) *connection {
	return &connection{
		codec: newCodec(bufio.NewReaderSize(conn, bufferSize), bufio.NewWriterSize(conn, bufferSize)),
	}
}

func (c *connection) init(queueManager *queue.Manager) (message.OperatingMode, *queue.Queue, error) {
	msg, err := c.codec.readMessage()

	if err != nil {
		return message.SendingOperatingMode, nil, err
	}

	initMessage, ok := msg.(*message.Init)

	if !ok {
		return message.SendingOperatingMode, nil, errors.New("expected message type \"Init\"")
	}

	q, ok := queueManager.Get(initMessage.QueueKey)

	if !ok {
		return message.SendingOperatingMode, nil,
			fmt.Errorf("queue with key \"%v\" does not exist", initMessage.QueueKey)
	}

	if err = c.codec.writeMessage(confirmMessage); err != nil {
		return message.SendingOperatingMode, nil, err
	}

	return initMessage.OperatingMode, q, nil
}

func (c *connection) run(operatingMode message.OperatingMode, q *queue.Queue) error {
	if operatingMode == message.SendingOperatingMode {
		return c.readMessages(q)
	}

	return c.writeMessages(q)
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

		if err = q.Add(payloadMessage.Data); err != nil {
			return err
		}

		if err = c.codec.writeMessage(confirmMessage); err != nil {
			return err
		}
	}
}

func (c *connection) writeMessages(q *queue.Queue) error {
	for {
		item, err := q.Take()

		if err != nil {
			return err
		}

		if err = c.codec.writeMessage(message.NewPayloadWithData(item)); err != nil {
			if err = q.PutBack(item); err != nil {
				return err
			}

			return err
		}

		msg, err := c.codec.readMessage()

		if err != nil {
			if err = q.PutBack(item); err != nil {
				return err
			}

			return err
		}

		_, ok := msg.(*message.Confirm)

		if !ok {
			if err = q.PutBack(item); err != nil {
				return err
			}

			return errors.New("expected message type \"Confirm\"")
		}
	}
}
