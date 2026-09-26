package store

import (
	"sort"
	"sync"

	"github.com/georgelopez7/peeko/internal/domain"
)

type Store struct {
	mu       sync.Mutex
	maxSize  int
	nextID   int
	requests []domain.Request
}

func NewStore(maxSize int) *Store {
	return &Store{maxSize: maxSize}
}

// Add - appends a request evicting the oldest when full.
func (s *Store) Add(r domain.Request) domain.Request {
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

// GetByID - returns the request by ID.
func (s *Store) GetByID(id int) (domain.Request, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.requests {
		if r.ID == id {
			return r, true
		}
	}

	return domain.Request{}, false
}

// List - returns all requests ordered by newest first.
func (s *Store) List() []domain.Request {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]domain.Request, len(s.requests))
	for i, r := range s.requests {
		out[i] = r
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })

	return out
}

// Reset - removes all requests.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = nil
}
