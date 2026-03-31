package queue

import (
	"errors"
	"log/slog"
	"sync"
)

var (
	ErrQueueDeleted   = errors.New("interaction with the deleted queue")
	ErrEmptyQueueItem = errors.New("queue item is empty")
)

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
	key       string
	mutex     sync.Mutex
	cond      *sync.Cond
}

func newQueue(chunkSize int, key string) *Queue {
	q := &Queue{chunkSize: chunkSize, key: key}

	q.head = newChunk(q.chunkSize)
	q.tail = q.head

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(item []byte) error {
	if len(item) == 0 {
		return ErrEmptyQueueItem
	}

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

		slog.Debug("new tail chunk created", slog.Any("queue", q.key),
			slog.Any("chunk_size", q.chunkSize), slog.Any("count", q.count))
	}

	q.tail.data[q.tailPos] = item

	q.tailPos++

	q.count++

	slog.Debug("item added", slog.Any("queue", q.key), slog.Any("chunk_size", q.chunkSize),
		slog.Any("count", q.count))

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

		slog.Debug("head chunk shifted", slog.Any("queue", q.key), slog.Any("chunk_size", q.chunkSize),
			slog.Any("count", q.count))
	}

	item := q.head.data[q.headPos]

	q.head.data[q.headPos] = nil

	q.headPos++

	q.count--

	slog.Debug("item taken", slog.Any("queue", q.key), slog.Any("chunk_size", q.chunkSize),
		slog.Any("count", q.count))

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

		slog.Debug("new head chunk created", slog.Any("queue", q.key),
			slog.Any("chunk_size", q.chunkSize), slog.Any("count", q.count))
	}

	q.headPos--

	q.head.data[q.headPos] = item

	q.count++

	slog.Debug("item put back", slog.Any("queue", q.key), slog.Any("chunk_size", q.chunkSize),
		slog.Any("count", q.count))

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
