package tcp

import (
	"bufio"
	"errors"
	"super-simple-queues/internal/server/tcp/message"
)

var messageFactory = map[message.Type]func() message.Message{
	message.InitType:    message.NewInit,
	message.PayloadType: message.NewPayload,
	message.ConfirmType: message.NewConfirm,
}

type codec struct {
	bufferedReader *bufio.Reader
	bufferedWriter *bufio.Writer
}

func newCodec(bufferedReader *bufio.Reader, bufferedWriter *bufio.Writer) *codec {
	return &codec{bufferedReader, bufferedWriter}
}

func (c *codec) readMessage() (message.Message, error) {
	headerByte, err := c.bufferedReader.ReadByte()

	if err != nil {
		return nil, err
	}

	constructor, ok := messageFactory[message.Type(headerByte)]

	if !ok {
		return nil, errors.New("unknown message type")
	}

	msg := constructor()

	err = msg.ReadBody(c.bufferedReader)

	if err != nil {
		return nil, err
	}

	return msg, err
}

func (c *codec) writeMessage(msg message.Message) error {
	err := c.bufferedWriter.WriteByte(byte(msg.GetType()))

	if err != nil {
		return err
	}

	err = msg.WriteBody(c.bufferedWriter)

	if err != nil {
		return err
	}

	err = c.bufferedWriter.Flush()

	if err != nil {
		return err
	}

	return nil
}
