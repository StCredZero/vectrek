package ecs

import "github.com/StCredZero/vectrek/geom"

// Entity represents the player's spaceship with position, rotation, and movement

type HelmInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

type SyncInput struct {
	Position geom.Vector
	Velocity geom.Vector
	Angle    geom.Angle
}

// Entity thrust
const (
	ThrustAccel = 0.2
	MaxVelocity = 5.0
)
