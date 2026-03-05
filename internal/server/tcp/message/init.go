package message

import "io"

type Init struct {
	Mode     bool
	QueueKey string
}

func NewInit() Message {
	return &Init{}
}

func (i *Init) GetType() Type {
	return InitType
}

func (i *Init) ReadBody(reader io.Reader) error {
	var oneByteArray [1]byte

	oneByteBuffer := oneByteArray[:]

	_, err := io.ReadFull(reader, oneByteBuffer)

	if err != nil {
		return err
	}

	i.Mode = oneByteBuffer[0] != 0

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
