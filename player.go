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

type state int

const (
	idle = iota
	moving
)

const animation_timeout = 0.25 * FPS

// const emit_timeout = 0.025 * FPS
const emit_timeout = 1

type Player struct {
	pos                    Vector2
	hitbox                 Vector2
	health                 int
	xp                     int
	lvl                    int
	sprites                []*ebiten.Image
	sprite                 *ebiten.Image
	speed                  float64
	state                  state
	animation_task         *Task
	animation_index        int
	animations             int
	moving_particle_emiter ParticleEmitter
	emit_task              *Task
	bullet_manager         BulletManager
	debug                  bool
	dir
}

func NewPlayer(pos Vector2, health int, tm *TextureManager) *Player {
	var sprites = []*ebiten.Image{
		tm.GetTexture("robot0"),
		tm.GetTexture("robot1"),
		tm.GetTexture("robot2"),
		tm.GetTexture("robot3"),
	}

	p := &Player{
		pos:                    pos,
		hitbox:                 Vector2{float64(sprites[0].Bounds().Dx()), float64(sprites[0].Bounds().Dy())},
		health:                 health,
		xp:                     0,
		lvl:                    0,
		sprites:                sprites,
		sprite:                 sprites[0],
		dir:                    left,
		speed:                  1.25,
		state:                  idle,
		animation_index:        0,
		animations:             len(sprites),
		moving_particle_emiter: *NewParticleEmitter(pos.Add(Vector2{float64(sprites[0].Bounds().Dx()) - 10, float64(sprites[0].Bounds().Dy()) - 4}), 45, 60, 0.4, 0.6, 4, 4, color.RGBA{60, 60, 75, 255}),
		bullet_manager:         *NewBulletManager(pos.Add(Vector2{-16, 0}), 90, 2, 69, tm.GetTexture("bullet")),
		debug:                  false,
	}

	p.animation_task = NewTask(animation_timeout, func() {
		p.animation_index = (p.animation_index + 1) % p.animations
		p.sprite = p.sprites[p.animation_index]
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
		op.GeoM.Translate(p.hitbox.x, 0)
	}

	op.GeoM.Translate(p.pos.x, p.pos.y)

	screen.DrawImage(p.sprite, op)

	if p.debug {
		vector.StrokeRect(screen, float32(p.pos.x), float32(p.pos.y), float32(p.hitbox.x), float32(p.hitbox.y), 1, color.RGBA{255, 0, 0, 255}, false)
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

	p.bullet_manager.Shoot(dir)
}

func (p *Player) Emit() {
	vel := Vector2{0, 0}
	pos := p.pos
	pos.y += p.hitbox.y
	switch p.dir {
	case left:
		vel.x = 1
		pos.x = p.pos.x + p.hitbox.x
	case right:
		vel.x = -1
		pos.x = p.pos.x
	}

	vel.y = (rand.Float64() / -2) - 0.5

	speed := 0.5
	vel = vel.Norm().Scale(speed)

	p.moving_particle_emiter.Emit(vel)
}

func (p *Player) Update(g *Game) {
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

	if move {
		p.Move(diff)
		p.emit_task.Update()
	} else {
		p.state = idle
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Shoot()
	}

	p.animation_task.Update()

	p.moving_particle_emiter.Update()
	p.bullet_manager.Update()
}

func (p *Player) Move(diff Vector2) {
	p.state = moving
	if diff.x < 0 {
		p.dir = left
		p.moving_particle_emiter.pos.x = p.pos.x + p.hitbox.x - 10
		p.bullet_manager.pos.x = p.pos.x - 16
	} else if diff.x > 0 {
		p.dir = right
		p.moving_particle_emiter.pos.x = p.pos.x + 4
		p.bullet_manager.pos.x = p.pos.x + p.hitbox.x
	}

	p.pos = p.pos.Add(diff.Norm().Scale(p.speed))
	p.moving_particle_emiter.pos.y = p.pos.y + p.hitbox.y - 4
	p.bullet_manager.pos.y = p.pos.y
}
