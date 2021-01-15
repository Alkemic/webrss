package repository

import (
	"fmt"
	"sync"
	"time"
)

type session struct {
	due  time.Time
	data map[string]interface{}
}

type SessionRepository struct {
	sync.RWMutex
	age   time.Duration
	store map[string]session // session-id => data
}

func NewSessionRepository(age time.Duration) *SessionRepository {
	return &SessionRepository{
		age:   age,
		store: make(map[string]session),
	}
}

func (r *SessionRepository) Get(sessionID string) (map[string]interface{}, error) {
	r.RLock()
	session, ok := r.store[sessionID]
	if !ok {
		r.RUnlock()
		return nil, ErrNotFound
	}
	r.RUnlock()
	if session.due.Before(time.Now()) {
		if err := r.Delete(sessionID); err != nil {
			return nil, fmt.Errorf("cannot delete expired session: %w", err)
		}
		return nil, ErrNotFound
	}
	return session.data, nil
}

func (r *SessionRepository) Set(sessionID string, data map[string]interface{}) error {
	r.Lock()
	r.store[sessionID] = session{data: data, due: time.Now().Add(r.age)}
	r.Unlock()
	return nil
}

func (r *SessionRepository) Delete(sessionID string) error {
	r.Lock()
	delete(r.store, sessionID)
	r.Unlock()
	return nil
}
