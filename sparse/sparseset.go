package sparse

import (
	"errors"
	"fmt"
)

var ErrMissing = errors.New("missing")

type Key uint64

type mapEntry[T any] struct {
	deleted bool
	value   T
}

// Map is a structure to hold entities with a specific component type.
type Map[T any] struct {
	// sparse maps entity IDs to their index in the dense list
	sparse map[Key]int
	// dense holds the actual components or entity IDs
	dense   []mapEntry[T]
	deleted map[int]struct{}
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		sparse:  make(map[Key]int),
		dense:   make([]mapEntry[T], 0, 16),
		deleted: make(map[int]struct{}),
	}
}

func (s *Map[T]) Add(key Key, value T) {
	for index := range s.deleted {
		s.sparse[key] = index
		delete(s.deleted, index)
		s.dense[index].value = value
		return
	}
	denseIndex := len(s.dense)
	s.sparse[key] = denseIndex
	s.dense = append(s.dense, mapEntry[T]{deleted: false, value: value})
	delete(s.deleted, denseIndex)
}

func (s *Map[T]) Get(key Key) (T, bool) {
	var result T
	denseIndex, ok := s.sparse[key]
	if !ok {
		return result, false
	}
	entry := s.dense[denseIndex]
	return entry.value, entry.deleted
}

func (s *Map[T]) Update(key Key, updateFunc func(T) (T, error)) error {
	denseIndex, ok := s.sparse[key]
	if !ok {
		return fmt.Errorf("could not find key %d: %w", key, ErrMissing)
	}
	entry := s.dense[denseIndex]
	if entry.deleted {
		return fmt.Errorf("key deleted %d: %w", key, ErrMissing)
	}
	result, err := updateFunc(entry.value)
	if err != nil {
		return fmt.Errorf("error in update: %w", err)
	}
	entry.value = result
	s.dense[denseIndex] = entry
	return nil
}

func (s *Map[T]) Delete(key Key) {
	denseIndex, found := s.sparse[key]
	if !found {
		return
	}
	s.dense[denseIndex].deleted = true
	s.deleted[denseIndex] = struct{}{}
}
