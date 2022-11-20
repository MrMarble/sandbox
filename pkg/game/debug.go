package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mrmarble/sandbox/pkg/ui"
)

func (g *Game) debugInfo(screen *ebiten.Image) {
	dbg := fmt.Sprintf("FPS: %0.2f\n", ebiten.ActualFPS())

	if g.debug {
		dbg += fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS())
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
