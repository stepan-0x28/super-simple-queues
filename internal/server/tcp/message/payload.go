package message

import (
	"encoding/binary"
	"io"
)

type Payload struct {
	Data []byte
}

func NewPayload() Message {
	return &Payload{}
}

func NewPayloadWithData(data []byte) Message {
	return &Payload{data}
}

func (p *Payload) GetType() Type {
	return PayloadType
}

func (p *Payload) ReadBody(reader io.Reader) error {
	var fourByteArray [4]byte

	fourByteBuffer := fourByteArray[:]

	_, err := io.ReadFull(reader, fourByteBuffer)

	if err != nil {
		return err
	}

	p.Data = make([]byte, binary.BigEndian.Uint32(fourByteBuffer))

	_, err = io.ReadFull(reader, p.Data)

	if err != nil {
		return err
	}

	return nil
}

func (p *Payload) WriteBody(writer io.Writer) error {
	var fourByteArray [4]byte

	fourByteBuffer := fourByteArray[:]

	binary.BigEndian.PutUint32(fourByteBuffer, uint32(len(p.Data)))

	_, err := writer.Write(fourByteBuffer)

	if err != nil {
		return err
	}

	_, err = writer.Write(p.Data)

	if err != nil {
		return err
	}

	return nil
}
