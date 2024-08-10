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
	player          *Player
	emitters        []*ParticleEmitter
	texture_manager *TextureManager
	background      *ebiten.Image
	camera          Camera
	enemies         []*Enemy
	enemies_grid    SpatialGrid
}

var game *Game

func NewGame() *Game {
	tm, err := NewTextureManager("./res/unknown.png")
	if err != nil {
		panic(err)
	}

	err = tm.LoadTextures("./res/")
	if err != nil {
		panic(err)
	}

	player := NewPlayer(Vector2{100, 100}, 100, tm)
	camera := NewCamera(Vector2{960, 600}, &player.rect.pos)

	enemies_grid := NewSpatialGrid(100, 100, 32)
	enemies := []*Enemy{}
	for i := 0; i < 200; i++ {
		enemies = append(enemies, NewEnemy(
			Vector2{
				(rand.Float64() - 0.5) * 2 * 500,
				(rand.Float64() - 0.5) * 2 * 500,
			},
			tm,
		))
	}

	return &Game{
		player: player,
		emitters: []*ParticleEmitter{
			NewParticleEmitter(Vector2{100, 150}, 90, 120, 0.3, 0.5, 2, 6, color.RGBA{255, 30, 150, 255}),
			NewParticleEmitter(Vector2{200, 200}, 60, 90, 0.6, 0.8, 2, 2, color.RGBA{30, 255, 150, 255}),
			NewParticleEmitter(Vector2{300, 150}, 120, 150, 0.2, 0.4, 3, 3, color.RGBA{150, 30, 255, 255}),
		},
		texture_manager: tm,
		background:      tm.GetTexture("background"),
		camera:          camera,
		enemies:         enemies,
		enemies_grid:    enemies_grid,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.player.debug = !g.player.debug
	}

	g.player.Update()
	g.camera.Update()
	for _, enemy := range g.enemies {
		enemy.Update()
	}

	for _, emitter := range g.emitters {
		dir := Vector2{x: (rand.Float64() - 0.5) * 2, y: (-rand.Float64() / 2) - 0.5}
		emitter.Emit(dir)
		emitter.Update()
	}

	game.enemies_grid.Clear()

	for _, enemy := range g.enemies {
		game.enemies_grid.Insert(enemy)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 55, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.camera.rect.pos.x, -g.camera.rect.pos.y)
	screen.DrawImage(g.background, op)

	for _, emitter := range g.emitters {
		emitter.Draw(screen)
	}
	g.player.Draw(screen)
	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}

	DebugDrawEnemies(screen, game.enemies_grid.GetNearbyEnemies(game.player.rect.pos))

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\nFPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth / 2, outsideHeight / 2
}

func main() {
	ebiten.SetWindowSize(1920, 1200)
	ebiten.SetTPS(FPS)
	game = NewGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
