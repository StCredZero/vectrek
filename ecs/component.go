package ecs

import (
	"fmt"
	"github.com/StCredZero/vectrek/geom"
	"github.com/StCredZero/vectrek/globals"
	"github.com/StCredZero/vectrek/vterr"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

type Component interface {
	Init(i *Instance, entity EntityID) error
	Update(i *Instance) error
	SystemID() SystemID
}

type Position struct {
	Entity EntityID
	geom.Vector
	geom.Angle
}

func (m *Position) Init(i *Instance, entity EntityID) error {
	m.Entity = entity
	i.Positions.Add(m.Entity, *m)
	return nil
}
func (m *Position) Update(_ *Instance) error {
	return nil
}
func (m *Position) SystemID() SystemID {
	return SystemPosition
}

type Motion struct {
	Entity   EntityID
	Position *Position
	Velocity geom.Vector
}

func (m *Motion) Init(i *Instance, entity EntityID) error {
	m.Entity = entity
	if p, gotPosition := i.Positions.Get(m.Entity); gotPosition {
		m.Position = p
	} else {
		return fmt.Errorf("no position found for %s: %w", m.Entity, vterr.ErrMissing)
	}
	i.Motions.Add(m.Entity, *m)
	return nil
}
func (m *Motion) Update(i *Instance) error {
	m.Position.Vector = m.Position.Vector.Add(m.Velocity)
	fmt.Printf("position update %f %f %f %f \n", m.Velocity.X, m.Velocity.Y, m.Position.X, m.Position.Y)
	// Wrap around screen edges (toroidal topology)
	if m.Position.X < 0 {
		m.Position.X += i.Parameters.ScreenWidth
	} else if m.Position.X >= i.Parameters.ScreenWidth {
		m.Position.X -= i.Parameters.ScreenWidth
	}
	if m.Position.Y < 0 {
		m.Position.Y += i.Parameters.ScreenHeight
	} else if m.Position.Y >= i.Parameters.ScreenHeight {
		m.Position.Y -= i.Parameters.ScreenHeight
	}
	return nil
}
func (m *Motion) SystemID() SystemID {
	return SystemMotion
}

type Helm struct {
	Entity   EntityID
	Position *Position
	Motion   *Motion
	Input    HelmInput
}

func (m *Helm) Init(i *Instance, entity EntityID) error {
	m.Entity = entity
	if m.Motion = i.Motions.MustGet(m.Entity); m.Motion == nil {
		return fmt.Errorf("no Motion found: %w", vterr.ErrMissing)
	}
	if m.Position = i.Positions.MustGet(m.Entity); m.Position == nil {
		return fmt.Errorf("no Position found: %w", vterr.ErrMissing)
	}
	i.Helms.Add(m.Entity, *m)
	return nil
}
func (m *Helm) Update(i *Instance) error {
	var input HelmInput
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		input.Left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		input.Right = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		input.Thrust = true
	}
	// Entity rotation (3 degrees per tick)
	if input.Left {
		m.Position.Angle -= 3 * (math.Pi / 180)
	}
	if input.Right {
		m.Position.Angle += 3 * (math.Pi / 180)
	}
	if input.Thrust {
		// Update velocity based on velocity and angle
		m.Motion.Velocity = m.Motion.Velocity.Add(m.Position.Angle.ToVector().Multiply(ThrustAccel))
		fmt.Println("accel velocity %f %f \n", m.Motion.Velocity.X, m.Motion.Velocity.Y)
	}
	return nil
}
func (m *Helm) SystemID() SystemID {
	return SystemHelm
}

type Sprite struct {
	Entity   EntityID
	Motion   *Motion
	Position *Position
	Vertices []ebiten.Vertex
	Indices  []uint16
}

func (m *Sprite) Init(i *Instance, entity EntityID) error {
	m.Entity = entity
	if m.Motion = i.Motions.MustGet(m.Entity); m.Motion == nil {
		return fmt.Errorf("no Motion found: %w", vterr.ErrMissing)
	}
	if m.Position = i.Positions.MustGet(m.Entity); m.Position == nil {
		return fmt.Errorf("no Position found: %w", vterr.ErrMissing)
	}
	i.Sprites.Add(m.Entity, *m)
	return nil
}
func (m *Sprite) Update(i *Instance) error {
	return nil
}
func (s *Sprite) Draw(screen *ebiten.Image, aa bool, line bool) {
	var path vector.Path

	// Define ship as a triangle
	length := float32(15.0)
	theta := float32(s.Position.Angle)

	// Front point
	path.MoveTo(
		float32(s.Motion.Position.X)+length*float32(math.Cos(float64(theta))),
		float32(s.Motion.Position.Y)+length*float32(math.Sin(float64(theta))),
	)

	// Right point (120 degrees from front)
	path.LineTo(
		float32(s.Motion.Position.X)+length*float32(math.Cos(float64(theta)+2.0944)), // 2.0944 rad = 120 deg
		float32(s.Motion.Position.Y)+length*float32(math.Sin(float64(theta)+2.0944)),
	)

	// Left point (-120 degrees from front)
	path.LineTo(
		float32(s.Motion.Position.X)+length*float32(math.Cos(float64(theta)-2.0944)),
		float32(s.Motion.Position.Y)+length*float32(math.Sin(float64(theta)-2.0944)),
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
func (m *Sprite) SystemID() SystemID {
	return SystemSprite
}
