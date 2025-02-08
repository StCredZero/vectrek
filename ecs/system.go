package ecs

import "github.com/StCredZero/vectrek/sparse"

type SystemID uint64

const (
	SystemPosition SystemID = iota
	SystemMotion
	SystemHelm
	SystemSprite
)

type System[T Component] struct {
	ID      SystemID
	Members sparse.Map[T]
}

func (s System[T]) Add(entity EntityID, obj T) {
	s.Members.Add(entity, obj)
}
func (s System[T]) Iterate(fn func(each *T) error) {
	s.Iterate(fn)
}
