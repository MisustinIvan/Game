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
	screen_pos := n.rect.pos.Sub(game.camera.rect.pos)
	vector.StrokeRect(screen, float32(screen_pos.x), float32(screen_pos.y), float32(n.rect.extents.x), float32(n.rect.extents.y), 1, color.RGBA{255, 255, 0, 255}, false)

	for _, val := range n.values {
		screen_pos = val.rect.pos.Sub(game.camera.rect.pos)
		//vector.StrokeRect(screen, float32(screen_pos.x), float32(screen_pos.y), float32(val.rect.extents.x), float32(val.rect.extents.y), 1, color.RGBA{255, 0, 0, 255}, false)
		ellen := game.texture_manager.GetTexture("mugshot")
		ew := float64(ellen.Bounds().Dx())
		eh := float64(ellen.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(val.rect.extents.x/ew, val.rect.extents.y/eh)
		op.GeoM.Translate(screen_pos.x, screen_pos.y)
		screen.DrawImage(ellen, op)
	}

	if !n.leaf {
		for _, child := range n.children {
			child.Draw(screen)
		}
	}
}

type QNodeDynamic struct {
	rect     Rect
	children [4]*QNodeDynamic
	capacity int
	values   []*Entity
	leaf     bool
	parent   *QNodeDynamic
}

func NewDynamicNode(rect Rect, capacity int) *QNodeDynamic {
	return &QNodeDynamic{
		rect:     rect,
		children: [4]*QNodeDynamic{},
		capacity: capacity,
		values:   []*Entity{},
		leaf:     true,
		parent:   nil,
	}
}

func (n *QNodeDynamic) Insert(e *Entity) bool {
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

func (n *QNodeDynamic) Subdivide() {
	hwidth := n.rect.extents.x / 2
	hheight := n.rect.extents.y / 2
	size := Vector2{hwidth, hheight}

	n.children[0] = NewDynamicNode(NewRect(n.rect.pos, size), n.capacity)
	n.children[0].parent = n
	n.children[1] = NewDynamicNode(NewRect(n.rect.pos.Add(Vector2{hwidth, 0}), size), n.capacity)
	n.children[1].parent = n
	n.children[2] = NewDynamicNode(NewRect(n.rect.pos.Add(Vector2{0, hheight}), size), n.capacity)
	n.children[2].parent = n
	n.children[3] = NewDynamicNode(NewRect(n.rect.pos.Add(Vector2{hwidth, hheight}), size), n.capacity)
	n.children[3].parent = n

	remains := []*Entity{}
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

func (n *QNodeDynamic) Update(e *Entity) bool {
	for i, val := range n.values {
		if val == e {
			if n.rect.Contains(e.rect) {
				return true
			}
			n.values = append(n.values[:i], n.values[i+1:]...)
			break
		}
	}

	node := n.parent
	for node != nil {
		if node.Insert(e) {
			return true
		}
		node = node.parent
	}
	return false
}

func (n QNodeDynamic) Query(area Rect) []*Entity {
	if !n.rect.Intersects(area) {
		return []*Entity{}
	}

	res := []*Entity{}
	for _, val := range n.values {
		if val.rect.Intersects(area) && val.rect != area {
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

func (n QNodeDynamic) Count() int {
	res := 0
	if !n.leaf {
		for _, child := range n.children {
			res += child.Count()
		}
	}
	return res + len(n.values)
}

func (n QNodeDynamic) Draw(screen *ebiten.Image) {
	screen_pos := n.rect.pos.Sub(game.camera.rect.pos)
	vector.StrokeRect(screen, float32(screen_pos.x), float32(screen_pos.y), float32(n.rect.extents.x), float32(n.rect.extents.y), 1, color.RGBA{255, 255, 0, 255}, false)

	for _, val := range n.values {
		screen_pos = val.rect.pos.Sub(game.camera.rect.pos)
		vector.StrokeRect(screen, float32(screen_pos.x), float32(screen_pos.y), float32(val.rect.extents.x), float32(val.rect.extents.y), 1, color.RGBA{255, 0, 0, 255}, false)
	}

	if !n.leaf {
		for _, child := range n.children {
			child.Draw(screen)
		}
	}
}
