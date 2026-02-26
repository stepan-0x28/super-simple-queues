package queue

import (
	"sync"
)

type Queue struct {
	mutex    sync.Mutex
	cond     *sync.Cond
	messages []map[string]interface{}
}

func NewQueue() *Queue {
	q := &Queue{messages: make([]map[string]interface{}, 0)}

	q.cond = sync.NewCond(&q.mutex)

	return q
}

func (q *Queue) Add(message map[string]interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, message)

	q.cond.Signal()
}

func (q *Queue) Take() map[string]interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.messages) == 0 {
		q.cond.Wait()
	}

	message := q.messages[0]

	q.messages = q.messages[1:]

	return message
}

func (q *Queue) GetQuantity() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.messages)
}
