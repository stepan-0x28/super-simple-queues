package message

import "io"

type Type byte

const (
	InitType    Type = 1
	PayloadType Type = 2
	ConfirmType Type = 3
)

type Message interface {
	GetType() Type
	ReadBody(io.Reader) error
	WriteBody(io.Writer) error
}
