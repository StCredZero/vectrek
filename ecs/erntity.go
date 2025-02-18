package ecs

// Entity represents the player's spaceship with position, rotation, and movement

type HelmInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

// Entity thrust
const (
	ThrustAccel = 0.2
	MaxVelocity = 5.0
)
