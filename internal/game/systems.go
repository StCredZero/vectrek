package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
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

// Ship control system (server-side)
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

func (s *ShipControlSystem) HandleInput(input struct{ Left, Right, Up bool }, dt float32) {
	if s.ship == nil {
		return
	}

	// Get components once
	rot := s.ship.GetComponent(&RotationComponent{}).(*RotationComponent)
	vel := s.ship.GetComponent(&VelocityComponent{}).(*VelocityComponent)

	// Server-side input validation and handling
	maxRotationSpeed := float32(4.0)
	if input.Left {
		rot.Angle -= maxRotationSpeed * dt
	}
	if input.Right {
		rot.Angle += maxRotationSpeed * dt
	}
	
	// Server-side thrust validation and handling
	maxThrust := float32(200.0)
	maxSpeed := float32(400.0)
	if input.Up {
		// Apply thrust in direction of rotation with speed limit
		thrust := maxThrust * dt
		vel.X += thrust * float32(math.Cos(float64(rot.Angle)))
		vel.Y += thrust * float32(math.Sin(float64(rot.Angle)))

		// Enforce maximum velocity (server-side authority)
		speed := float32(math.Sqrt(float64(vel.X*vel.X + vel.Y*vel.Y)))
		if speed > maxSpeed {
			ratio := maxSpeed / speed
			vel.X *= ratio
			vel.Y *= ratio
		}
	}
}

func (s *ShipControlSystem) Remove(e ecs.BasicEntity) {}

// Render system (client-side)
type RenderSystem struct {
	world *ecs.World
}

func (s *RenderSystem) New(w *ecs.World) {
	s.world = w
	
	// Set up OpenGL
	gl.ClearColor(0, 0, 0, 1)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, 800, 600, 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
}

func (s *RenderSystem) Remove(e ecs.BasicEntity) {}
