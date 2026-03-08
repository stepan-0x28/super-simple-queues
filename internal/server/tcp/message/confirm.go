package message

import "io"

type Confirm struct{}

func NewConfirm() Message {
	return &Confirm{}
}

func (c *Confirm) Type() Type {
	return ConfirmType
}

func (c *Confirm) ReadBody(_ io.Reader) error {
	return nil
}

func (c *Confirm) WriteBody(_ io.Writer) error {
	return nil
}
