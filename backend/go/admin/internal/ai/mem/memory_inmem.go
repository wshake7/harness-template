package mem

import (
	"context"
	"sync"

	"github.com/cloudwego/eino/schema"
)

type inMemMemory struct {
	maxWindowSize int
	sessions      map[string]*inMemSession
	mu            sync.Mutex
}

func NewInMemMemory(cfg InMemConfig) Memory {
	return &inMemMemory{
		maxWindowSize: cfg.MaxWindowSize,
		sessions:      make(map[string]*inMemSession),
	}
}

func (m *inMemMemory) GetSession(id string) Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[id]; ok {
		return s
	}
	s := &inMemSession{
		id:            id,
		maxWindowSize: m.maxWindowSize,
	}
	m.sessions[id] = s
	return s
}

type inMemSession struct {
	id            string
	messages      []*schema.Message
	maxWindowSize int
	mu            sync.Mutex
}

func (s *inMemSession) SetMessages(_ context.Context, msg *schema.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
	if len(s.messages) > s.maxWindowSize {
		excess := len(s.messages) - s.maxWindowSize
		if excess%2 != 0 {
			excess++
		}
		s.messages = s.messages[excess:]
	}
	return nil
}

func (s *inMemSession) GetMessages(_ context.Context) ([]*schema.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages, nil
}
