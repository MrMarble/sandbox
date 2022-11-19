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
	margin       = 50

	offscrenOptions = &ebiten.DrawImageOptions{}
)

type game struct {
	pixels  []byte
	sandbox *sandbox.Sandbox
	menu    *ui.Menu

	pause       bool
	debug       bool
	tempOverlay bool

	brushSize int

	prevPos   [2]int
	cursorPos [2]int

	cellQueue [][2][2]int

	offscreen *ebiten.Image
}

func New() *game {
	ebiten.SetWindowTitle("Sandbox")
	ebiten.SetWindowResizable(false)
	ebiten.SetMaxTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	return &game{
		sandbox:     sandbox.NewSandbox(screenWidth-margin, screenHeight-margin),
		brushSize:   10,
		tempOverlay: true,
		offscreen:   ebiten.NewImage(screenWidth-margin, screenHeight-margin),
		menu:        ui.NewMenu(margin/2, screenHeight-margin/2+5),
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	return screenWidth * int(s), screenHeight * int(s)
}

func (g *game) Update() error {
	g.updateCursor()
	g.handleInput()
	g.menu.Update()
	g.placeQueueParticles()

	if !g.pause {
		g.sandbox.Update(g.tempOverlay)
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		fmt.Println("init")
		g.pixels = make([]byte, g.offscreen.Bounds().Dx()*g.offscreen.Bounds().Dy()*4)
		offscrenOptions.GeoM.Translate(float64(margin/2), float64(margin/2))
	}

	g.sandbox.Draw(g.pixels, g.offscreen.Bounds().Dx(), g.tempOverlay)
	g.offscreen.WritePixels(g.pixels)

	// Brush size
	offscreenX, offscreenY := offscreenCursor(g.cursorPos[0], g.cursorPos[1])
	ui.Rect(g.offscreen, offscreenX-g.brushSize/2, offscreenY-g.brushSize/2, g.brushSize, g.brushSize, color.White, false)
	// Border
	ui.Rect(g.offscreen, 0, 0, screenWidth-margin, screenHeight-margin, color.RGBA{20, 20, 20, 100}, false)

	g.menu.Draw(screen)
	g.debugInfo(screen)
	screen.DrawImage(g.offscreen, offscrenOptions)

}

func (g *game) updateCursor() {
	g.prevPos = g.cursorPos
	x, y := ebiten.CursorPosition()
	g.cursorPos = [2]int{x, y}
}

func (g *game) handleInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		prevX, prevY := offscreenCursor(g.prevPos[0], g.prevPos[1])
		x, y := offscreenCursor(g.cursorPos[0], g.cursorPos[1])
		g.cellQueue = append(g.cellQueue, [2][2]int{{prevX, prevY}, {x, y}})
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

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.tempOverlay = !g.tempOverlay
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.sandbox = sandbox.NewSandbox(screenWidth-margin, screenHeight-margin)
		g.pixels = nil
		offscrenOptions.GeoM.Reset()
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
						if g.menu.GetSelected() == sandbox.AIR || g.sandbox.IsEmpty(x, y) {
							g.sandbox.SetCell(x, y, sandbox.NewCell(g.menu.GetSelected()))
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
		curx, cury := offscreenCursor(g.cursorPos[0], g.cursorPos[1])
		if g.sandbox.InBounds(curx, cury) {
			cell := g.sandbox.GetCell(curx, cury)
			if cell != nil {
				dbg += fmt.Sprintf("Particle: %+v\n", cell)
			}
		}
		dbg += fmt.Sprintf("X: %d Y: %d\n", curx, cury)
		for _, chunk := range g.sandbox.Chunks {
			ui.Rect(g.offscreen, chunk.X*chunk.Width, chunk.Y*chunk.Height, chunk.Width, chunk.Height, color.RGBA{100, 0, 0, 100}, false)
			text.Draw(g.offscreen, fmt.Sprintf("%d,%d", chunk.X, chunk.Y), bitmapfont.Gothic12r, chunk.X*chunk.Width+12, chunk.Y*chunk.Height+12, color.White)
			if chunk.MaxX > 0 {
				ui.Rect(g.offscreen, chunk.X*chunk.Width+chunk.MinX, chunk.Y*chunk.Height+chunk.MinY, chunk.MaxX-chunk.MinX, chunk.MaxY-chunk.MinY, color.RGBA{0, 100, 0, 100}, false)
			}
		}
	}
	ebitenutil.DebugPrint(screen, dbg)

}

func offscreenCursor(x, y int) (int, int) {
	return x - margin/2, y - margin/2
}
