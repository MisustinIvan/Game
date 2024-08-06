package main

import "github.com/hajimehoshi/ebiten/v2"

const BulletLifetime = 1 * FPS
const BulletSpeed = 2

type Bullet struct {
	pos      Vector2
	vel      Vector2
	lifetime int
	sprite   *ebiten.Image
}

func NewBullet(pos Vector2, vel Vector2, sprite *ebiten.Image) *Bullet {
	return &Bullet{
		pos:      pos,
		vel:      vel.Norm().Scale(BulletSpeed),
		lifetime: BulletLifetime,
		sprite:   sprite,
	}
}

func (b *Bullet) Update() {
	b.pos.AddEq(b.vel)
	b.lifetime -= 1
	if b.lifetime <= 0 {
		b = nil
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.sprite, op)
}
