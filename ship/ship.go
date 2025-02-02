package ship

import (
	"math"
)

// NewShip creates a new ship at the center of the screen
func NewShip(screenWidth, screenHeight float64) *Ship {
	return &Ship{
		X:            screenWidth / 2,
		Y:            screenHeight / 2,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		Angle:        0.0,
	}
}

// Ship represents the player's spaceship with position, rotation, and movement
type Ship struct {
	X            float64 // X position on screen
	Y            float64 // Y position on screen
	XV           float64
	YV           float64
	Angle        float64 // rotation Angle in radians
	ScreenWidth  float64
	ScreenHeight float64
}

type ShipInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

// Ship thrust
const (
	ThrustAccel = 0.2
	maxVelocity = 5.0
)

func (s *Ship) Update(input ShipInput) error {
	// Ship rotation (3 degrees per tick)
	if input.Left {
		s.Angle -= 3 * (math.Pi / 180)
	}
	if input.Right {
		s.Angle += 3 * (math.Pi / 180)
	}
	if input.Thrust {
		// Update ship position based on velocity and angle
		s.XV += math.Cos(s.Angle) * ThrustAccel
		s.YV += math.Sin(s.Angle) * ThrustAccel
	}
	s.X += s.XV
	s.Y += s.YV
	// Wrap around screen edges (toroidal topology)
	if s.X < 0 {
		s.X += s.ScreenWidth
	} else if s.X >= s.ScreenWidth {
		s.X -= s.ScreenWidth
	}
	if s.Y < 0 {
		s.Y += s.ScreenHeight
	} else if s.Y >= s.ScreenHeight {
		s.Y -= s.ScreenHeight
	}
	return nil
}
