package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const FPS = 120

type Game struct {
	player          *Player
	emitters        []*ParticleEmitter
	texture_manager *TextureManager
}

func NewGame() *Game {
	tm, err := NewTextureManager("./res/unknown.png")
	if err != nil {
		panic(err)
	}

	tm.LoadTexture("bullet", "./res/bullet_pink.png")
	tm.LoadTexture("robot_idle_0", "./res/robot_idle_0.png")
	tm.LoadTexture("robot_idle_1", "./res/robot_idle_1.png")
	tm.LoadTexture("robot_idle_2", "./res/robot_idle_2.png")
	tm.LoadTexture("robot_idle_3", "./res/robot_idle_3.png")
	tm.LoadTexture("robot_moving_0", "./res/robot_moving_0.png")
	tm.LoadTexture("robot_moving_1", "./res/robot_moving_1.png")
	tm.LoadTexture("robot_moving_2", "./res/robot_moving_2.png")
	tm.LoadTexture("robot_moving_3", "./res/robot_moving_3.png")
	tm.LoadTexture("robot_attack_0", "./res/robot_attack_0.png")
	tm.LoadTexture("robot_attack_1", "./res/robot_attack_1.png")
	tm.LoadTexture("robot_attack_2", "./res/robot_attack_2.png")
	tm.LoadTexture("robot_attack_3", "./res/robot_attack_3.png")
	tm.LoadTexture("robot_attack_moving_0", "./res/robot_attack_moving_0.png")
	tm.LoadTexture("robot_attack_moving_1", "./res/robot_attack_moving_1.png")
	tm.LoadTexture("robot_attack_moving_2", "./res/robot_attack_moving_2.png")
	tm.LoadTexture("robot_attack_moving_3", "./res/robot_attack_moving_3.png")

	return &Game{
		player: NewPlayer(Vector2{100, 100}, 100, tm),
		emitters: []*ParticleEmitter{
			NewParticleEmitter(Vector2{100, 150}, 90, 120, 0.3, 0.5, 2, 6, color.RGBA{255, 30, 150, 255}),
			NewParticleEmitter(Vector2{200, 200}, 60, 90, 0.6, 0.8, 2, 2, color.RGBA{30, 255, 150, 255}),
			NewParticleEmitter(Vector2{300, 150}, 120, 150, 0.2, 0.4, 3, 3, color.RGBA{150, 30, 255, 255}),
		},
		texture_manager: tm,
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

	for _, emitter := range g.emitters {
		emitter.Draw(screen)
	}
	g.player.Draw(screen)
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\nFPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS()))
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
