package store

import (
	"strconv"
	"sync"
	"testing"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestStore_Add(t *testing.T) {
	s := NewStore(5)

	t.Run("should assign sequential ids", func(t *testing.T) {
		first := s.Add(domain.Request{Method: "GET", Path: "/a"})
		second := s.Add(domain.Request{Method: "POST", Path: "/b"})

		require.Equal(t, 0, first.ID)
		require.Equal(t, 1, second.ID)
	})

	t.Run("should evict oldest beyond capacity", func(t *testing.T) {
		full := NewStore(3)

		for i := 0; i < 5; i++ {
			full.Add(domain.Request{Path: "/" + strconv.Itoa(i)})
		}

		list := full.List()
		require.Len(t, list, 3)
		require.Equal(t, "/4", list[0].Path)
		require.Equal(t, "/2", list[2].Path)

		_, ok := full.GetByID(0)
		require.False(t, ok)
	})
}

func TestStore_GetByID(t *testing.T) {
	s := NewStore(5)
	added := s.Add(domain.Request{Method: "GET", Path: "/a"})

	t.Run("should get existing id", func(t *testing.T) {
		got, ok := s.GetByID(added.ID)
		require.True(t, ok)
		require.Equal(t, "/a", got.Path)
	})

	t.Run("should not get missing id", func(t *testing.T) {
		_, ok := s.GetByID(99)
		require.False(t, ok)
	})
}

func TestStore_List(t *testing.T) {
	s := NewStore(5)

	t.Run("should list newest-first", func(t *testing.T) {
		s.Add(domain.Request{Method: "GET", Path: "/a"})
		s.Add(domain.Request{Method: "POST", Path: "/b"})

		list := s.List()
		require.Len(t, list, 2)
		require.Equal(t, 1, list[0].ID)
		require.Equal(t, 0, list[1].ID)
	})
}

func TestStore_Reset(t *testing.T) {
	s := NewStore(5)

	t.Run("should clear all requests", func(t *testing.T) {
		s.Add(domain.Request{Path: "/a"})
		s.Reset()

		require.Empty(t, s.List())
	})
}

func TestStore_Concurrency(t *testing.T) {
	s := NewStore(10)

	t.Run("should handle concurrent adds and lists", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.Add(domain.Request{Path: "/x"})
				s.List()
			}()
		}

		wg.Wait()

		require.Len(t, s.List(), 10)
	})
}
