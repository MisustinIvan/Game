package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Entity struct {
	id   int
	rect Rect
}

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
	return !(r.pos.x+r.extents.x < o.pos.x ||
		r.pos.x > o.pos.x+o.extents.x ||
		r.pos.y+r.extents.y < o.pos.y ||
		r.pos.y > o.pos.y+o.extents.y)
}

func (r Rect) Contains(o Rect) bool {
	return (r.pos.x <= o.pos.x && o.pos.x+o.extents.x <= r.pos.x+r.extents.x &&
		r.pos.y <= o.pos.y && o.pos.y+o.extents.y <= r.pos.y+r.extents.y)
}

type QNodeStatic struct {
	rect     Rect
	children [4]*QNodeStatic
	capacity int
	values   []Entity
	leaf     bool
}

func NewStaticNode(rect Rect, capacity int) *QNodeStatic {
	return &QNodeStatic{
		rect:     rect,
		children: [4]*QNodeStatic{},
		capacity: capacity,
		values:   []Entity{},
		leaf:     true,
	}
}

func (n *QNodeStatic) Insert(e Entity) bool {
	// if entity does not fit, abort
	if !n.rect.Contains(e.rect) {
		return false
	}

	// insert if node not full leaf
	if n.leaf {
		if len(n.values) < n.capacity {
			n.values = append(n.values, e)
			return true
		} else {
			n.Subdivide() // if full subdivide
		}
	}

	// try to insert to the children
	for _, child := range n.children {
		if child.Insert(e) {
			return true
		}
	}

	// if entity does not fit into any children, put it into this branch
	n.values = append(n.values, e)
	return true
}

func (n *QNodeStatic) Subdivide() {
	hwidth := n.rect.extents.x / 2
	hheight := n.rect.extents.y / 2
	size := Vector2{hwidth, hheight}

	n.children[0] = NewStaticNode(NewRect(n.rect.pos, size), n.capacity)
	n.children[1] = NewStaticNode(NewRect(n.rect.pos.Add(Vector2{hwidth, 0}), size), n.capacity)
	n.children[2] = NewStaticNode(NewRect(n.rect.pos.Add(Vector2{0, hheight}), size), n.capacity)
	n.children[3] = NewStaticNode(NewRect(n.rect.pos.Add(Vector2{hwidth, hheight}), size), n.capacity)

	remains := []Entity{}
	for _, val := range n.values {
		inserted := false
		for _, child := range n.children {
			if child.Insert(val) {
				inserted = true
				break
			}
		}
		if !inserted {
			remains = append(remains, val) // gather the nodes that dont fit
		}
	}

	n.values = remains
	n.leaf = false
}

func (n QNodeStatic) Query(area Rect) []Entity {
	if !n.rect.Intersects(area) {
		return []Entity{}
	}

	res := []Entity{}
	for _, val := range n.values {
		if val.rect.Intersects(area) {
			res = append(res, val)
		}
	}

	if !n.leaf {
		for _, child := range n.children {
			res = append(res, child.Query(area)...)
		}
	}

	return res
}

func (n QNodeStatic) Count() int {
	res := 0
	if !n.leaf {
		for _, child := range n.children {
			res += child.Count()
		}
	}
	return res + len(n.values)
}

func (n QNodeStatic) Draw(screen *ebiten.Image) {
	vector.StrokeRect(screen, float32(n.rect.pos.x), float32(n.rect.pos.y), float32(n.rect.extents.x), float32(n.rect.extents.y), 2, color.RGBA{255, 0, 0, 255}, false)

	for _, val := range n.values {
		vector.StrokeRect(screen, float32(val.rect.pos.x), float32(val.rect.pos.y), float32(val.rect.extents.x), float32(val.rect.extents.y), 2, color.RGBA{255, 0, 0, 255}, false)
	}

	if !n.leaf {
		for _, child := range n.children {
			child.Draw(screen)
		}
	}
}
