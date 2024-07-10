package main

import (
	"fmt"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
)

const (
	SCREEN_WIDTH  = 320
	SCREEN_HEIGHT = 240
	TILE_SIZE     = 16
)

type Position struct {
	X, Y int
}

type Player struct {
	Position Position
	Speed int
}

type Game struct {
	Player  Player
	mapData [][]int
}

func (g *Game) CheckCollision(pos Position) bool {
	posX, posY := int(math.Round(float64(pos.X) / float64(TILE_SIZE))), int(math.Round(float64(pos.Y) / float64(TILE_SIZE)))
    if posX < 0 || posX >= len(g.mapData[0]) || posY < 0 || posY >= len(g.mapData) {
        return false
    }

    return g.mapData[posY][posX] != 0
}


func (g *Game) Update() error {
	player_pos := &g.Player.Position
	init_player_pos := *player_pos

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		player_pos.Y -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		player_pos.Y += g.Player.Speed
	}

	if g.CheckCollision(*player_pos){ player_pos.Y = init_player_pos.Y }
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player_pos.X -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		player_pos.X += g.Player.Speed
	}
	if g.CheckCollision(*player_pos){ player_pos.X = init_player_pos.X }

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	tps := ebiten.ActualTPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f", tps))
	for y, row := range g.mapData {
		for x, cell := range row {
			var clr color.Color
			switch cell {
			case 0:
				clr = color.RGBA{128, 128, 0, 255} // Walkable
			case 1:
				clr = color.RGBA{0, 0, 128, 255} // Obstacle
			default:
				clr = color.RGBA{128, 128, 128, 255} // Other types
			}
			ebitenutil.DrawRect(screen, float64(x*TILE_SIZE), float64(y*TILE_SIZE), TILE_SIZE, TILE_SIZE, clr)
		}
	}

	rect_color := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	ebitenutil.DrawRect(screen, float64(g.Player.Position.X), float64(g.Player.Position.Y), 16, 16, rect_color)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	mapData := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 1, 0, 1, 1, 1, 1, 0, 1, 0},
		{0, 1, 0, 1, 0, 0, 1, 0, 1, 0},
		{0, 0, 0, 1, 0, 0, 1, 0, 1, 0},
		{0, 1, 0, 1, 1, 1, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	game := Game{Player: Player{Speed: 4}, mapData: mapData}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
