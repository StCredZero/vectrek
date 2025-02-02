package ship

import (
	"fmt"
	"github.com/StCredZero/vectrek/globals"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	Input        PilotInput

	Vertices []ebiten.Vertex
	Indices  []uint16
}

type PilotInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

// Ship thrust
const (
	ThrustAccel = 0.2
	maxVelocity = 5.0
)

func (s *Ship) Update() error {
	input := s.Input

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

func (s *Ship) Draw(screen *ebiten.Image, aa bool, line bool) {
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
