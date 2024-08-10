package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CircleCollider struct {
	pos Vector2
	r   float64
}

func (c *CircleCollider) Collides(o CircleCollider) bool {
	return (c.pos.Sub(o.pos).Mag() < c.r+o.r)
}

type SpatialGrid struct {
	cells    [][]*Enemy
	cellSize float64
	width    int
	height   int
}

func NewSpatialGrid(width, height int, cellSize float64) SpatialGrid {
	cells := make([][]*Enemy, width*height)
	return SpatialGrid{
		cells:    cells,
		cellSize: cellSize,
		height:   height,
		width:    width,
	}
}
func (s *SpatialGrid) Insert(enemy *Enemy) {
	index := s.PosToIndex(enemy.cc.pos)
	if index >= 0 && index < len(s.cells) {
		s.cells[index] = append(s.cells[index], enemy)
	}
}

func (g *SpatialGrid) PosToIndex(pos Vector2) int {
	x := int(math.Floor(pos.x / g.cellSize))
	y := int(math.Floor(pos.y / g.cellSize))
	return y*g.width + x
}

func (s *SpatialGrid) Clear() {
	for i := range s.cells {
		s.cells[i] = []*Enemy{}
	}
}

func (s *SpatialGrid) GetNearbyEnemies(pos Vector2) []*Enemy {
	var nearbyEnemies []*Enemy
	cellX := int(math.Floor(pos.x / s.cellSize))
	cellY := int(math.Floor(pos.y / s.cellSize))

	for x := cellX - 1; x <= cellX+1; x++ {
		for y := cellY - 1; y <= cellY+1; y++ {
			if x >= 0 && x < s.width && y >= 0 && y < s.height {
				index := y*s.width + x
				nearbyEnemies = append(nearbyEnemies, s.cells[index]...)
			}
		}
	}

	return nearbyEnemies
}

func DebugDrawEnemies(screen *ebiten.Image, enemies []*Enemy) {
	for _, enemy := range enemies {
		screen_pos := enemy.cc.pos.Sub(game.camera.rect.pos)
		vector.StrokeRect(screen, float32(screen_pos.x)-EnemySize/2, float32(screen_pos.y)-EnemySize/2, EnemySize, EnemySize, 1, color.RGBA{0, 0, 255, 255}, false)
	}
}
