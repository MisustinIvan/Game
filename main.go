package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const FPS = 120

type Game struct {
	player Player
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 55, 255})

	g.player.Draw(screen)
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
