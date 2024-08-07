package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const FPS = 120

type Game struct {
	player   Player
	emitters []*ParticleEmitter
}

var bulletSprite *ebiten.Image

func NewGame() *Game {

	image0, _, err := ebitenutil.NewImageFromFile("./res/robot0.png")
	if err != nil {
		panic(err)
	}

	image1, _, err := ebitenutil.NewImageFromFile("./res/robot1.png")
	if err != nil {
		panic(err)
	}

	image2, _, err := ebitenutil.NewImageFromFile("./res/robot2.png")
	if err != nil {
		panic(err)
	}

	image3, _, err := ebitenutil.NewImageFromFile("./res/robot3.png")
	if err != nil {
		panic(err)
	}

	bulletSprite, _, err = ebitenutil.NewImageFromFile("./res/bullet.png")
	if err != nil {
		panic(err)
	}

	return &Game{
		player: NewPlayer(Vector2{100, 100}, 100, []*ebiten.Image{image0, image1, image2, image3}),
		emitters: []*ParticleEmitter{
			NewParticleEmitter(Vector2{100, 150}, 90, 120, 0.3, 0.5, 2, 6, color.RGBA{255, 30, 150, 255}),
			NewParticleEmitter(Vector2{200, 200}, 60, 90, 0.6, 0.8, 2, 2, color.RGBA{30, 255, 150, 255}),
			NewParticleEmitter(Vector2{300, 150}, 120, 150, 0.2, 0.4, 3, 3, color.RGBA{150, 30, 255, 255}),
		},
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.player.debug = !g.player.debug
	}

	g.player.Update(g)

	for _, emitter := range g.emitters {
		dir := Vector2{x: (rand.Float64() - 0.5) * 2, y: (-rand.Float64() / 2) - 0.5}
		emitter.Emit(dir)
		emitter.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 55, 255})

	g.player.Draw(screen)
	for _, emitter := range g.emitters {
		emitter.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\nFPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth / 2, outsideHeight / 2
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetTPS(FPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
