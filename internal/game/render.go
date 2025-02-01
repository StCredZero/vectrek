package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/gl"
	"math"
)

type RenderComponent struct {
	ecs.BasicComponent
	Points []float32 // x1,y1, x2,y2, etc.
	Color  struct{ R, G, B float32 }
}

// EntityType identifies the type of entity for rendering
type EntityType int

const (
	EntityTypeShip EntityType = iota
	EntityTypeRectangle
)

// CreateRenderComponent creates a RenderComponent based on entity type
func CreateRenderComponent(entityType EntityType) *RenderComponent {
	render := &RenderComponent{}
	render.Color.R = 1
	render.Color.G = 1
	render.Color.B = 1

	switch entityType {
	case EntityTypeShip:
		render.Points = []float32{
			0, -20, // nose
			15, 20,  // right
			-15, 20, // left
		}
	case EntityTypeRectangle:
		render.Points = []float32{
			-25, -25, // top left
			25, -25,  // top right
			25, 25,   // bottom right
			-25, 25,  // bottom left
		}
	}

	return render
}

func (s *RenderSystem) Update(dt float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.LoadIdentity()

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
		
		// Translate to position with screen wrapping
		for i := 0; i < len(transformedPoints); i += 2 {
			x := pos.X + transformedPoints[i]
			y := pos.Y + transformedPoints[i+1]

			// Handle screen wrapping
			if x < 0 {
				x += 800
			} else if x > 800 {
				x -= 800
			}
			if y < 0 {
				y += 600
			} else if y > 600 {
				y -= 600
			}

			transformedPoints[i] = x
			transformedPoints[i+1] = y
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
