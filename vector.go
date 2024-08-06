package main

import "math"

type Vector2 struct {
	x float64
	y float64
}

func (v Vector2) Add(o Vector2) Vector2 {
	return Vector2{
		v.x + o.x,
		v.y + o.y,
	}
}

func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{
		v.x - o.x,
		v.y - o.y,
	}
}

func (v Vector2) Scale(a float64) Vector2 {
	return Vector2{
		v.x * a,
		v.y * a,
	}
}

func (v Vector2) Mag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v Vector2) Norm() Vector2 {
	mag := v.Mag()
	return Vector2{
		v.x / mag,
		v.y / mag,
	}
}

// do i really need this??? we'l se in the future
func (v Vector2) Dot(o Vector2) float64 {
	return v.x*o.x + v.y*o.y
}
