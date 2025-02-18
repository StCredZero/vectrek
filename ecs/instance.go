package ecs

import (
	"errors"
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/slices"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

type Parameters struct {
	ScreenWidth  float64
	ScreenHeight float64
}
type Instance struct {
	Entities map[ecstypes.EntityID]struct{}

	Position *SMSystem[Position]
	Motion   *SMSystem[Motion]
	Helm     *SMSystem[Helm]
	Sprite   *SMSystem[Sprite]

	Counter    uint64
	Parameters Parameters

	Pipe     *Pipe
	receiver Receiver
}

func NewInstance(parameters Parameters) *Instance {
	result := &Instance{
		Entities: make(map[ecstypes.EntityID]struct{}),

		Position: NewSMSystem[Position](func(each *Position) error {
			return each.Update()
		}),
		Motion: NewSMSystem[Motion](func(each *Motion) error {
			return each.Update()
		}),
		Helm: NewSMSystem[Helm](func(each *Helm) error {
			return each.Update()
		}),

		// sprites are like a system, but they are
		// executed by Draw
		Sprite: NewSMSystem[Sprite](func(each *Sprite) error {
			return each.Update()
		}),

		Parameters: parameters,
	}
	return result
}
func (i *Instance) GetSystem(id ecstypes.SystemID) (ecstypes.System, error) {
	switch id {
	case ecstypes.SystemPosition:
		return i.Position, nil
	case ecstypes.SystemMotion:
		return i.Motion, nil
	case ecstypes.SystemHelm:
		return i.Helm, nil
	default:
		return nil, fmt.Errorf("invalid system id: %w", ErrType)
	}
}
func (i *Instance) AddComponent(e ecstypes.EntityID, component ecstypes.Component) error {
	switch c := component.(type) {
	case *Position:
		if err := i.Position.AddComponent(e, c); err != nil {
			return err
		}
	case *Motion:
		if err := i.Motion.AddComponent(e, c); err != nil {
			return err
		}
	case *Helm:
		if err := i.Helm.AddComponent(e, c); err != nil {
			return err
		}
	case *Sprite:
		if err := i.Sprite.AddComponent(e, c); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid system type %v: %w", component, ErrType)
	}
	return nil
}
func (i *Instance) GetComponent(systemID ecstypes.SystemID, e ecstypes.EntityID) (ecstypes.Component, bool) {
	switch systemID {
	case ecstypes.SystemPosition:
		return i.Position.GetComponent(e)
	case ecstypes.SystemMotion:
		return i.Motion.GetComponent(e)
	case ecstypes.SystemHelm:
		return i.Helm.GetComponent(e)
	case ecstypes.SystemSprite:
		return i.Sprite.GetComponent(e)
	default:
		return nil, false
	}
}
func (i *Instance) SetPipe(pipe *Pipe) {
	i.Pipe = pipe
	i.receiver = pipe
}
func (i *Instance) Update() error {
	i.Counter++

	var hasMessage bool
	var msg ComponentMessage
	for {
		if msg, hasMessage = i.receiver.Receive(); !hasMessage {
			break
		}
		switch obj := msg.Payload.(type) {
		case HelmInput:
			helm, ok := i.Helm.GetComponent(msg.Entity)
			if ok {
				helm.Input = obj
			}
		default:
		}
	}

	// systems must be executed in reverse dependency order
	var errs []error
	errs = append(errs, i.Helm.iterate(func(each *Helm) error {
		return each.Update()
	})...)
	errs = append(errs, i.Motion.iterate(func(each *Motion) error {
		return each.Update()
	})...)

	errs = slices.Select(errs, func(err error) bool {
		return err != nil
	})
	return errors.Join(errs...)
}
func (i *Instance) Draw(screen *ebiten.Image) {
	i.Sprite.iterate(func(sprite *Sprite) error {
		sprite.Draw(screen, false, false)
		return nil
	})
}
func (i *Instance) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
func (i *Instance) AddEntity(
	entity ecstypes.EntityID,
	components ...ecstypes.Component,
) error {
	i.Entities[entity] = struct{}{}
	sort.Slice(components, func(i, j int) bool {
		return components[i].SystemID() < components[j].SystemID()
	})
	for _, component := range components {
		if err := component.Init(i, entity); err != nil {
			return err
		}
	}
	return nil
}
