package queue

import (
	"errors"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	if m == nil {
		t.Fatal("a non-nil manager is expected")
	}

	if len(m.GetAll()) != 0 {
		t.Fatal("it is expected that the manager does not have any queues yet")
	}
}

func TestManager_Delete(t *testing.T) {
	m := NewManager()

	queueKey := "test"

	deleted := m.Delete(queueKey)

	if deleted {
		t.Fatal("it was expected that there would be nothing to delete")
	}

	m.Create(queueKey)

	q, _ := m.Get(queueKey)

	const interactingCount = 3

	errChan := make(chan error, interactingCount)

	addedItem := []byte("test")

	go queueInteraction(errChan, func() error { return q.Add(addedItem) })

	go queueInteraction(errChan, func() error { _, err := q.Take(); return err })

	go queueInteraction(errChan, func() error { return q.PutBack(addedItem) })

	time.Sleep(time.Second)

	deleted = m.Delete(queueKey)

	if !deleted {
		t.Fatal("it is expected that the queue has been successfully deleted")
	}

	_, ok := m.Get(queueKey)

	if ok {
		t.Fatal("it is expected that the queue will not exist")
	}

	for i := 0; i < interactingCount; i++ {
		select {
		case err := <-errChan:
			if !errors.Is(err, ErrQueueDeleted) {
				t.Fatal("the error was expected to match ErrQueueDeleted")
			}
		case <-time.After(time.Second):
			t.Fatal("it was expected that an error would be received")
		}
	}
}

func queueInteraction(ch chan error, fn func() error) {
	for {
		if err := fn(); err != nil {
			ch <- err

			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}
