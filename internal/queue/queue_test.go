package queue

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue()

	if q == nil {
		t.Fatal("non-nil queue expected")
	}

	actualCount := q.Count()

	if actualCount != 0 {
		t.Fatalf("an empty queue is expected, there are %d items in the queue", actualCount)
	}
}

func TestQueue_Concurrency(t *testing.T) {
	var wg sync.WaitGroup

	// values have been selected at which the queue at the end will not be empty
	const (
		senderCount, sentItemsCount         = 64, 128
		receiverCount, receivedItemsCount   = 16, 32
		receiverCount2, receivedItemsCount2 = 16, 32
		returnerCount, returnedItemsCount   = 4, 8
	)

	q := NewQueue()

	wg.Add(1)
	go func() {
		defer wg.Done()
		interactWithQueue(t, receiverCount, receivedItemsCount, func() { q.Take() })
	}()

	// let's wait to increase the likelihood that q.Take() will call q.cond.Wait()
	time.Sleep(time.Second)

	addedItem := []byte("test")

	wg.Add(1)
	go func() {
		defer wg.Done()
		interactWithQueue(t, senderCount, sentItemsCount, func() { q.Add(addedItem) })
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		interactWithQueue(t, receiverCount2, receivedItemsCount2, func() { q.Take() })
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		interactWithQueue(t, returnerCount, returnedItemsCount, func() { q.PutBack(addedItem) })
	}()

	wg.Wait()

	expectedCount := senderCount*sentItemsCount -
		receiverCount*receivedItemsCount -
		receiverCount2*receivedItemsCount2 +
		returnerCount*returnedItemsCount

	actualCount := q.Count()

	if expectedCount != actualCount {
		t.Fatalf("%d elements were expected, but there were %d", expectedCount, actualCount)
	}
}

func TestQueue_Sequence(t *testing.T) {
	q := NewQueue()

	for i := 0; i < 10; i++ {
		q.Add([]byte(strconv.Itoa(i)))
	}

	for i := 0; i < 10; i++ {
		q.PutBack([]byte(strconv.Itoa(i)))
	}

	check := func(i int) {
		t.Helper()

		expected := strconv.Itoa(i)
		received := string(q.Take())

		if received != expected {
			t.Fatalf("expected \"%v\", received \"%v\"", expected, received)
		}
	}

	for i := 9; i >= 0; i-- {
		check(i)
	}

	for i := 0; i < 10; i++ {
		check(i)
	}
}

func interactWithQueue(t *testing.T, goroutineCount int, itemCount int, fn func()) {
	t.Helper()

	var wg sync.WaitGroup

	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < itemCount; j++ {
				fn()
			}
		}()
	}

	wg.Wait()
}
