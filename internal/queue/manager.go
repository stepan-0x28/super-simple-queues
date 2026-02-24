package queue

type Manager struct {
	queues map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{make(map[string]*Queue)}
}

func (m *Manager) CreateQueue(key string) {
	m.queues[key] = NewQueue()
}

func (m *Manager) GetQueue(key string) *Queue {
	queue, ok := m.queues[key]

	if !ok {
		return nil
	}

	return queue
}

func (m *Manager) GetEverything() map[string]int {
	data := make(map[string]int)

	for k, v := range m.queues {
		data[k] = v.ElementsNumber()
	}

	return data
}
