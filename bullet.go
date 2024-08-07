package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BulletManager struct {
	pos             Vector2
	bullet_lifetime int
	bullet_velocity float64
	bullet_damage   int
	bullet_sprite   *ebiten.Image
	bullets         *Bullet
	last            *Bullet
	debug           *bool
}

func NewBulletManager(
	pos Vector2,
	bullet_lifetime int,
	bullet_velocity float64,
	bullet_damage int,
	bullet_sprite *ebiten.Image,
) *BulletManager {
	return &BulletManager{
		pos:             pos,
		bullet_lifetime: bullet_lifetime,
		bullet_velocity: bullet_velocity,
		bullet_damage:   bullet_damage,
		bullet_sprite:   bullet_sprite,
		bullets:         nil,
		last:            nil,
	}
}

func (bm *BulletManager) Shoot(dir Vector2) {
	b := NewBullet(bm.pos, dir.Norm().Scale(bm.bullet_velocity), bm.bullet_sprite, bm.bullet_lifetime, bm.bullet_damage)

	if bm.bullets == nil {
		bm.bullets = b
	}

	if bm.last != nil {
		bm.last.next = b
	}

	bm.last = b
}

func (bm *BulletManager) Update() {
	var pb *Bullet = nil
	var cb *Bullet = bm.bullets

	for cb != nil {
		cb.Update()

		if cb.Decayed() {
			if pb != nil {
				pb.next = cb.next
			} else {
				bm.bullets = cb.next
			}
		} else {
			pb = cb
		}
		cb = cb.next
	}
}

func (bm *BulletManager) Draw(screen *ebiten.Image, debug bool) {
	cb := bm.bullets
	for cb != nil {
		cb.Draw(screen, debug)
		cb = cb.next
	}
}

type Bullet struct {
	pos      Vector2
	hitbox   Vector2
	vel      Vector2
	lifetime int
	damage   int
	next     *Bullet
	sprite   *ebiten.Image
}

func NewBullet(pos Vector2, vel Vector2, sprite *ebiten.Image, lifetime int, damage int) *Bullet {
	return &Bullet{
		pos:      pos,
		hitbox:   Vector2{float64(sprite.Bounds().Dx()), float64(sprite.Bounds().Dy())},
		vel:      vel,
		lifetime: lifetime,
		sprite:   sprite,
		damage:   damage,
	}
}

func (b *Bullet) Update() {
	b.pos.AddEq(b.vel)
	b.lifetime -= 1
}

func (b *Bullet) Decayed() bool {
	return b.lifetime <= 0
}

func (b *Bullet) Draw(screen *ebiten.Image, debug bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.sprite, op)
	if debug {
		vector.StrokeRect(screen, float32(b.pos.x), float32(b.pos.y), float32(b.hitbox.x), float32(b.hitbox.y), 2, color.RGBA{255, 0, 0, 255}, false)
	}
}
