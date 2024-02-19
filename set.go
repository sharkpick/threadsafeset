package threadsafeset

import (
	"sync"

	"github.com/sharkpick/simpleset"

	"golang.org/x/exp/constraints"
)

type Set[T constraints.Ordered] struct {
	set   *simpleset.Set[T]
	mutex sync.RWMutex
}

func New[T constraints.Ordered]() *Set[T] {
	return &Set[T]{
		set: simpleset.New[T](),
	}
}

func NewFromSlice[T constraints.Ordered](slice []T) *Set[T] {
	return &Set[T]{
		set: simpleset.NewFromSlice(slice),
	}
}

func (s *Set[T]) Add(t T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set.Add(t)
}

func (s *Set[T]) AddSlice(slice []T) []bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set.AddSlice(slice)
}

func (s *Set[T]) Contains(t T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.Contains(t)
}

func (s *Set[T]) ContainsSlice(slice []T) []bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.ContainsSlice(slice)
}

func (s *Set[T]) Drop(t T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set.Drop(t)
}

func (s *Set[T]) DropSlice(slice []T) []bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set.DropSlice(slice)
}

func (s *Set[T]) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.Len()
}

func (s *Set[T]) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set.Reset()
}

func (s *Set[T]) Slice() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.Slice()
}
