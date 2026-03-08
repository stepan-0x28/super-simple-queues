package queue

import "sync"

type Manager struct {
	mutex  sync.Mutex
	queues map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{queues: make(map[string]*Queue)}
}

func (m *Manager) Create(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.queues[key]

	if ok {
		return false
	}

	m.queues[key] = NewQueue()

	return true
}

func (m *Manager) Get(key string) (*Queue, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, ok := m.queues[key]

	return queue, ok
}

func (m *Manager) ItemsCounts() map[string]int {
	m.mutex.Lock()

	queuesSnapshot := make(map[string]*Queue, len(m.queues))

	for key, queue := range m.queues {
		queuesSnapshot[key] = queue
	}

	m.mutex.Unlock()

	itemsCount := make(map[string]int, len(queuesSnapshot))

	for key, queue := range queuesSnapshot {
		itemsCount[key] = queue.Count()
	}

	return itemsCount
}
