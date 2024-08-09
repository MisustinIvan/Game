package main

type Rect struct {
	pos     Vector2
	extents Vector2
}

func NewRect(pos Vector2, extents Vector2) Rect {
	return Rect{
		pos:     pos,
		extents: extents,
	}
}

func (r Rect) Intersects(o Rect) bool {
	return !(r.pos.x+r.extents.x <= o.pos.x ||
		r.pos.x >= o.pos.x+o.extents.x ||
		r.pos.y+r.extents.y <= o.pos.y ||
		r.pos.y >= o.pos.y+o.extents.y)
}

func (r Rect) Contains(o Rect) bool {
	return (r.pos.x <= o.pos.x && o.pos.x+o.extents.x <= r.pos.x+r.extents.x &&
		r.pos.y <= o.pos.y && o.pos.y+o.extents.y <= r.pos.y+r.extents.y)
}
