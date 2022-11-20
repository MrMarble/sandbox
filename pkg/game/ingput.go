package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mrmarble/sandbox/pkg/misc"
	"github.com/mrmarble/sandbox/pkg/sandbox"
)

func (g *Game) updateCursor() {
	g.prevPos = g.cursorPos
	x, y := ebiten.CursorPosition()
	g.cursorPos = [2]int{x, y}
}

func (g *Game) handleInput() {
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

func (g *Game) placeQueueParticles() {
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
