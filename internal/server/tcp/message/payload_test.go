package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestPayload_WriteBody(t *testing.T) {
	tests := []struct {
		name               string
		data               []byte
		expectedDataLength uint32
	}{{
		name:               "correct message (4 bytes)",
		data:               []byte("test"),
		expectedDataLength: 4,
	}, {
		name:               "correct message (6 bytes)",
		data:               []byte("test-2"),
		expectedDataLength: 6,
	}, {
		name:               "correct message (no data)",
		data:               []byte{},
		expectedDataLength: 0,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := test.data

			buffer := bytes.Buffer{}

			err := NewPayloadWithData(data).WriteBody(&buffer)

			if err != nil {
				t.Fatalf("writing error, %v", err)
			}

			bufferBytes := buffer.Bytes()

			writtenDataLength := binary.BigEndian.Uint32(bufferBytes[:4])
			expectedDataLength := test.expectedDataLength

			if writtenDataLength != expectedDataLength {
				t.Fatalf("data length is %v, expected %v", writtenDataLength, expectedDataLength)
			}

			writtenData := string(bufferBytes[4:])
			expectedData := string(data)

			if writtenData != expectedData {
				t.Fatalf("data is %v, expected %v", writtenData, expectedData)
			}
		})
	}
}

func TestPayload_ReadBody(t *testing.T) {
	tests := []struct {
		name          string
		readableBytes []byte
		expectedError error
		expectedData  []byte
	}{{
		name:          "correct message (4 bytes)",
		readableBytes: []byte{0, 0, 0, 4, 't', 'e', 's', 't'},
		expectedData:  []byte("test"),
	}, {
		name:          "correct message (6 bytes)",
		readableBytes: []byte{0, 0, 0, 6, 't', 'e', 's', 't', '-', '2'},
		expectedData:  []byte("test-2"),
	}, {
		name:          "incomplete message",
		readableBytes: []byte{0, 0, 0, 4, 't', 'e'},
		expectedError: io.ErrUnexpectedEOF,
	}, {
		name:          "empty message",
		readableBytes: []byte{},
		expectedError: io.EOF,
	}, {
		name:          "incorrect message",
		readableBytes: []byte{0},
		expectedError: io.ErrUnexpectedEOF,
	}, {
		name:          "message without data",
		readableBytes: []byte{0, 0, 0, 0},
		expectedData:  []byte{},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payloadMessage := NewPayload().(*Payload)

			err := payloadMessage.ReadBody(bytes.NewReader(test.readableBytes))
			expectedError := test.expectedError

			if !errors.Is(err, expectedError) {
				t.Fatalf("expected error \"%v\", received \"%v\"", expectedError, err)
			}

			expectedData := test.expectedData

			if expectedData != nil {
				data := string(payloadMessage.Data)
				expectedDataString := string(expectedData)

				if data != expectedDataString {
					t.Fatalf("expected data \"%v\", received \"%v\"", expectedDataString, data)
				}
			}
		})
	}
}
