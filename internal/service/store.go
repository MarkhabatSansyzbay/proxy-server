package service

import (
	"sync"

	"proxyserver/internal/models"
)

type Store struct {
	store map[string]models.ProxyResponse
	mu    *sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]models.ProxyResponse),
		mu:    &sync.RWMutex{},
	}
}

func (s *Store) save(request string, response models.ProxyResponse) {
	s.mu.Lock()
	s.store[request] = response
	s.mu.Unlock()
}

func (s *Store) get(request string) (models.ProxyResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, ok := s.store[request]
	return res, ok
}
