package queue

import "sync"

type Manager struct {
	mutex  sync.Mutex
	queues map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{queues: make(map[string]*Queue)}
}

func (m *Manager) Create(key string) {
	queue := NewQueue()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.queues[key]; ok {
		return
	}

	m.queues[key] = queue
}

func (m *Manager) Get(key string) (*Queue, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, ok := m.queues[key]

	return queue, ok
}

func (m *Manager) MessagesCount() map[string]int {
	m.mutex.Lock()

	queuesSnapshot := make(map[string]*Queue, len(m.queues))

	for key, queue := range m.queues {
		queuesSnapshot[key] = queue
	}

	m.mutex.Unlock()

	messagesCount := make(map[string]int, len(queuesSnapshot))

	for key, queue := range queuesSnapshot {
		messagesCount[key] = queue.Count()
	}

	return messagesCount
}
