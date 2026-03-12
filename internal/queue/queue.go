package queue

import (
	"errors"
	"sync"
)

var ErrQueueDeleted = errors.New("interaction with the deleted queue")

type Queue struct {
	mutex   sync.Mutex
	cond    *sync.Cond
	deleted bool
	// TODO the slices need to be replaced with a different structure
	items [][]byte
}

func newQueue() *Queue {
	q := &Queue{items: make([][]byte, 0, 1024)}

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(item []byte) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return ErrQueueDeleted
	}

	q.items = append(q.items, item)

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

		if len(q.items) == 0 {
			q.cond.Wait()
		} else {
			break
		}
	}

	item := q.items[0]

	q.items[0] = nil

	q.items = q.items[1:]

	return item, nil
}

func (q *Queue) PutBack(item []byte) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return ErrQueueDeleted
	}

	q.items = append(q.items, nil)

	copy(q.items[1:], q.items[:len(q.items)-1])

	q.items[0] = item

	q.cond.Signal()

	return nil
}

func (q *Queue) Count() (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.deleted {
		return 0, ErrQueueDeleted
	}

	return len(q.items), nil
}

func (q *Queue) delete() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.deleted = true

	q.cond.Broadcast()
}
