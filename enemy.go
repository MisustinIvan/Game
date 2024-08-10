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
	cc             CircleCollider
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
		cc:     CircleCollider{pos: pos.Add(Vector2{EnemySize / 2, EnemySize / 2}), r: EnemySize / 2},
		damage: 10,
		health: 50,
		speed:  0.33,
		sprite: idle_sprites[0],
		target: Vector2{0, 0},
		state:  EnemyIdle,
	}
	animator := NewAnimator[EnemyState](&e.sprite)
	animator.AddAnimation(EnemyIdle, NewAnimation(idle_sprites, EnemyAnimationTimeout))
	animator.SetAnimation(EnemyIdle)

	e.animator = animator

	return e
}

func (e *Enemy) Update() {
	e.target = game.player.rect.pos

	diff := e.target.Sub(e.cc.pos)
	// DONT REMOVE, ELSE ENEMIES VANISH INTO THE IEEE.754 SHADOW REALM
	// Division by zero happens...
	if diff.Mag() != 0 {
		e.Move(diff)
	}

	for _, other := range game.enemies_grid.GetNearbyEnemies(e.cc.pos) {
		if e.cc.Collides(other.cc) {
			e.Resolve(other)
		}
	}
}

func (e *Enemy) Resolve(o *Enemy) {
	diff := e.cc.pos.Sub(o.cc.pos)
	distance := diff.Mag()

	totalRadius := e.cc.r + o.cc.r

	if distance < totalRadius && distance != 0 {
		overlap := totalRadius - distance
		normal := diff.Norm()
		separation := normal.Scale(overlap / 2)
		e.cc.pos.AddEq(separation)
		o.cc.pos.SubEq(separation)
	}
}

func (e *Enemy) Move(dir Vector2) {
	if e.state == EnemyAttacking {
		e.state = EnemyAttackingMoving
	} else {
		e.state = EnemyMoving
	}

	e.cc.pos.AddEq(dir.Norm().Scale(e.speed))

	if e.cc.pos.x < 0 {
		e.cc.pos.x = 0
	}

	if e.cc.pos.y < 0 {
		e.cc.pos.y = 0
	}

	if dir.x < 0 {
		e.dir = left
	} else if dir.x > 0 {
		e.dir = right
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen_pos := e.cc.pos.Sub(game.camera.rect.pos)

	w := float64(e.sprite.Bounds().Dx())
	h := float64(e.sprite.Bounds().Dy())
	op.GeoM.Scale(EnemySize/w, EnemySize/h)
	if e.dir == right {
		op.GeoM.Scale(-1, 1)
		screen_pos.x += EnemySize - 1
	}
	op.GeoM.Translate(screen_pos.x-EnemySize/2, screen_pos.y-EnemySize/2)
	screen.DrawImage(e.sprite, op)
}
