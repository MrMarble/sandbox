package main

import (
	"fmt"
	"image/color"
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
	pixels  []byte
	sandbox *Sandbox

	pause bool
	debug bool

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

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.pause = !g.pause
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.sandbox = NewSandbox(screenWidth, screenHeight)
		g.pixels = make([]byte, screenWidth*screenHeight*4)
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

func clamp[T constraints.Ordered](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func (g *Game) Update() error {
	g.updateCursor()
	g.HandleInput()
	g.placeParticles()

	if !g.pause {
		g.sandbox.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		fmt.Println("init")
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.sandbox.Draw(g.pixels)
	screen.WritePixels(g.pixels)
	rect(screen, g.cursorPos[0]-g.brushSize/2, g.cursorPos[1]-g.brushSize/2, g.brushSize, g.brushSize, color.White)
	// DEBUG
	if g.debug {
		for _, chunk := range g.sandbox.chunks {
			rect(screen, chunk.x*chunk.width, chunk.y*chunk.height, chunk.width, chunk.height, color.RGBA{100, 0, 0, 100})
			if chunk.maxX > 0 {
				rect(screen, chunk.x*chunk.width+chunk.minX, chunk.y*chunk.height+chunk.minY, chunk.maxX-chunk.minX, chunk.maxY-chunk.minY, color.RGBA{0, 100, 0, 100})

			}
		}
	}

	dbg := fmt.Sprintf("Brush size: %d\n", g.brushSize)
	dbg += fmt.Sprintf("Cell type: %d\n", g.selectedCellType)
	dbg += fmt.Sprintf("Pause: %v\n", g.pause)
	if g.debug {
		dbg += fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS())
		dbg += fmt.Sprintf("FPS: %0.2f\n", ebiten.CurrentFPS())
		dbg += fmt.Sprintf("Chunks: %d\n", len(g.sandbox.chunks))
		dbg += fmt.Sprintf("ChunksL: %d\n", g.sandbox.chunkLookup.Count())
	}
	ebitenutil.DebugPrint(screen, dbg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		sandbox:          NewSandbox(screenWidth, screenHeight),
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
						if g.selectedCellType == EMPTY || g.sandbox.IsEmpty(x, y) {
							variable := byte(rand.Intn(20))
							g.sandbox.SetCell(x, y, *NewCell(g.selectedCellType).WithVariation(variable))
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