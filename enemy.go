package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyState int

const (
	EnemyIdle EnemyState = iota
	EnemyMoving
	EnemyAttacking
	EnemyAttackingMoving
	EnemyDying
)

const EnemyAnimationTimeout = 0.25 * FPS
const EnemySize = 32

type Enemy struct {
	rect           Rect
	damage         int
	health         int
	speed          float64
	sprite         *ebiten.Image
	target         Vector2
	dir            dir
	state          EnemyState
	animator       *Animator[EnemyState]
	animation_task *Task
}

func NewEnemy(pos Vector2, tm *TextureManager) *Enemy {

	idle_sprites := []*ebiten.Image{
		tm.GetTexture("ellen"),
		tm.GetTexture("ellen"),
		tm.GetTexture("ellen"),
	}

	e := &Enemy{
		rect:   NewRect(pos, Vector2{EnemySize, EnemySize}),
		damage: 10,
		health: 50,
		speed:  1,
		sprite: idle_sprites[0],
		target: Vector2{0, 0},
		state:  EnemyIdle,
	}
	animator := NewAnimator[EnemyState](&e.sprite)
	animator.AddAnimation(EnemyIdle, NewAnimation(idle_sprites, EnemyAnimationTimeout))
	animator.SetAnimation(EnemyIdle)

	e.animator = animator

	// TODO ANIMATOR TASK

	return e
}

func (e *Enemy) Update() {
	e.target = game.player.rect.pos

	diff := e.target.Sub(e.rect.pos)
	// DONT REMOVE, ELSE ENEMIES VANISH INTO THE IEEE.754 SHADOW REALM
	// Division by zero happens...
	if diff.Mag() != 0 {
		e.Move(diff)
	}
}

func (e *Enemy) Move(dir Vector2) {
	if e.state == EnemyAttacking {
		e.state = EnemyAttackingMoving
	} else {
		e.state = EnemyMoving
	}

	movement_vec := dir.Norm().Scale(e.speed)
	// horizontal
	e.rect.pos.x += movement_vec.x
	walls := append(game.walls_quadtree.Query(e.rect), e.GetOtherEnemiesRects()...)
	for _, wall := range walls {
		if wall.rect.Intersects(e.rect) {
			e.rect.ResolveX(wall.rect, movement_vec)
		}
	}

	// vertical
	e.rect.pos.y += movement_vec.y
	walls = append(game.walls_quadtree.Query(e.rect), e.GetOtherEnemiesRects()...)
	for _, wall := range walls {
		if wall.rect.Intersects(e.rect) {
			e.rect.ResolveY(wall.rect, movement_vec)
		}
	}

	if dir.x < 0 {
		e.dir = left
	} else if dir.x > 0 {
		e.dir = right
	}
}

func (e Enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen_pos := e.rect.pos.Sub(game.camera.rect.pos)

	w := float64(e.sprite.Bounds().Dx())
	h := float64(e.sprite.Bounds().Dy())
	op.GeoM.Scale(e.rect.extents.x/w, e.rect.extents.y/h)
	if e.dir == right {
		op.GeoM.Scale(-1, 1)
		screen_pos.x += e.rect.extents.x - 1
	}
	op.GeoM.Translate(screen_pos.x, screen_pos.y)
	screen.DrawImage(e.sprite, op)
}

func (e *Enemy) GetOtherEnemiesRects() []Entity {
	res := []Entity{}
	for _, enemy := range game.enemies {
		if enemy != e {
			res = append(res, Entity{69420, enemy.rect})
		}
	}
	return res
}
