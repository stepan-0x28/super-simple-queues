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

	actualCount, _ := q.Count()

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

	concurrentQueueInteraction(&wg, receiverCount, receivedItemsCount, func() { _, _ = q.Take() })

	// let's wait to increase the likelihood that q.Take() will call q.cond.Wait()
	time.Sleep(time.Second)

	addedItem := []byte("test")

	go concurrentQueueInteraction(&wg, senderCount, sentItemsCount, func() { _ = q.Add(addedItem) })

	go concurrentQueueInteraction(&wg, receiverCount2, receivedItemsCount2, func() { _, _ = q.Take() })

	go concurrentQueueInteraction(&wg, returnerCount, returnedItemsCount, func() { _ = q.PutBack(addedItem) })

	wg.Wait()

	expectedCount := senderCount*sentItemsCount -
		receiverCount*receivedItemsCount -
		receiverCount2*receivedItemsCount2 +
		returnerCount*returnedItemsCount

	actualCount, _ := q.Count()

	if expectedCount != actualCount {
		t.Fatalf("%d elements were expected, but there were %d", expectedCount, actualCount)
	}
}

func TestQueue_Sequence(t *testing.T) {
	q := newQueue()

	const itemsCount = 10

	for i := 0; i < itemsCount; i++ {
		_ = q.Add([]byte(strconv.Itoa(i)))
	}

	for i := 0; i < itemsCount; i++ {
		_ = q.PutBack([]byte(strconv.Itoa(i)))
	}

	check := func(i int) {
		t.Helper()

		expected := strconv.Itoa(i)

		item, _ := q.Take()

		received := string(item)

		if received != expected {
			t.Fatalf("expected \"%v\", received \"%v\"", expected, received)
		}
	}

	for i := itemsCount - 1; i >= 0; i-- {
		check(i)
	}

	for i := 0; i < itemsCount; i++ {
		check(i)
	}
}

func TestQueue_Take(t *testing.T) {
	q := newQueue()

	takenItemChan := make(chan []byte)

	go func() {
		for {
			item, _ := q.Take()

			takenItemChan <- item
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

	_ = q.Add(addedItem)

	waitData()

	_ = q.PutBack(addedItem)

	waitData()
}

func concurrentQueueInteraction(wg *sync.WaitGroup, goroutineCount int, itemCount int, fn func()) {
	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < itemCount; j++ {
				fn()
			}
		}()
	}
}
