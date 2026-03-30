package message

import (
	"errors"
	"io"
)

type OperatingMode bool

const (
	ReceivingOperatingMode OperatingMode = false
	SendingOperatingMode   OperatingMode = true
)

var ErrIncorrectOperatingMode = errors.New("incorrect operating mode")

type Init struct {
	OperatingMode OperatingMode
	QueueKey      string
}

func NewInit() Message {
	return &Init{}
}

func (i *Init) Type() Type {
	return InitType
}

func (i *Init) ReadBody(reader io.Reader) error {
	var oneByteArray [1]byte

	oneByteBuffer := oneByteArray[:]

	_, err := io.ReadFull(reader, oneByteBuffer)

	if err != nil {
		return err
	}

	operatingMode := oneByteBuffer[0]

	if operatingMode > 1 {
		return ErrIncorrectOperatingMode
	}

	i.OperatingMode = operatingMode == 1

	_, err = io.ReadFull(reader, oneByteBuffer)

	if err != nil {
		return err
	}

	queueKeyBuffer := make([]byte, oneByteBuffer[0])

	_, err = io.ReadFull(reader, queueKeyBuffer)

	if err != nil {
		return err
	}

	i.QueueKey = string(queueKeyBuffer)

	return nil
}

func (i *Init) WriteBody(_ io.Writer) error {
	return nil
}
