package queue

import (
	"errors"
	"sync"
)

var ErrQueueDeleted = errors.New("interaction with the deleted queue")

type chunk struct {
	data [][]byte
	next *chunk
}

func newChunk(size int) *chunk {
	return &chunk{data: make([][]byte, size)}
}

type Queue struct {
	head      *chunk
	tail      *chunk
	headPos   int
	tailPos   int
	count     int
	deleted   bool
	chunkSize int
	mutex     sync.Mutex
	cond      *sync.Cond
}

func newQueue(chunkSize int) *Queue {
	q := &Queue{chunkSize: chunkSize}

	q.head = newChunk(q.chunkSize)
	q.tail = q.head

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(item []byte) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return ErrQueueDeleted
	}

	if q.tailPos == q.chunkSize {
		q.tailPos = 0

		c := newChunk(q.chunkSize)

		q.tail.next = c
		q.tail = c
	}

	q.tail.data[q.tailPos] = item

	q.tailPos++

	q.count++

	q.cond.Signal()

	return nil
}

func (q *Queue) Take() ([]byte, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for {
		if q.deleted {
			return nil, ErrQueueDeleted
		}

		if q.count == 0 {
			q.cond.Wait()
		} else {
			break
		}
	}

	if q.headPos == q.chunkSize {
		q.headPos = 0

		q.head = q.head.next
	}

	item := q.head.data[q.headPos]

	q.head.data[q.headPos] = nil

	q.headPos++

	q.count--

	return item, nil
}

func (q *Queue) PutBack(item []byte) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return ErrQueueDeleted
	}

	if q.headPos == 0 {
		q.headPos = q.chunkSize

		c := newChunk(q.chunkSize)

		c.next = q.head

		q.head = c
	}

	q.headPos--

	q.head.data[q.headPos] = item

	q.count++

	q.cond.Signal()

	return nil
}

func (q *Queue) Count() (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return 0, ErrQueueDeleted
	}

	return q.count, nil
}

func (q *Queue) delete() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.deleted = true

	q.cond.Broadcast()
}
