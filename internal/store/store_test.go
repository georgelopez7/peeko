package store

import (
	"strconv"
	"sync"
	"testing"

	"github.com/george-lopez/peeko/internal/domain"
)

func TestStoreAddListGet(t *testing.T) {
	s := NewStore(5)

	first := s.Add(domain.CapturedRequest{Method: "GET", Path: "/a"})
	second := s.Add(domain.CapturedRequest{Method: "POST", Path: "/b"})

	if first.ID != 0 || second.ID != 1 {
		t.Fatalf("want ids 0,1 got %d,%d", first.ID, second.ID)
	}

	list := s.List()
	if len(list) != 2 || list[0].ID != 1 || list[1].ID != 0 {
		t.Fatalf("want newest-first [1,0], got %v", list)
	}

	got, ok := s.Get(0)
	if !ok || got.Path != "/a" {
		t.Fatalf("want /a, got %q ok=%v", got.Path, ok)
	}

	if _, ok := s.Get(99); ok {
		t.Fatal("want missing id to not be found")
	}
}

func TestStoreEviction(t *testing.T) {
	s := NewStore(3)

	for i := 0; i < 5; i++ {
		s.Add(domain.CapturedRequest{Path: "/" + strconv.Itoa(i)})
	}

	list := s.List()
	if len(list) != 3 {
		t.Fatalf("want 3 entries, got %d", len(list))
	}
	if list[0].Path != "/4" || list[2].Path != "/2" {
		t.Fatalf("want newest /4 oldest /2, got %q..%q", list[0].Path, list[2].Path)
	}

	if _, ok := s.Get(0); ok {
		t.Fatal("want evicted id to not be found")
	}
}

func TestStoreClear(t *testing.T) {
	s := NewStore(5)
	s.Add(domain.CapturedRequest{Path: "/a"})

	s.Clear()
	if got := s.List(); len(got) != 0 {
		t.Fatalf("want empty after clear, got %d", len(got))
	}
}

func TestStoreConcurrent(t *testing.T) {
	s := NewStore(10)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Add(domain.CapturedRequest{Path: "/x"})
			s.List()
		}()
	}
	wg.Wait()

	if got := s.List(); len(got) != 10 {
		t.Fatalf("want 10 after concurrent adds, got %d", len(got))
	}
}
