package queue

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := newQueue()

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

	q := newQueue()

	wg.Add(senderCount + receiverCount + receiverCount2 + returnerCount)

	interactWithQueue(&wg, receiverCount, receivedItemsCount, func() { q.Take() })

	// let's wait to increase the likelihood that q.Take() will call q.cond.Wait()
	time.Sleep(time.Second)

	addedItem := []byte("test")

	go interactWithQueue(&wg, senderCount, sentItemsCount, func() { q.Add(addedItem) })

	go interactWithQueue(&wg, receiverCount2, receivedItemsCount2, func() { q.Take() })

	go interactWithQueue(&wg, returnerCount, returnedItemsCount, func() { q.PutBack(addedItem) })

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
	q := newQueue()

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

func TestQueue_Take(t *testing.T) {
	q := newQueue()

	takenItemChan := make(chan []byte)

	go func() {
		for {
			takenItemChan <- q.Take()
		}
	}()

	select {
	case <-takenItemChan:
		t.Fatal("the item take function was expected to block on q.cond.Wait()")
	case <-time.After(time.Second):
	}

	addedItem := []byte("test")

	waitData := func() {
		t.Helper()

		select {
		case <-takenItemChan:
		case <-time.After(time.Second):
			t.Fatal("the item was expected to be received immediately")
		}
	}

	q.Add(addedItem)

	waitData()

	q.PutBack(addedItem)

	waitData()
}

func interactWithQueue(wg *sync.WaitGroup, goroutineCount int, itemCount int, fn func()) {
	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < itemCount; j++ {
				fn()
			}
		}()
	}
}
