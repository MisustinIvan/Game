package main

type Camera struct {
	rect Rect
}

func NewCamera(Rect) Camera {
	return Camera{
		rect: Rect{},
	}
}
