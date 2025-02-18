package ecs

import (
	"errors"
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/geom"
	"github.com/StCredZero/vectrek/globals"
	"github.com/StCredZero/vectrek/vterr"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

var ErrType = errors.New("type error")

type Position struct {
	Entity ecstypes.EntityID
	geom.Vector
	geom.Angle
}

func (comp *Position) Init(sm ecstypes.SystemManager, entity ecstypes.EntityID) error {
	comp.Entity = entity
	if err := sm.AddComponent(entity, comp); err != nil {
		return fmt.Errorf("adding position: %w", err)
	}
	return nil
}
func (comp *Position) Update() error {
	return nil
}
func (comp Position) SystemID() ecstypes.SystemID {
	return ecstypes.SystemPosition
}

type Motion struct {
	Entity   ecstypes.EntityID
	Position *Position
	Velocity geom.Vector
}

func (comp *Motion) Init(sm ecstypes.SystemManager, entity ecstypes.EntityID) error {
	var err error
	comp.Entity = entity
	if comp.Position, err = GetComponent[Position](sm, entity); err != nil {
		return fmt.Errorf("adding position: %w", err)
	}
	if err = sm.AddComponent(entity, comp); err != nil {
		return fmt.Errorf("adding motion: %w", err)
	}
	return nil
}
func (comp *Motion) Update() error {
	comp.Position.Vector = comp.Position.Vector.Add(comp.Velocity)

	// Wrap around screen edges (toroidal topology)
	if comp.Position.X < 0 {
		comp.Position.X += constants.ScreenWidth
	} else if comp.Position.X >= constants.ScreenWidth {
		comp.Position.X -= constants.ScreenWidth
	}
	if comp.Position.Y < 0 {
		comp.Position.Y += constants.ScreenHeight
	} else if comp.Position.Y >= constants.ScreenHeight {
		comp.Position.Y -= constants.ScreenHeight
	}
	return nil
}
func (comp Motion) SystemID() ecstypes.SystemID {
	return ecstypes.SystemMotion
}

type Helm struct {
	Entity   ecstypes.EntityID
	Position *Position
	Motion   *Motion
	Input    HelmInput
}

func (comp *Helm) Init(sm ecstypes.SystemManager, entity ecstypes.EntityID) error {
	var err error
	comp.Entity = entity
	if comp.Motion, err = GetComponent[Motion](sm, entity); comp.Motion == nil {
		return fmt.Errorf("no Motion found: %w", vterr.ErrMissing)
	}
	if comp.Position, err = GetComponent[Position](sm, entity); comp.Position == nil {
		return fmt.Errorf("no Position found: %w", vterr.ErrMissing)
	}
	if err = sm.AddComponent(entity, comp); err != nil {
		return fmt.Errorf("adding motion: %w", err)
	}
	return nil
}
func (comp *Helm) Update() error {
	input := comp.Input
	if input.Left {
		comp.Position.Angle -= 3 * (math.Pi / 180)
	}
	if input.Right {
		comp.Position.Angle += 3 * (math.Pi / 180)
	}
	if input.Thrust {
		// Update velocity based on velocity and angle
		comp.Motion.Velocity = comp.Motion.Velocity.Add(comp.Position.Angle.ToVector().Multiply(ThrustAccel))
		fmt.Printf("accel velocity %f %f \n", comp.Motion.Velocity.X, comp.Motion.Velocity.Y)
	}
	return nil
}
func (comp Helm) SystemID() ecstypes.SystemID {
	return ecstypes.SystemHelm
}

type Sprite struct {
	Entity   ecstypes.EntityID
	Motion   *Motion
	Position *Position
	Vertices []ebiten.Vertex
	Indices  []uint16
}

func (comp *Sprite) Init(sm ecstypes.SystemManager, entity ecstypes.EntityID) error {
	var err error
	comp.Entity = entity
	if comp.Motion, err = GetComponent[Motion](sm, entity); comp.Motion == nil {
		return fmt.Errorf("no Motion found: %w", vterr.ErrMissing)
	}
	if comp.Position, err = GetComponent[Position](sm, entity); comp.Position == nil {
		return fmt.Errorf("no Position found: %w", vterr.ErrMissing)
	}
	if err = sm.AddComponent(entity, comp); err != nil {
		return fmt.Errorf("adding motion: %w", err)
	}
	return nil
}
func (comp *Sprite) Update() error {
	return nil
}
func (comp *Sprite) Draw(screen *ebiten.Image, aa bool, line bool) {
	var path vector.Path

	// Define ship as a triangle
	length := float32(15.0)
	theta := float32(comp.Position.Angle)

	// Front point
	path.MoveTo(
		float32(comp.Motion.Position.X)+length*float32(math.Cos(float64(theta))),
		float32(comp.Motion.Position.Y)+length*float32(math.Sin(float64(theta))),
	)

	// Right point (120 degrees from front)
	path.LineTo(
		float32(comp.Motion.Position.X)+length*float32(math.Cos(float64(theta)+2.0944)), // 2.0944 rad = 120 deg
		float32(comp.Motion.Position.Y)+length*float32(math.Sin(float64(theta)+2.0944)),
	)

	// Left point (-120 degrees from front)
	path.LineTo(
		float32(comp.Motion.Position.X)+length*float32(math.Cos(float64(theta)-2.0944)),
		float32(comp.Motion.Position.Y)+length*float32(math.Sin(float64(theta)-2.0944)),
	)

	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 2
		op.LineJoin = vector.LineJoinRound
		comp.Vertices, comp.Indices = path.AppendVerticesAndIndicesForStroke(comp.Vertices[:0], comp.Indices[:0], op)
	} else {
		comp.Vertices, comp.Indices = path.AppendVerticesAndIndicesForFilling(comp.Vertices[:0], comp.Indices[:0])
	}

	for i := range comp.Vertices {
		comp.Vertices[i].SrcX = 1
		comp.Vertices[i].SrcY = 1
		comp.Vertices[i].ColorR = 1
		comp.Vertices[i].ColorG = 1
		comp.Vertices[i].ColorB = 1
		comp.Vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(comp.Vertices, comp.Indices, globals.WhiteSubImage, op)
}
func (comp Sprite) SystemID() ecstypes.SystemID {
	return ecstypes.SystemSprite
}
