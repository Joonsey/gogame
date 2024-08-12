package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	Speed    int
	MoveCooldown float32
	BaseMoveCooldown float32
}

type Game struct {
	Player  Player
	mapData [][]int
	Client  *Client
	FrameCount uint64
}

func (g *Game) CheckCollision(pos Position) bool {
	return false
}

func (g *Game) Update() error {
	player_pos := &g.Player.Position
	init_player_pos := *player_pos
	g.FrameCount++

	g.Player.MoveCooldown = max(0, g.Player.MoveCooldown - 1)

	if g.Player.MoveCooldown == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			player_pos.Y -= g.Player.Speed
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			player_pos.Y += g.Player.Speed
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}

		if g.CheckCollision(*player_pos) {
			player_pos.Y = init_player_pos.Y
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			player_pos.X -= g.Player.Speed
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			player_pos.X += g.Player.Speed
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}
		if g.CheckCollision(*player_pos) {
			player_pos.X = init_player_pos.X
			g.Player.MoveCooldown = g.Player.BaseMoveCooldown
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	g.Client.SendPosition(CoordinateData{float32(player_pos.X), float32(player_pos.Y)})
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
			vector.DrawFilledRect(screen, float32(x*TILE_SIZE), float32(y*TILE_SIZE), TILE_SIZE, TILE_SIZE, clr, false)
		}
	}

	rect_color := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	vector.DrawFilledRect(screen, float32(g.Player.Position.X), float32(g.Player.Position.Y), 16, 16, rect_color, false)

	other_rect_color := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	vector.DrawFilledRect(screen, g.Client.other_pos.X, g.Client.other_pos.Y, 16, 16, other_rect_color, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func createMap() [][]int {
	width_count := SCREEN_WIDTH / TILE_SIZE
	height_count := SCREEN_HEIGHT / TILE_SIZE

	mapData := make([][]int, height_count)

	for y := 0; y < height_count; y++ {
		mapData[y] = make([]int, width_count)
		for x := 0; x < width_count; x++ {
			val := 0
			if rand.Intn(10) < 1 {
				val = 1
			}
			mapData[y][x] = val
		}
	}

	return mapData
}

func expandMap(mapData [][]int, additionalRows int) [][]int {
	width_count := len(mapData[0])

	for i := 0; i < additionalRows; i++ {
		newRow := make([]int, width_count)
		mapData = append(mapData, newRow)
	}

	return mapData
}

func main() {
	is_server := flag.String("server", "y", "run server")
	server_ip := flag.String("ip", "172.20.10.2", "ip")

	flag.Parse()

	if *is_server == "y" {
		RunServer()
		return
	}

	mapData := createMap()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	client := Client{}
	game := Game{Player: Player{Speed: TILE_SIZE, BaseMoveCooldown: 15}, mapData: mapData, Client: &client}

	go client.RunClient(*server_ip)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
