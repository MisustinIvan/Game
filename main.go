package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const FPS = 120

type Game struct {
	player          *Player
	emitters        []*ParticleEmitter
	texture_manager *TextureManager
	background      *ebiten.Image
	walls_quadtree  *QNodeStatic
	camera          Camera
}

func NewGame() *Game {
	tm, err := NewTextureManager("./res/unknown.png")
	if err != nil {
		panic(err)
	}

	tm.LoadTexture("bullet", "./res/bullet_pink.png")
	tm.LoadTexture("bullet0", "./res/bullet_pink0.png")
	tm.LoadTexture("bullet1", "./res/bullet_pink1.png")
	tm.LoadTexture("bullet2", "./res/bullet_pink2.png")
	tm.LoadTexture("bullet3", "./res/bullet_pink3.png")
	tm.LoadTexture("bullet4", "./res/bullet_pink4.png")
	tm.LoadTexture("bullet5", "./res/bullet_pink5.png")
	tm.LoadTexture("bullet6", "./res/bullet_pink6.png")
	tm.LoadTexture("bullet7", "./res/bullet_pink7.png")
	tm.LoadTexture("bullet_decay_0", "./res/bullet_decay_0.png")
	tm.LoadTexture("bullet_decay_1", "./res/bullet_decay_1.png")
	tm.LoadTexture("bullet_decay_2", "./res/bullet_decay_2.png")
	tm.LoadTexture("bullet_decay_3", "./res/bullet_decay_3.png")
	tm.LoadTexture("bullet_decay_4", "./res/bullet_decay_4.png")
	tm.LoadTexture("bullet_decay_5", "./res/bullet_decay_5.png")
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
	tm.LoadTexture("ellen", "./res/ellen.png")

	tm.LoadTexture("background", "./res/background.png")
	wq := NewStaticNode(NewRect(Vector2{0, 0}, Vector2{960, 600}), 4)
	for i := 0; i < 50; i++ {

		rect := Rect{
			pos: Vector2{
				x: rand.Float64() * 940,
				y: rand.Float64() * 580,
			},
			extents: Vector2{
				x: rand.Float64() * 80,
				y: rand.Float64() * 80,
			},
		}

		wq.Insert(Entity{
			id:   i,
			rect: rect,
		})
	}

	player := NewPlayer(Vector2{100, 100}, 100, tm)
	camera := NewCamera(Vector2{960, 600}, &player.pos)

	return &Game{
		player: player,
		emitters: []*ParticleEmitter{
			NewParticleEmitter(Vector2{100, 150}, 90, 120, 0.3, 0.5, 2, 6, color.RGBA{255, 30, 150, 255}),
			NewParticleEmitter(Vector2{200, 200}, 60, 90, 0.6, 0.8, 2, 2, color.RGBA{30, 255, 150, 255}),
			NewParticleEmitter(Vector2{300, 150}, 120, 150, 0.2, 0.4, 3, 3, color.RGBA{150, 30, 255, 255}),
		},
		texture_manager: tm,
		background:      tm.GetTexture("background"),
		walls_quadtree:  wq,
		camera:          camera,
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
	g.camera.Update()

	for _, emitter := range g.emitters {
		dir := Vector2{x: (rand.Float64() - 0.5) * 2, y: (-rand.Float64() / 2) - 0.5}
		emitter.Emit(dir)
		emitter.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 55, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.camera.rect.pos.x, -g.camera.rect.pos.y)
	screen.DrawImage(g.background, op)

	for _, emitter := range g.emitters {
		emitter.Draw(screen, g)
	}
	g.player.Draw(screen, g)
	g.walls_quadtree.Draw(screen, g)

	for _, val := range g.walls_quadtree.Query(NewRect(g.player.pos, g.player.hitbox)) {
		sp := val.rect.pos.Sub(g.camera.rect.pos)
		vector.StrokeRect(screen, float32(sp.x), float32(sp.y), float32(val.rect.extents.x), float32(val.rect.extents.y), 2, color.RGBA{0, 0, 255, 255}, false)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\nFPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth / 2, outsideHeight / 2
}

func main() {
	ebiten.SetWindowSize(1920, 1200)
	ebiten.SetTPS(FPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
