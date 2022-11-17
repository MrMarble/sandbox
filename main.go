package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/constraints"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	pixels []byte
	sanbox *Sandbox

	selectedCellType CellType
	brushSize        int

	prevPos   [2]int
	cursorPos [2]int

	cellQueue [][2][2]int
}

func (g *Game) updateCursor() {
	g.prevPos = g.cursorPos
	x, y := ebiten.CursorPosition()
	g.cursorPos = [2]int{x, y}
}

func (g *Game) HandleInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.cellQueue = append(g.cellQueue, [2][2]int{g.prevPos, g.cursorPos})
	}

	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.brushSize = min(50, g.brushSize+2)
	} else if scrollY < 0 {
		g.brushSize = max(2, g.brushSize-2)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selectedCellType = SAND
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selectedCellType = WATER
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.selectedCellType = WALL
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.selectedCellType = STONE
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.selectedCellType = EMPTY
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.sanbox = NewSandbox(screenWidth, screenHeight)
	}
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func (g *Game) Update() error {
	g.updateCursor()
	g.HandleInput()
	g.placeParticles()

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
	rect(screen, g.cursorPos[0]-g.brushSize/2, g.cursorPos[1]-g.brushSize/2, g.brushSize, g.brushSize)
	// DEBUG
	for _, chunk := range g.sanbox.chunks {
		rect(screen, chunk.x*chunk.width, chunk.y*chunk.height, chunk.width, chunk.height)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nType: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), g.selectedCellType))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		sanbox:           NewSandbox(screenWidth, screenHeight),
		selectedCellType: SAND,
		brushSize:        10,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sandbox")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) placeParticles() {
	for {
		if len(g.cellQueue) == 0 {
			break
		}

		p1, p2 := g.cellQueue[0][0], g.cellQueue[0][1]
		g.cellQueue = g.cellQueue[1:]
		p1x, p1y := p1[0], p1[1]
		p2x, p2y := p2[0], p2[1]

		dx := int(math.Abs(float64(p2x - p1x)))
		sx := -1
		if p1x < p2x {
			sx = 1
		}
		dy := -int(math.Abs(float64(p2y - p1y)))
		sy := -1
		if p1y < p2y {
			sy = 1
		}
		err := dx + dy
		for {
			for x := p1x - (g.brushSize / 2); x < (p1x + g.brushSize/2); x++ {
				for y := p1y - (g.brushSize / 2); y < (p1y + g.brushSize/2); y++ {
					if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
						if g.selectedCellType == EMPTY || g.sanbox.IsEmpty(x, y) {
							variable := byte(rand.Intn(20))
							g.sanbox.SetCell(x, y, *NewCell(g.selectedCellType).WithVariation(variable))
						}
					}
				}
			}
			if p1x == p2x && p1y == p2y {
				break
			}
			e2 := 2 * err
			if e2 >= dy {
				err += dy
				p1x += sx
			}
			if e2 <= dx {
				err += dx
				p1y += sy
			}
		}
	}
}
