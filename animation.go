package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	keyframes     []*ebiten.Image
	current_frame int
	delay         int
}

func NewAnimation(keyframes []*ebiten.Image, delay int) *Animation {
	return &Animation{
		keyframes:     keyframes,
		current_frame: 0,
		delay:         delay,
	}
}

func (a *Animation) Step() *ebiten.Image {
	a.current_frame = (a.current_frame + 1) % len(a.keyframes)
	return a.keyframes[a.current_frame]
}

type PlayerState int

const (
	PlayerIdle PlayerState = iota
	PlayerMoving
	PlayerAttacking
	PlayerMovingAttacking
	PlayerTakingDamage
)

type Animator[T ~int] struct {
	target     **ebiten.Image
	animations map[T]*Animation
	ckey       T
	animation  *Animation
	elapsed    int
}

func NewAnimator[T ~int](target **ebiten.Image) *Animator[T] {
	return &Animator[T]{
		target:     target,
		animations: map[T]*Animation{},
		ckey:       0,
		animation:  nil,
		elapsed:    0,
	}
}

func (a *Animator[T]) Update() {
	a.elapsed = (a.elapsed + 1) % a.animation.delay
	if a.elapsed == 0 {
		*a.target = a.animation.Step()
	}
}

func (a *Animator[T]) AddAnimation(key T, animation *Animation) {
	a.animations[key] = animation
}

func (a *Animator[T]) SetAnimation(key T) {
	a.ckey = key
	a.animation = a.animations[key]
	a.animation.current_frame = 0
}
