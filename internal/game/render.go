package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
	"math"
)

type RenderComponent struct {
	ecs.BasicComponent
	Points []float32 // x1,y1, x2,y2, etc.
	Color  struct{ R, G, B float32 }
}

func (s *RenderSystem) Update(dt float32) {
	for _, e := range s.world.Entities() {
		pos, hasPos := e.GetComponent(&PositionComponent{}).(*PositionComponent)
		render, hasRender := e.GetComponent(&RenderComponent{}).(*RenderComponent)
		
		if !hasPos || !hasRender {
			continue
		}
		
		// Transform points based on position
		transformedPoints := make([]float32, len(render.Points))
		copy(transformedPoints, render.Points)
		
		// Apply rotation if entity has it
		if rot, hasRot := e.GetComponent(&RotationComponent{}).(*RotationComponent); hasRot {
			sin := float32(math.Sin(float64(rot.Angle)))
			cos := float32(math.Cos(float64(rot.Angle)))
			
			for i := 0; i < len(transformedPoints); i += 2 {
				x := transformedPoints[i]
				y := transformedPoints[i+1]
				transformedPoints[i] = x*cos - y*sin
				transformedPoints[i+1] = x*sin + y*cos
			}
		}
		
		// Translate to position
		for i := 0; i < len(transformedPoints); i += 2 {
			transformedPoints[i] += pos.X
			transformedPoints[i+1] += pos.Y
		}
		
		// Draw the shape
		gl.Begin(gl.LINE_LOOP)
		gl.Color3f(render.Color.R, render.Color.G, render.Color.B)
		
		for i := 0; i < len(transformedPoints); i += 2 {
			gl.Vertex2f(transformedPoints[i], transformedPoints[i+1])
		}
		
		gl.End()
	}
}

func (s *ShipControlSystem) New(w *ecs.World) {
	s.world = w
	
	// Create ship entity with triangle shape
	ship := ecs.NewEntity("Ship")
	shipRender := &RenderComponent{
		Points: []float32{
			0, -20, // nose
			15, 20, // right
			-15, 20, // left
		},
	}
	shipRender.Color.R = 1
	shipRender.Color.G = 1
	shipRender.Color.B = 1
	
	ship.AddComponent(&PositionComponent{X: 400, Y: 300})
	ship.AddComponent(&VelocityComponent{})
	ship.AddComponent(&RotationComponent{})
	ship.AddComponent(shipRender)
	
	s.ship = ship
	w.AddEntity(ship)
	
	// Create rectangle entity
	rect := ecs.NewEntity("Rectangle")
	rectRender := &RenderComponent{
		Points: []float32{
			-25, -25, // top left
			25, -25,  // top right
			25, 25,   // bottom right
			-25, 25,  // bottom left
		},
	}
	rectRender.Color.R = 1
	rectRender.Color.G = 1
	rectRender.Color.B = 1
	
	rect.AddComponent(&PositionComponent{X: 200, Y: 200})
	rect.AddComponent(rectRender)
	
	w.AddEntity(rect)
}
