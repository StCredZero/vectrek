package ecs

import (
	"errors"
	"github.com/StCredZero/vectrek/slices"
	"reflect"
)

var ErrMissingPrerequisite = errors.New("missing prerequisite")

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	// Only call IsNil() on types that support it
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func HasPrerequisites(objs ...any) error {
	if slices.Detect(objs, func(obj any) bool {
		return isNil(obj)
	}) {
		return ErrMissingPrerequisite
	}
	return nil
}
