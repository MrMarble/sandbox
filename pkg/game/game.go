package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mrmarble/sandbox/pkg/misc"
	"github.com/mrmarble/sandbox/pkg/sandbox"
	"github.com/mrmarble/sandbox/pkg/ui"
)

var (
	screenWidth  = 640
	screenHeight = 480
)

type game struct {
	pixels  []byte
	sandbox *sandbox.Sandbox

	pause bool
	debug bool

	selectedCellType sandbox.CellType
	brushSize        int

	prevPos   [2]int
	cursorPos [2]int

	cellQueue [][2][2]int
}

func New() *game {
	ebiten.SetWindowTitle("Sandbox")
	ebiten.SetWindowResizable(false)
	ebiten.SetMaxTPS(120)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	return &game{
		sandbox:          sandbox.NewSandbox(screenWidth, screenHeight),
		selectedCellType: sandbox.SAND,
		brushSize:        10,
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	return screenWidth * int(s), screenHeight * int(s)
}

func (g *game) Update() error {
	g.updateCursor()
	g.handleInput()
	g.placeQueueParticles()

	g.sandbox.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		fmt.Println("init")
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.sandbox.Draw(g.pixels, screenWidth)
	screen.WritePixels(g.pixels)
	g.debugInfo(screen)
	g.DrawUI(screen)
}

func (g *game) updateCursor() {
	g.prevPos = g.cursorPos
	x, y := ebiten.CursorPosition()
	g.cursorPos = [2]int{x, y}
}

func (g *game) handleInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.cellQueue = append(g.cellQueue, [2][2]int{g.prevPos, g.cursorPos})
	}

	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.brushSize = misc.Min(50, g.brushSize+2)
	} else if scrollY < 0 {
		g.brushSize = misc.Max(2, g.brushSize-2)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.pause = !g.pause
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.sandbox = sandbox.NewSandbox(screenWidth, screenHeight)
		g.pixels = nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.cursorPos[1] > screenHeight-20 {
			idx := int(math.Floor(float64(g.cursorPos[0]) / 30))
			if idx < 0 || idx > int(sandbox.AIR) {
				return
			}
			g.selectedCellType = sandbox.CellType(idx)
		}
	}
}

func (g *game) placeQueueParticles() {
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
						if g.selectedCellType == sandbox.AIR || g.sandbox.IsEmpty(x, y) {
							g.sandbox.SetCell(x, y, sandbox.NewCell(g.selectedCellType))
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

func (g *game) debugInfo(screen *ebiten.Image) {
	dbg := fmt.Sprintf("FPS: %0.2f\n", ebiten.CurrentFPS())

	if g.debug {
		dbg += fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS())
		//dbg += fmt.Sprintf("Chunks: %d\n", len(g.sandbox.chunks))
		//dbg += fmt.Sprintf("ChunksL: %d\n", g.sandbox.chunkLookup.Count())
		/*cell := g.sandbox.GetCell(g.cursorPos[0], g.cursorPos[1])
		if cell != nil {
			dbg += fmt.Sprintf("Particle: %s\n", cell.cType)
		}*/
		dbg += fmt.Sprintf("X: %d Y: %d\n", g.cursorPos[0], g.cursorPos[1])
		for _, chunk := range g.sandbox.Chunks {
			ui.Rect(screen, chunk.X*chunk.Width, chunk.Y*chunk.Height, chunk.Width, chunk.Height, color.RGBA{100, 0, 0, 100}, false)
			text.Draw(screen, fmt.Sprintf("%d,%d", chunk.X, chunk.Y), bitmapfont.Gothic12r, chunk.X*chunk.Width+12, chunk.Y*chunk.Height+12, color.White)
			if chunk.MaxX > 0 {
				ui.Rect(screen, chunk.X*chunk.Width+chunk.MinX, chunk.Y*chunk.Height+chunk.MinY, chunk.MaxX-chunk.MinX, chunk.MaxY-chunk.MinY, color.RGBA{0, 100, 0, 100}, false)
			}
		}
	}
	ebitenutil.DebugPrint(screen, dbg)

}

func (g *game) DrawUI(screen *ebiten.Image) {
	ui.Rect(screen, g.cursorPos[0]-g.brushSize/2, g.cursorPos[1]-g.brushSize/2, g.brushSize, g.brushSize, color.White, false)

	for i := 0; i <= int(sandbox.AIR); i++ {
		ui.Button(screen, sandbox.CellType(i).String(), 5+30*i, screenHeight-18, sandbox.CellType(i).Color(), g.selectedCellType == sandbox.CellType(i))
	}
}
