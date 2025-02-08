package ecs

import (
	"errors"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/slices"
	"github.com/StCredZero/vectrek/sparse"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

type Parameters struct {
	ScreenWidth  float64
	ScreenHeight float64
}
type Instance struct {
	Entities  map[EntityID]struct{}
	Positions *sparse.Map[Position]
	Motions   *sparse.Map[Motion]
	Helms     *sparse.Map[Helm]
	Sprites   *sparse.Map[Sprite]

	Parameters Parameters
}

func NewInstance(parameters Parameters) *Instance {
	return &Instance{
		Entities:  make(map[EntityID]struct{}),
		Positions: sparse.NewMap[Position](),
		Motions:   sparse.NewMap[Motion](),
		Helms:     sparse.NewMap[Helm](),

		// sprites are like a system, but they are
		// executed by Draw
		Sprites: sparse.NewMap[Sprite](),

		Parameters: parameters,
	}
}
func (i *Instance) Update() error {
	// systems must be executed in reverse dependency order
	var errs []error
	errs = append(errs, i.Helms.Iterate(func(each *Helm) error {
		return each.Update(i)
	})...)
	errs = append(errs, i.Motions.Iterate(func(each *Motion) error {
		return each.Update(i)
	})...)

	errs = slices.Select(errs, func(err error) bool {
		return err != nil
	})
	return errors.Join(errs...)
}
func (i *Instance) Draw(screen *ebiten.Image) {
	i.Sprites.Iterate(func(sprite *Sprite) error {
		sprite.Draw(screen, false, false)
		return nil
	})
}
func (i *Instance) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
func (i *Instance) AddEntity(
	entity EntityID,
	components ...Component,
) error {
	i.Entities[entity] = struct{}{}
	sort.Slice(components, func(i, j int) bool {
		return int(components[i].SystemID()) < components[j].SystemID()
	})
	for _, component := range components {
		if err := component.Init(i, entity); err != nil {
			return err
		}
	}
	return nil
}
