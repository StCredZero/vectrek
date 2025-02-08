package geom

import "math"

type Angle float64

func (angle Angle) ToVector() Vector {
	return Vector{
		X: math.Cos(float64(angle)),
		Y: math.Sin(float64(angle)),
	}
}

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Add(ov Vector) Vector {
	return Vector{
		X: v.X + ov.X,
		Y: v.Y + ov.Y,
	}
}
func (v Vector) Multiply(w float64) Vector {
	return Vector{
		X: v.X * w,
		Y: v.Y * w,
	}
}
