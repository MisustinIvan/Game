package main

type Camera struct {
	rect Rect
	tgt  *Vector2
}

func NewCamera(size Vector2, tgt *Vector2) Camera {
	return Camera{
		rect: Rect{
			Vector2{tgt.x - size.x/2,
				tgt.y - size.y/2,
			},
			size,
		},
		tgt: tgt,
	}
}

func (c *Camera) Update() {
	c.rect.pos.x = c.tgt.x - c.rect.extents.x/2
	c.rect.pos.y = c.tgt.y - c.rect.extents.y/2
}
