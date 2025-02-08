package ecs

import (
	"errors"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/slices"
	"github.com/StCredZero/vectrek/sparse"
	"github.com/hajimehoshi/ebiten/v2"
)

type Parameters struct {
	ScreenWidth  float64
	ScreenHeight float64
}
type Instance struct {
	Entities map[EntityID]struct{}
	Motions  *sparse.Map[Motion]
	Helms    *sparse.Map[Helm]
	Sprites  *sparse.Map[Sprite]

	Parameters Parameters
}

func NewInstance(parameters Parameters) *Instance {
	return &Instance{
		Entities: make(map[EntityID]struct{}),
		Motions:  sparse.NewMap[Motion](),
		Helms:    sparse.NewMap[Helm](),

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
func (g *Instance) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
func (i *Instance) AddEntity(
	entity EntityID,
	motion *Motion,
	helm *Helm,
	sprite *Sprite,
) error {
	i.Entities[entity] = struct{}{}

	if motion != nil {
		motion.Entity = entity
		i.Motions.Add(entity, *motion)
		if err := motion.VerifyInit(); err != nil {
			return err
		}
	}
	if helm != nil {
		helm.Entity = entity
		helm.Motion = i.Motions.MustGet(entity)
		i.Helms.Add(entity, *helm)
		if err := helm.VerifyInit(); err != nil {
			return err
		}
	}
	if sprite != nil {
		sprite.Entity = entity
		sprite.Motion = i.Motions.MustGet(entity)
		i.Sprites.Add(entity, *sprite)
		if err := sprite.VerifyInit(); err != nil {
			return err
		}
	}
	return nil
}
