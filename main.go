package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	pixels []byte
	sanbox *Sandbox

	selectedCellType CellType
}

func (g *Game) HandleInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		g.sanbox.SetCell(cursorX, cursorY, *NewCell(g.selectedCellType))
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selectedCellType = SAND
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selectedCellType = WATER
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.selectedCellType = WALL
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.selectedCellType = STONE
	}
}

func (g *Game) Update() error {
	g.HandleInput()
	g.sanbox.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		fmt.Println("init")
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.sanbox.Draw(g.pixels)

	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nType: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), g.selectedCellType))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		sanbox:           NewSandbox(screenWidth, screenHeight),
		selectedCellType: SAND,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sandbox")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
