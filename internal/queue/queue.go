package queue

import (
	"sync"
)

type Queue struct {
	mutex    sync.Mutex
	cond     *sync.Cond
	messages []map[string]any
}

func NewQueue() *Queue {
	q := &Queue{messages: make([]map[string]any, 0)}

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(message map[string]any) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, message)

	q.cond.Signal()
}

func (q *Queue) Take() map[string]any {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.messages) == 0 {
		q.cond.Wait()
	}

	message := q.messages[0]

	q.messages = q.messages[1:]

	return message
}

func (q *Queue) Count() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.messages)
}
