package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"math"
)

// Basic components
type PositionComponent struct {
	ecs.BasicComponent
	X, Y float32
}

type VelocityComponent struct {
	ecs.BasicComponent
	X, Y float32
}

type RotationComponent struct {
	ecs.BasicComponent
	Angle float32
}

// Ship control system
type ShipControlSystem struct {
	world *ecs.World
	ship  *ecs.Entity
}

func (s *ShipControlSystem) New(w *ecs.World) {
	s.world = w
	
	// Create ship entity
	ship := ecs.NewEntity("Ship")
	
	// Add components
	ship.AddComponent(&PositionComponent{X: 400, Y: 300})
	ship.AddComponent(&VelocityComponent{})
	ship.AddComponent(&RotationComponent{})
	
	s.ship = ship
	w.AddEntity(ship)
	
	// Create rectangle entity
	rect := ecs.NewEntity("Rectangle")
	rect.AddComponent(&PositionComponent{X: 200, Y: 200})
	w.AddEntity(rect)
}

func (s *ShipControlSystem) Update(dt float32) {
	if s.ship == nil {
		return
	}
	
	// Handle rotation
	if engo.Input.Button("ArrowLeft").Down() {
		rot := s.ship.GetComponent(&RotationComponent{}).(*RotationComponent)
		rot.Angle -= 4 * dt
	}
	if engo.Input.Button("ArrowRight").Down() {
		rot := s.ship.GetComponent(&RotationComponent{}).(*RotationComponent)
		rot.Angle += 4 * dt
	}
	
	// Handle thrust
	if engo.Input.Button("ArrowUp").Down() {
		rot := s.ship.GetComponent(&RotationComponent{}).(*RotationComponent)
		vel := s.ship.GetComponent(&VelocityComponent{}).(*VelocityComponent)
		
		// Apply thrust in direction of rotation
		thrust := float32(200)
		vel.X += thrust * float32(math.Cos(float64(rot.Angle))) * dt
		vel.Y += thrust * float32(math.Sin(float64(rot.Angle))) * dt
	}
	
	// Update position based on velocity
	pos := s.ship.GetComponent(&PositionComponent{}).(*PositionComponent)
	vel := s.ship.GetComponent(&VelocityComponent{}).(*VelocityComponent)
	
	pos.X += vel.X * dt
	pos.Y += vel.Y * dt
	
	// Wrap around screen edges (toroidal topology)
	if pos.X < 0 {
		pos.X = 800
	} else if pos.X > 800 {
		pos.X = 0
	}
	if pos.Y < 0 {
		pos.Y = 600
	} else if pos.Y > 600 {
		pos.Y = 0
	}
}

func (s *ShipControlSystem) Remove(e ecs.BasicEntity) {}

// Render system
type RenderSystem struct {
	world *ecs.World
}

func (s *RenderSystem) New(w *ecs.World) {
	s.world = w
	
	// Register input controls
	engo.Input.RegisterButton("ArrowLeft", engo.KeyArrowLeft)
	engo.Input.RegisterButton("ArrowRight", engo.KeyArrowRight)
	engo.Input.RegisterButton("ArrowUp", engo.KeyArrowUp)
}

func (s *RenderSystem) Update(dt float32) {
	// Rendering will be implemented here
}

func (s *RenderSystem) Remove(e ecs.BasicEntity) {}
