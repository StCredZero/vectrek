package ecs

import (
	"fmt"
	"github.com/StCredZero/vectrek/ecstypes"
)

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
	case ecstypes.SystemPlayer:
		return i.Player.GetComponent(e)
	case ecstypes.SystemSyncReceiver:
		return i.SyncReceiver.GetComponent(e)
	case ecstypes.SystemSyncSender:
		return i.SyncSender.GetComponent(e)
	default:
		return nil, false
	}
}
func (i *Instance) GetSystem(id ecstypes.SystemID) (ecstypes.System, error) {
	switch id {
	case ecstypes.SystemPosition:
		return i.Position, nil
	case ecstypes.SystemMotion:
		return i.Motion, nil
	case ecstypes.SystemHelm:
		return i.Helm, nil
	case ecstypes.SystemPlayer:
		return i.Player, nil
	case ecstypes.SystemSyncReceiver:
		return i.SyncReceiver, nil
	case ecstypes.SystemSyncSender:
		return i.SyncSender, nil
	default:
		return nil, fmt.Errorf("invalid system id: %w", ErrType)
	}
}
func (i *Instance) AddComponent(e ecstypes.EntityID, component ecstypes.Component) error {
	switch c := component.(type) {
	case Position:
		if err := i.Position.AddComponent(e, c); err != nil {
			return err
		}
	case Motion:
		if err := i.Motion.AddComponent(e, c); err != nil {
			return err
		}
	case Helm:
		if err := i.Helm.AddComponent(e, c); err != nil {
			return err
		}
	case Sprite:
		if err := i.Sprite.AddComponent(e, c); err != nil {
			return err
		}
	case Player:
		if err := i.Player.AddComponent(e, c); err != nil {
			return err
		}
	case SyncReceiver:
		if err := i.SyncReceiver.AddComponent(e, c); err != nil {
			return err
		}
	case SyncSender:
		if err := i.SyncSender.AddComponent(e, c); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid system type %v: %w", component, ErrType)
	}
	return nil
}
