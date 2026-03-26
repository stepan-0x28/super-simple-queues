package queue

import (
	"log/slog"
	"sync"
)

type Manager struct {
	queues         map[string]*Queue
	queueChunkSize int
	mutex          sync.Mutex
}

func NewManager(queueChunkSize int) *Manager {
	return &Manager{queues: make(map[string]*Queue, 64), queueChunkSize: queueChunkSize}
}

func (m *Manager) Create(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.queues[key]

	if ok {
		return false
	}

	m.queues[key] = newQueue(m.queueChunkSize, key)

	slog.Info("queue created", slog.Any("key", key))

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

func (m *Manager) Delete(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, ok := m.queues[key]

	if !ok {
		return false
	}

	queue.delete()

	delete(m.queues, key)

	slog.Info("queue deleted", slog.Any("key", key))

	return true
}
