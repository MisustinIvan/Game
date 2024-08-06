package main

import (
	"image/color"

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

type Player struct {
	pos             Vector2
	hitbox          Vector2
	health          int
	xp              int
	lvl             int
	sprites         []*ebiten.Image
	sprite          *ebiten.Image
	speed           float64
	state           state
	animation_tick  int
	animation_index int
	animations      int
	bullets         []*Bullet
	debug           bool
	dir
}

func NewPlayer(pos Vector2, health int, sprites []*ebiten.Image) Player {
	return Player{
		pos:             pos,
		hitbox:          Vector2{float64(sprites[0].Bounds().Dx()), float64(sprites[0].Bounds().Dy())},
		health:          health,
		xp:              0,
		lvl:             0,
		sprites:         sprites,
		sprite:          sprites[0],
		dir:             left,
		speed:           4,
		state:           idle,
		animation_tick:  0,
		animation_index: 0,
		animations:      len(sprites),
		bullets:         []*Bullet{},
		debug:           false,
	}
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

	for _, b := range p.bullets {
		b.Draw(screen)
	}
}

func (p *Player) Shoot(bs *ebiten.Image) {
	vel := Vector2{0, 0}
	if p.dir == left {
		vel.x = -1
	} else {
		vel.x = 1
	}
	p.bullets = append(p.bullets, NewBullet(p.pos, vel, bs))
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
	} else {
		p.state = idle
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Shoot(bulletSprite)
	}

	p.animation_tick += 1
	if p.animation_tick >= animation_timeout {
		p.animation_index = (p.animation_index + 1) % p.animations
		p.sprite = p.sprites[p.animation_index]
		p.animation_tick = 0
	}

	// TODO linked list or something idk
	for _, b := range p.bullets {
		if b != nil {
			b.Update()
		}
	}
}

func (p *Player) Move(diff Vector2) {
	p.state = moving
	if diff.x < 0 {
		p.dir = left
	} else if diff.x > 0 {
		p.dir = right
	}

	p.pos = p.pos.Add(diff.Norm().Scale(p.speed))
}
