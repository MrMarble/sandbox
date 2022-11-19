package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mrmarble/sandbox/pkg/sandbox"
)

type Menu struct {
	x, y             int
	selectedCellType sandbox.CellType
}

func NewMenu(x, y int) *Menu {
	return &Menu{
		x:                x,
		y:                y,
		selectedCellType: sandbox.SAND,
	}
}

func (m *Menu) GetSelected() sandbox.CellType {
	return m.selectedCellType
}

func (m *Menu) Draw(screen *ebiten.Image) {
	for i := 0; i <= int(sandbox.AIR); i++ {
		Button(screen, sandbox.CellType(i).String(), m.x+35*i, m.y, sandbox.CellType(i).Color(), m.selectedCellType == sandbox.CellType(i))
	}
}

func (m *Menu) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		curX, curY := ebiten.CursorPosition()
		if curY > m.y && curX > m.x {
			x := curX - m.x
			idx := int(math.Floor(float64(x) / 35))
			if idx < 0 || idx > int(sandbox.AIR) {
				return
			}
			m.selectedCellType = sandbox.CellType(idx)
		}
	}
}
