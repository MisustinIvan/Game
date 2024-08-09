package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type dir int

const (
	left = iota
	right
)

const animation_timeout = 0.25 * FPS
const emit_timeout = 1
const player_size = 32

type Player struct {
	rect                   Rect
	health                 int
	xp                     int
	lvl                    int
	sprite                 *ebiten.Image
	speed                  float64
	state                  PlayerState
	animator               *Animator[PlayerState]
	animation_task         *Task
	moving_particle_emiter ParticleEmitter
	emit_task              *Task
	attack_timer           int
	bullet_manager         BulletManager
	debug                  bool
	dir
}

func NewPlayer(pos Vector2, health int, tm *TextureManager) *Player {
	var idle_sprites = []*ebiten.Image{
		tm.GetTexture("robot_idle_0"),
		tm.GetTexture("robot_idle_1"),
		tm.GetTexture("robot_idle_2"),
		tm.GetTexture("robot_idle_3"),
	}

	var moving_spites = []*ebiten.Image{
		tm.GetTexture("robot_moving_0"),
		tm.GetTexture("robot_moving_1"),
		tm.GetTexture("robot_moving_2"),
		tm.GetTexture("robot_moving_3"),
	}

	var attack_sprites = []*ebiten.Image{
		tm.GetTexture("robot_attack_0"),
		tm.GetTexture("robot_attack_1"),
		tm.GetTexture("robot_attack_2"),
		tm.GetTexture("robot_attack_3"),
	}

	var attack_moving_sprites = []*ebiten.Image{
		tm.GetTexture("robot_attack_moving_0"),
		tm.GetTexture("robot_attack_moving_1"),
		tm.GetTexture("robot_attack_moving_2"),
		tm.GetTexture("robot_attack_moving_3"),
	}

	p := &Player{
		rect:                   NewRect(pos, Vector2{float64(player_size), float64(player_size)}),
		health:                 health,
		xp:                     0,
		lvl:                    0,
		sprite:                 idle_sprites[0],
		dir:                    left,
		speed:                  1.25,
		state:                  PlayerIdle,
		attack_timer:           0,
		moving_particle_emiter: *NewParticleEmitter(pos.Add(Vector2{float64(player_size) - 10, float64(player_size) - 4}), 45, 60, 0.4, 0.6, 4, 4, color.RGBA{60, 60, 75, 255}),
		bullet_manager:         *NewBulletManager(pos.Add(Vector2{-16, 0}), 120, 3, 69, tm),
		debug:                  false,
	}

	p.animator = NewAnimator[PlayerState](&p.sprite)
	p.animator.AddAnimation(
		PlayerIdle,
		NewAnimation(
			idle_sprites, animation_timeout,
		),
	)

	p.animator.AddAnimation(
		PlayerMoving,
		NewAnimation(
			moving_spites, animation_timeout,
		),
	)

	p.animator.AddAnimation(
		PlayerAttacking,
		NewAnimation(
			attack_sprites, animation_timeout,
		),
	)

	p.animator.AddAnimation(
		PlayerMovingAttacking,
		NewAnimation(
			attack_moving_sprites, animation_timeout,
		),
	)

	p.animator.SetAnimation(PlayerIdle)

	p.animation_task = NewTask(1, func() {
		if p.animator.ckey != p.state {
			p.animator.SetAnimation(p.state)
		}

		p.animator.Update()
	})

	p.emit_task = NewTask(emit_timeout, func() {
		p.Emit()
	})

	return p
}

func (p Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if p.dir == right {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(p.rect.extents.x-1, 0)
	}

	screen_pos := p.rect.pos.Sub(game.camera.rect.pos)

	op.GeoM.Translate(screen_pos.x, screen_pos.y)

	screen.DrawImage(p.sprite, op)

	if p.debug {
		vector.StrokeRect(screen, float32(screen_pos.x), float32(screen_pos.y), float32(p.rect.extents.x), float32(p.rect.extents.y), 1, color.RGBA{255, 0, 0, 255}, false)
	}

	p.moving_particle_emiter.Draw(screen)
	p.bullet_manager.Draw(screen, p.debug)
}

func (p *Player) Shoot() {
	dir := Vector2{0, 0}
	if p.dir == left {
		dir.x = -1
	} else {
		dir.x = 1
	}

	dir.y = (rand.Float64() - 0.5) * 2

	p.bullet_manager.Shoot(dir)

	if p.state == PlayerMoving {
		p.state = PlayerMovingAttacking
	} else {
		p.state = PlayerAttacking
	}

	p.attack_timer = 2 * animation_timeout
}

func (p *Player) Emit() {
	vel := Vector2{0, 0}
	pos := p.rect.pos
	pos.y += p.rect.extents.y
	switch p.dir {
	case left:
		vel.x = 1
		pos.x = p.rect.pos.x + p.rect.extents.x
	case right:
		vel.x = -1
		pos.x = p.rect.pos.x
	}

	vel.y = (rand.Float64() / -2) - 0.5

	speed := 0.5
	vel = vel.Norm().Scale(speed)

	p.moving_particle_emiter.Emit(vel)
}

func (p *Player) Update() {
	diff := Vector2{0, 0}
	move := false
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		diff.y = -1
		move = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		diff.y = 1
		move = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		diff.x = -1
		move = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		diff.x = 1
		move = true
	}

	p.animation_task.Update()

	if move {
		p.Move(diff)
		p.emit_task.Update()
	} else if p.attack_timer == 0 {
		p.state = PlayerIdle
	}

	if p.attack_timer > 0 {
		p.attack_timer -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Shoot()
	}

	p.animation_task.Update()

	p.moving_particle_emiter.Update()
	p.bullet_manager.Update()
}

func (p *Player) Move(dir Vector2) {
	if p.attack_timer > 0 {
		p.state = PlayerMovingAttacking
	} else {
		p.state = PlayerMoving
	}

	movement_vec := dir.Norm().Scale(p.speed)
	// horizontal
	p.rect.pos.x += movement_vec.x
	walls := game.walls_quadtree.Query(p.rect)
	for _, wall := range walls {
		if wall.rect.Intersects(p.rect) {
			p.rect.ResolveX(wall.rect, movement_vec)
		}
	}
	// vertical
	p.rect.pos.y += movement_vec.y
	walls = game.walls_quadtree.Query(p.rect)
	for _, wall := range walls {
		if wall.rect.Intersects(p.rect) {
			p.rect.ResolveY(wall.rect, movement_vec)
		}
	}

	if dir.x < 0 {
		p.dir = left
		p.moving_particle_emiter.pos.x = p.rect.pos.x + p.rect.extents.x - 10
		p.bullet_manager.pos.x = p.rect.pos.x - 16
	} else if dir.x > 0 {
		p.dir = right
		p.moving_particle_emiter.pos.x = p.rect.pos.x + 4
		p.bullet_manager.pos.x = p.rect.pos.x + p.rect.extents.x
	}

	p.moving_particle_emiter.pos.y = p.rect.pos.y + p.rect.extents.y - 4
	p.bullet_manager.pos.y = p.rect.pos.y
}
