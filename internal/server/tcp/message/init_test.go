package message

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestInit_ReadBody(t *testing.T) {
	tests := []struct {
		name                  string
		readableBytes         []byte
		expectedError         error
		expectedOperatingMode OperatingMode
		expectedQueueKey      string
	}{{
		name:                  "correct message (receiving operating mode)",
		readableBytes:         []byte{0, 4, 't', 'e', 's', 't'},
		expectedOperatingMode: ReceivingOperatingMode,
		expectedQueueKey:      "test",
	}, {
		name:                  "correct message (sending operating mode)",
		readableBytes:         []byte{1, 6, 't', 'e', 's', 't', '-', '2'},
		expectedOperatingMode: SendingOperatingMode,
		expectedQueueKey:      "test-2",
	}, {
		name:          "incomplete queue key",
		readableBytes: []byte{0, 4, 't', 'e'},
		expectedError: io.ErrUnexpectedEOF,
	}, {
		name:          "empty message",
		readableBytes: []byte{},
		expectedError: io.EOF,
	}, {
		name:          "incorrect message",
		readableBytes: []byte{0},
		expectedError: io.EOF,
	}, {
		name:          "message without a queue key",
		readableBytes: []byte{0, 0},
	}, {
		name:          "incorrect operating mode",
		readableBytes: []byte{255},
		expectedError: ErrIncorrectOperatingMode,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initMessage := NewInit().(*Init)

			err := initMessage.ReadBody(bytes.NewReader(test.readableBytes))
			expectedError := test.expectedError

			if !errors.Is(err, expectedError) {
				t.Fatalf("expected error \"%v\", received \"%v\"", expectedError, err)
			}

			expectedQueueKey := test.expectedQueueKey

			if expectedQueueKey != "" {
				operatingMode := initMessage.OperatingMode
				expectedOperatingMode := test.expectedOperatingMode

				if operatingMode != expectedOperatingMode {
					t.Fatalf("expected operating mode \"%v\", received \"%v\"", expectedOperatingMode,
						operatingMode)
				}

				queueKey := initMessage.QueueKey

				if queueKey != expectedQueueKey {
					t.Fatalf("expected queue key \"%v\", received \"%v\"", expectedQueueKey, queueKey)
				}
			}
		})
	}
}
