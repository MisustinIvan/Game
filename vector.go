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

func (v *Vector2) AddEq(o Vector2) {
	v.x += o.x
	v.y += o.y
}

func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{
		v.x - o.x,
		v.y - o.y,
	}
}

func (v *Vector2) SubEq(o Vector2) {
	v.x -= o.x
	v.y -= o.y
}

func (v Vector2) Scale(s float64) Vector2 {
	return Vector2{
		v.x * s,
		v.y * s,
	}
}

func (v *Vector2) ScaleEq(s float64) {
	v.x *= s
	v.y *= s
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

func (v *Vector2) NormEq() {
	mag := v.Mag()
	v.x /= mag
	v.y /= mag
}
