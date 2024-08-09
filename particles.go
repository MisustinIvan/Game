package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ParticleEmitter struct {
	pos          Vector2
	lifetime_min int
	lifetime_max int
	speed_min    float64
	speed_max    float64
	size_x       float32
	size_y       float32
	color        color.Color
	particles    *Particle
	last         *Particle
}

func NewParticleEmitter(pos Vector2, lifetime_min int, lifetime_max int, speed_min float64, speed_max float64, size_x float32, size_y float32, color color.Color) *ParticleEmitter {
	return &ParticleEmitter{
		pos:          pos,
		lifetime_min: lifetime_min,
		lifetime_max: lifetime_max,
		speed_min:    speed_min,
		speed_max:    speed_max,
		size_x:       size_x,
		size_y:       size_y,
		color:        color,
		particles:    nil,
		last:         nil,
	}
}

func (e *ParticleEmitter) Emit(dir Vector2) {
	lft := e.lifetime_min + int(float64(e.lifetime_max-e.lifetime_min)*rand.Float64())
	spd := e.speed_min + (e.speed_max-e.speed_min)*rand.Float64()

	p := NewParticle(e.pos, dir.Norm().Scale(spd), lft)

	if e.particles == nil {
		e.particles = p
	}
	if e.last != nil {
		e.last.next = p
	}
	e.last = p
}

func (e *ParticleEmitter) Update() {
	var pp *Particle = nil
	var cp *Particle = e.particles

	for cp != nil {
		cp.Update()

		if cp.Decayed() {
			if pp != nil {
				pp.next = cp.next
			} else {
				e.particles = cp.next
			}
		} else {
			pp = cp
		}
		cp = cp.next
	}
}

func (e *ParticleEmitter) Draw(screen *ebiten.Image) {
	cp := e.particles
	for cp != nil {
		sp := cp.pos.Sub(game.camera.rect.pos)
		vector.DrawFilledRect(screen, float32(sp.x), float32(sp.y), e.size_x, e.size_y, e.color, false)
		cp = cp.next
	}
}

type Particle struct {
	pos      Vector2
	vel      Vector2
	next     *Particle
	lifetime int
}

func NewParticle(pos Vector2, vel Vector2, lifetime int) *Particle {
	return &Particle{
		pos:      pos,
		vel:      vel,
		lifetime: lifetime,
		next:     nil,
	}
}

func (p *Particle) Update() {
	p.pos.AddEq(p.vel)
	p.lifetime -= 1
}

func (p Particle) Decayed() bool {
	return p.lifetime <= 0
}
