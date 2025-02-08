package ecs

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/globals"
	"github.com/StCredZero/vectrek/sparse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

// NewEntity creates a new ship at the center of the screen
func NewEntity(x, y float64) *Entity {
	return &Entity{
		X:            x,
		Y:            y,
		ScreenWidth:  constants.ScreenWidth,
		ScreenHeight: constants.ScreenHeight,
		Angle:        0.0,
	}
}

type EntityID = sparse.Key

// Entity represents the player's spaceship with position, rotation, and movement
type Entity struct {
	// The following are fundamental to Entity
	ID EntityID // Stays in Entity
	X  float64  // X position on screen // Stays in Entity
	Y  float64  // Y position on screen // Stays in Entity

	// The following are to be moved to Motion Component
	XV    float64 // Motion Component
	YV    float64 // Motion Component
	Angle float64 // rotation Angle in radians // Motion component

	// The following are to be moved to Helm Component
	Input HelmInput // Helm component

	// The following are to be moved to Screen Component
	ScreenWidth  float64         // Screen Componen
	ScreenHeight float64         // Screen Componen
	Vertices     []ebiten.Vertex // Screen Componen
	Indices      []uint16        // Screen Componen
}

type HelmInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

// Entity thrust
const (
	ThrustAccel = 0.2
	maxVelocity = 5.0
)

func (s *Entity) Update(instance *Instance) error {
	input := s.Input

	// Entity rotation (3 degrees per tick)
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

func (s *Entity) Draw(screen *ebiten.Image, aa bool, line bool) {
	var path vector.Path
	fmt.Printf("ship draw\n")

	// Define ship as a triangle
	length := float32(15.0)
	theta := float32(s.Angle)

	// Front point
	path.MoveTo(
		float32(s.X)+length*float32(math.Cos(float64(theta))),
		float32(s.Y)+length*float32(math.Sin(float64(theta))),
	)

	// Right point (120 degrees from front)
	path.LineTo(
		float32(s.X)+length*float32(math.Cos(float64(theta)+2.0944)), // 2.0944 rad = 120 deg
		float32(s.Y)+length*float32(math.Sin(float64(theta)+2.0944)),
	)

	// Left point (-120 degrees from front)
	path.LineTo(
		float32(s.X)+length*float32(math.Cos(float64(theta)-2.0944)),
		float32(s.Y)+length*float32(math.Sin(float64(theta)-2.0944)),
	)

	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 2
		op.LineJoin = vector.LineJoinRound
		s.Vertices, s.Indices = path.AppendVerticesAndIndicesForStroke(s.Vertices[:0], s.Indices[:0], op)
	} else {
		s.Vertices, s.Indices = path.AppendVerticesAndIndicesForFilling(s.Vertices[:0], s.Indices[:0])
	}

	for i := range s.Vertices {
		s.Vertices[i].SrcX = 1
		s.Vertices[i].SrcY = 1
		s.Vertices[i].ColorR = 1
		s.Vertices[i].ColorG = 1
		s.Vertices[i].ColorB = 1
		s.Vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(s.Vertices, s.Indices, globals.WhiteSubImage, op)
}
