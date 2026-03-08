package queue

import "sync"

type Queue struct {
	mutex sync.Mutex
	cond  *sync.Cond
	// TODO the slices need to be replaced with a different structure
	items [][]byte
}

func NewQueue() *Queue {
	q := &Queue{items: make([][]byte, 0, 1024)}

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(item []byte) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, item)

	q.cond.Signal()
}

func (q *Queue) Take() []byte {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item := q.items[0]

	q.items[0] = nil

	q.items = q.items[1:]

	return item
}

func (q *Queue) PutBack(item []byte) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, nil)

	copy(q.items[1:], q.items[:len(q.items)-1])

	q.items[0] = item
}

func (q *Queue) Count() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items)
}
