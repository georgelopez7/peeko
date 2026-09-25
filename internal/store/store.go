package store

import (
	"sort"
	"sync"

	"github.com/george-lopez/peeko/internal/domain"
)

// Store - a thread-safe, bounded in-memory collection of captured requests.
type Store struct {
	mu       sync.Mutex
	maxSize  int
	nextID   int
	requests []domain.CapturedRequest
}

// NewStore - creates a store that keeps the last maxSize requests.
func NewStore(maxSize int) *Store {
	return &Store{maxSize: maxSize}
}

// Add - appends a captured request, evicting the oldest when full.
func (s *Store) Add(r domain.CapturedRequest) domain.CapturedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	r.ID = s.nextID
	s.nextID++
	s.requests = append(s.requests, r)
	if len(s.requests) > s.maxSize {
		s.requests = s.requests[len(s.requests)-s.maxSize:]
	}

	return r
}

// List - returns all captured requests, newest first.
func (s *Store) List() []domain.CapturedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]domain.CapturedRequest, len(s.requests))
	for i, r := range s.requests {
		out[i] = r
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })

	return out
}

// Get - returns the captured request with the given id; ok is false when missing.
func (s *Store) Get(id int) (domain.CapturedRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.requests {
		if r.ID == id {
			return r, true
		}
	}

	return domain.CapturedRequest{}, false
}

// Clear - removes all captured requests.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = nil
}
