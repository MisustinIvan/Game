package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const BulletSize = 16
const BulletAnimationTimeout int = 0.1 * FPS

type BulletManager struct {
	pos             Vector2
	bullet_lifetime int
	bullet_velocity float64
	bullet_damage   int
	texture_manager *TextureManager
	bullets         *Bullet
	last            *Bullet
	debug           *bool
}

func NewBulletManager(
	pos Vector2,
	bullet_lifetime int,
	bullet_velocity float64,
	bullet_damage int,
	texture_manager *TextureManager,
) *BulletManager {
	return &BulletManager{
		pos:             pos,
		bullet_lifetime: bullet_lifetime,
		bullet_velocity: bullet_velocity,
		bullet_damage:   bullet_damage,
		texture_manager: texture_manager,
		bullets:         nil,
		last:            nil,
	}
}

func (bm *BulletManager) Shoot(dir Vector2) {
	b := NewBullet(bm.pos, dir.Norm().Scale(bm.bullet_velocity), bm.bullet_lifetime, bm.bullet_damage, bm.texture_manager)

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
	pos        Vector2
	hitbox     Vector2
	vel        Vector2
	lifetime   int
	decay_time int
	damage     int
	next       *Bullet
	sprite     *ebiten.Image
	animator   *Animator[int]
	emitter    *ParticleEmitter
	emit_task  *Task
}

func NewBullet(pos Vector2, vel Vector2, lifetime int, damage int, tm *TextureManager) *Bullet {
	animation_sprites := []*ebiten.Image{
		tm.GetTexture("bullet0"),
		tm.GetTexture("bullet1"),
		tm.GetTexture("bullet2"),
		tm.GetTexture("bullet3"),
		tm.GetTexture("bullet4"),
		tm.GetTexture("bullet5"),
		tm.GetTexture("bullet6"),
		tm.GetTexture("bullet7"),
	}

	decay_spirtes := []*ebiten.Image{
		tm.GetTexture("bullet_decay_0"),
		tm.GetTexture("bullet_decay_1"),
		tm.GetTexture("bullet_decay_2"),
		tm.GetTexture("bullet_decay_3"),
		tm.GetTexture("bullet_decay_4"),
		tm.GetTexture("bullet_decay_5"),
	}

	b := &Bullet{
		pos:        pos,
		hitbox:     Vector2{float64(BulletSize), float64(BulletSize)},
		vel:        vel,
		lifetime:   lifetime,
		decay_time: len(decay_spirtes) * BulletAnimationTimeout,
		sprite:     animation_sprites[0],
		damage:     damage,
	}

	b.emitter = NewParticleEmitter(b.pos.Add(b.hitbox.Scale(0.5)), 20, 40, 0.6, 0.9, 3, 2, color.RGBA{215, 0, 255, 255})
	b.emit_task = NewTask(1, func() {
		vely := (rand.Float64() - 0.5) * 2
		b.emitter.Emit(b.vel.Scale(-0.25).Add(Vector2{0, vely}))
	})

	b.animator = NewAnimator[int](&b.sprite)
	b.animator.AddAnimation(0, NewAnimation(animation_sprites, BulletAnimationTimeout))
	b.animator.AddAnimation(1, NewAnimation(decay_spirtes, BulletAnimationTimeout))
	b.animator.SetAnimation(0)

	return b
}

func (b *Bullet) Update() {
	if b.StartedDecaying() {
		b.animator.SetAnimation(1)
		b.vel.ScaleEq(0.5)
	}

	if !b.Decaying() {
		b.emitter.pos = b.pos.Add(b.hitbox.Scale(0.5))
		b.emit_task.Update()
	}

	b.pos.AddEq(b.vel)
	b.lifetime -= 1
	b.animator.Update()
	b.emitter.Update()
}

func (b *Bullet) Decaying() bool {
	return b.lifetime < 0
}

func (b *Bullet) StartedDecaying() bool {
	return b.lifetime == 0
}

func (b *Bullet) Decayed() bool {
	return b.lifetime <= -b.decay_time
}

func (b *Bullet) Draw(screen *ebiten.Image, debug bool) {
	b.emitter.Draw(screen)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.sprite, op)
	if debug {
		vector.StrokeRect(screen, float32(b.pos.x), float32(b.pos.y), float32(b.hitbox.x), float32(b.hitbox.y), 2, color.RGBA{255, 0, 0, 255}, false)
	}
}
