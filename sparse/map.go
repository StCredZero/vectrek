package sparse

import (
	"errors"
	"fmt"
)

var ErrMissing = errors.New("missing")

type Key uint64

// Map is a structure to hold entities with a specific component type.
type Map[T any] struct {
	// sparse maps entity IDs to their index in the dense list
	sparse map[Key]int
	// dense holds the actual components or entity IDs
	dense   []T
	deleted map[int]struct{}
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		sparse:  make(map[Key]int),
		dense:   make([]T, 0, 16),
		deleted: make(map[int]struct{}),
	}
}

func (s *Map[T]) Add(key Key, value T) {
	for index := range s.deleted {
		s.sparse[key] = index
		delete(s.deleted, index)
		s.dense[index] = value
		return
	}
	denseIndex := len(s.dense)
	s.sparse[key] = denseIndex
	s.dense = append(s.dense, value)
	delete(s.deleted, denseIndex)
}

func (s *Map[T]) Delete(key Key) {
	denseIndex, found := s.sparse[key]
	if !found {
		return
	}
	s.deleted[denseIndex] = struct{}{}
}

func (s *Map[T]) Iterate(fn func(value T) (T, error)) []error {
	errs := make([]error, 0, len(s.dense))
	for i := range s.dense {
		if _, deleted := s.deleted[i]; !deleted {
			updated, err := fn(s.dense[i])
			s.dense[i] = updated
			errs = append(errs, err)
		}
	}
	return errs
}

func (s *Map[T]) Get(key Key) (*T, bool) {
	var result T
	denseIndex, ok := s.sparse[key]
	if !ok {
		return &result, false
	}
	return &s.dense[denseIndex], true
}

func (s *Map[T]) GetErr(key Key) (*T, error) {
	if result, ok := s.Get(key); !ok {
		return nil, fmt.Errorf("mising key %d: %w", key, ErrMissing)
	} else {
		return result, nil
	}
}

func (s *Map[T]) MustGet(key Key) *T {
	result, ok := s.Get(key)
	if !ok {
		return nil
	}
	return result
}
