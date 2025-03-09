package ecs

import (
	"fmt"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/sparse"
)

type SMSystem[T ecstypes.Component] struct {
	Map    *sparse.Map[T]
	Update func(each T) (T, error)
}

func NewSMSystem[T ecstypes.Component](update func(each T) (T, error)) *SMSystem[T] {
	return &SMSystem[T]{
		Map:    sparse.NewMap[T](),
		Update: update,
	}
}

func (s *SMSystem[T]) IsSystem() {}

//goland:noinspection GoDfaNilDereference
func (s *SMSystem[T]) SystemID() ecstypes.SystemID {
	var zero T
	return zero.SystemID()
}

func (s *SMSystem[T]) AddComponent(e ecstypes.EntityID, component T) error {
	s.Map.Add(sparse.Key(e), component)
	return nil
}

func (s *SMSystem[T]) GetComponent(e ecstypes.EntityID) (*T, bool) {
	result, ok := s.Map.Get(sparse.Key(e))
	if !ok {
		return nil, false
	}
	return result, true
}

func (s *SMSystem[T]) Iterate() []error {
	return s.doIterate(s.Update)
}

func (s *SMSystem[T]) doIterate(fn func(each T) (T, error)) []error {
	var errs = make([]error, 10)
	errs = append(errs, s.Map.Iterate(fn)...)
	return errs
}

func GetComponent[T ecstypes.Component](sm ecstypes.SystemManager, entity ecstypes.EntityID) (*T, error) {
	var zero T
	var sys ecstypes.System
	var err error
	if sys, err = sm.GetSystem(zero.SystemID()); err != nil {
		return nil, fmt.Errorf("error geting system: %w", err)
	}

	switch concreteSystem := sys.(type) {
	case *SMSystem[T]:
		c, ok := concreteSystem.GetComponent(entity)
		if !ok {
			return nil, nil
		}
		return c, nil
	default:
		return nil, fmt.Errorf("component storage: %w", ErrType)
	}
}
