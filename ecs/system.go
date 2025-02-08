package ecs

import "github.com/StCredZero/vectrek/sparse"

type SystemID uint64

const (
	SystemPosition SystemID = iota
	SystemMotion
	SystemHelm
	SystemSprite
)

type System interface {
	SystemID() SystemID
}

type MSystem[T any] struct {
	SID     SystemID
	Members sparse.Map[T]
}

func (s MSystem[T]) SystemID() SystemID {
	return s.SID
}
func (s MSystem[T]) Add(entity EntityID, obj T) {
	s.Members.Add(entity, obj)
}
func (s MSystem[T]) Iterate(fn func(each *T) error) {
	s.Iterate(fn)
}
func (s MSystem[T]) Get(entity EntityID) (*T, bool) {
	return s.Members.Get(entity)
}
