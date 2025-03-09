package ecs

import (
	"errors"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/slices"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
	"time"
)

type Parameters struct {
	ScreenWidth  float64
	ScreenHeight float64
}
type Instance struct {
	Entities map[ecstypes.EntityID]struct{}

	Name string

	Position     *SMSystem[Position]
	Motion       *SMSystem[Motion]
	Helm         *SMSystem[Helm]
	Sprite       *SMSystem[Sprite]
	Player       *SMSystem[Player]
	SyncReceiver *SMSystem[SyncReceiver]
	SyncSender   *SMSystem[SyncSender]

	Counter    uint64
	Parameters Parameters

	Pipe     *Pipe
	Receiver ecstypes.Receiver
	Sender   ecstypes.Sender
}

func NewInstance(parameters Parameters) *Instance {
	var result = new(Instance)
	result.Entities = make(map[ecstypes.EntityID]struct{})
	result.Position = NewSMSystem[Position](func(each Position) (Position, error) {
		return each.Update(result)
	})
	result.Motion = NewSMSystem[Motion](func(each Motion) (Motion, error) {
		return each.Update(result)
	})
	result.Helm = NewSMSystem[Helm](func(each Helm) (Helm, error) {
		return each.Update(result)
	})
	result.Sprite = NewSMSystem[Sprite](func(each Sprite) (Sprite, error) {
		return each.Update(result)
	})
	result.Player = NewSMSystem[Player](func(each Player) (Player, error) {
		return each.Update(result)
	})
	result.SyncReceiver = NewSMSystem[SyncReceiver](func(each SyncReceiver) (SyncReceiver, error) {
		return each.Update(result)
	})
	result.SyncSender = NewSMSystem[SyncSender](func(each SyncSender) (SyncSender, error) { return each.Update(result) })
	result.Parameters = parameters
	return result
}
func (i *Instance) GetSender() ecstypes.Sender {
	return i.Sender
}
func (i *Instance) SetSender(sender ecstypes.Sender) {
	i.Sender = sender
}
func (i *Instance) GetReceiver() ecstypes.Receiver {
	return i.Receiver
}
func (i *Instance) SetReceiver(pipe ecstypes.Receiver) {
	i.Receiver = pipe
}
func (i *Instance) GetName() string {
	return i.Name
}
func (i *Instance) RunServer(done chan bool) {
	ticker := time.NewTicker(16667 * time.Microsecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i.Update()
		case <-done:
			return
		}
	}
}
func (i *Instance) GetCounter() uint64 {
	return i.Counter
}
func (i *Instance) Update() error {
	i.Counter++

	var hasMessage bool
	var msg ecstypes.ComponentMessage
	for {
		if msg, hasMessage = i.Receiver.Receive(); !hasMessage {
			break
		}
		switch obj := msg.Payload.(type) {
		case HelmInput:
			if helm, ok := i.Helm.GetComponent(msg.Entity); ok {
				helm.Input = obj
			}
		case SyncInput:
			if sync, ok := i.SyncReceiver.GetComponent(msg.Entity); ok {
				sync.Input <- obj
			}
		default:
		}
	}

	// systems must be executed in reverse dependency order
	var errs []error
	errs = append(errs, i.Helm.Iterate()...)
	errs = append(errs, i.Motion.Iterate()...)
	//errs = append(errs, i.Sprite.Iterate()...)
	errs = append(errs, i.Player.Iterate()...)
	errs = append(errs, i.SyncSender.Iterate()...)
	errs = append(errs, i.SyncReceiver.Iterate()...)

	errs = slices.Select(errs, func(err error) bool {
		return err != nil
	})
	return errors.Join(errs...)
}
func (i *Instance) Draw(screen *ebiten.Image) {
	i.Sprite.doIterate(func(sprite Sprite) (Sprite, error) {
		sprite.Draw(screen, false, false)
		return sprite, nil
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
