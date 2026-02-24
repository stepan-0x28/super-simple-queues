package queue

type Queue struct {
	elements []map[string]interface{}
}

func NewQueue() *Queue {
	return &Queue{make([]map[string]interface{}, 0)}
}

func (q *Queue) AddElement(element map[string]interface{}) {
	q.elements = append(q.elements, element)
}

func (q *Queue) ElementsNumber() int {
	return len(q.elements)
}
