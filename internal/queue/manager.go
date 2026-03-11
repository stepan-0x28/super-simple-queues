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

	m.queues[key] = newQueue()

	return true
}

func (m *Manager) Get(key string) (*Queue, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, ok := m.queues[key]

	return queue, ok
}

func (m *Manager) GetAll() map[string]*Queue {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queues := make(map[string]*Queue, len(m.queues))

	for key, queue := range m.queues {
		queues[key] = queue
	}

	return queues
}
