package sandbox

import (
	"github.com/mrmarble/sandbox/pkg/misc"
	"pgregory.net/rand"
)

func (w *Worker) UpdateSteam(x, y int) {
	cell := w.GetCell(x, y)
	if cell.temp < 100 {
		w.SetCell(x, y, NewCell(WATER))
	}
}

func (w *Worker) UpdateReplicator(x, y int) {
	cell := w.GetCell(x, y)
	cell.extraData1 = 0

	if w.InBounds(x, y+2) {
		other := w.GetCell(x, y+1)
		if !isEmpty(other) && other.CType != CLNE {
			cell.extraData1 = 1
			if w.IsEmpty(x, y+2) {
				w.SetCell(x, y+2, NewCell(other.CType))
			}
		}
	}
	if w.InBounds(x, y-2) {
		other := w.GetCell(x, y-1)
		if !isEmpty(other) && other.CType != CLNE {
			cell.extraData1 = 1
			if w.IsEmpty(x, y-2) {
				w.SetCell(x, y-2, NewCell(other.CType))
			}
		}
	}
	if w.InBounds(x+2, y) {
		other := w.GetCell(x+2, y)
		if !isEmpty(other) && other.CType != CLNE {
			cell.extraData1 = 1
			if w.IsEmpty(x+2, y) {
				w.SetCell(x+2, y, NewCell(other.CType))
			}
		}
	}
	if w.InBounds(x-2, y) {
		other := w.GetCell(x-2, y)
		if !isEmpty(other) && other.CType != CLNE {
			cell.extraData1 = 1
			if w.IsEmpty(x-2, y) {
				w.SetCell(x-2, y, NewCell(other.CType))
			}
		}
	}
}

func (w *Worker) UpdateFire(x, y int) {
	cell := w.GetCell(x, y)
	if cell.temp < 40 || cell.extraData2 > 60 {
		if rand.Intn(10) > 1 {
			w.SetCell(x, y, nil)
		} else {
			smk := NewCell(SMOKE)
			smk.extraData1 = 0
			smk.extraData2 = 30 + (rand.Intn(30) + -15)
			w.SetCell(x, y, smk)
		}
	}
	if w.InBounds(x, y+1) {
		other := w.GetCell(x, y+1)
		if !isEmpty(other) && other.IsFlamable() && rand.Intn(100) > 65 {
			w.SetCell(x, y+1, NewCell(FIRE))
			w.SetCell(x, y, nil)
		}
	}
	if w.InBounds(x, y-1) {
		other := w.GetCell(x, y-1)
		if !isEmpty(other) && other.IsFlamable() && rand.Intn(100) > 65 {
			w.SetCell(x, y-1, NewCell(FIRE))
			w.SetCell(x, y, nil)
		}
	}
	if w.InBounds(x+1, y) {
		other := w.GetCell(x+1, y)
		if !isEmpty(other) && other.IsFlamable() && rand.Intn(100) > 65 {
			w.SetCell(x+1, y, NewCell(FIRE))
			w.SetCell(x, y, nil)
		}
	}
	if w.InBounds(x-1, y) {
		other := w.GetCell(x-1, y)
		if !isEmpty(other) && other.IsFlamable() && rand.Intn(100) > 65 {
			w.SetCell(x-1, y, NewCell(FIRE))
			w.SetCell(x, y, nil)
		}
	}
}

func (w *Worker) UpdateSmoke(x, y int) {
	cell := w.GetCell(x, y)

	if cell.extraData1 == 0 {
		cell.extraData2--
		if cell.extraData2 == 0 {
			w.SetCell(x, y, nil)
		}
	} else {
		cell.extraData1--
	}
}

func (w *Worker) UpdateWater(x, y int) {
	cell := w.GetCell(x, y)

	if cell.temp >= 100 {
		t := misc.Clamp(float64(cell.temp/150), 0, 1)
		chance := (1-t)*0.3 + t*0.7
		if rand.Intn(100) > int(chance) {
			w.SetCell(x, y, NewCell(STEAM))
		}
	}

	y2 := y + 1
	for {
		if !w.InBounds(x, y2) {
			return
		}
		cellBelow := w.GetCell(x, y2)
		if !isEmpty(cellBelow) && cellBelow.CType == SAND {
			if cellBelow.extraData1 == 0 {
				w.SetCell(x, y, nil)
				cellBelow.extraData1 = 1
				return
			}
		} else {
			return
		}
		y2++
	}
}

func (w *Worker) MoveFire(x, y int, cell *Cell) {
	nx, ny := w.MoveGas(x, y, cell)
	if other := w.GetCell(nx, ny); !isEmpty(other) {
		extraData2 := other.extraData2
		if ny == y {
			extraData2++
		} else {
			extraData2--
		}
		cell.extraData2 = extraData2
	}
}

func (w *Worker) UpdateSand(x, y int) {
	cell := w.GetCell(x, y)
	if !isEmpty(cell) && cell.extraData1 == 1 && cell.temp >= 30 {
		cell.extraData1 = 0
	}

	if cell.temp >= 120 {
		w.SetCell(x, y, NewCell(GLASS))
	}
}

func (w *Worker) MovePowder(x, y int, cell *Cell) {
	if w.IsEmpty(x, y+1) {
		w.MoveCell(x, y, x, y+1)
	} else {
		xn, yn := w.randomNeighbour(x, y, 1)
		if xn != -1 && yn != -1 {
			w.MoveCell(x, y, xn, yn)
		}
	}
}

func (w *Worker) MoveLiquid(x, y int, cell *Cell) {
	if w.IsEmpty(x, y+1) {
		w.MoveCell(x, y, x, y+1)
	} else {
		xn, yn := w.randomNeighbour(x, y, 1)
		if xn != -1 && yn != -1 {
			w.MoveCell(x, y, xn, yn)
		} else {
			xn, yn = w.randomNeighbour(x, y, 0)
			if xn != -1 && yn != -1 {
				w.MoveCell(x, y, xn, yn)
			}
		}
	}
}

// TODO: refactor this to return position.
func (w *Worker) MoveGas(x, y int, cell *Cell) (int, int) {
	if w.IsEmpty(x, y-1) && rand.Intn(100) < 50 {
		w.MoveCell(x, y, x, y-1)
		return x, y - 1
	}

	xn, yn := w.randomNeighbour(x, y, -1)
	if xn != -1 && yn != -1 {
		w.MoveCell(x, y, xn, yn)
		return xn, yn
	}

	xn, yn = w.randomNeighbour(x, y, 0)
	if xn != -1 && yn != -1 {
		w.MoveCell(x, y, xn, yn)
		return xn, yn
	}

	return x, y
}

func (w *Worker) MoveSolid(x, y int, cell *Cell) {
	if w.IsEmpty(x, y+1) {
		w.MoveCell(x, y, x, y+1)
	}
}

func (w *Worker) randomNeighbour(x, y, yOffset int) (int, int) {
	// check if left is air
	leftFree := false
	if w.InBounds(x-1, y+yOffset) && w.IsEmpty(x-1, y) && w.IsEmpty(x-1, y+yOffset) {
		leftFree = true
	}

	// check if right is air
	rightFree := false
	if w.InBounds(x+1, y+yOffset); w.IsEmpty(x+1, y) && w.IsEmpty(x+1, y+yOffset) {
		rightFree = true
	}

	if leftFree || rightFree {
		if leftFree && rightFree {
			if rand.Intn(2) == 1 {
				return x - 1, y + yOffset
			}
			return x + 1, y + yOffset
		} else if leftFree {
			return x - 1, y + yOffset
		} else {
			return x + 1, y + yOffset
		}
	}
	return -1, -1
}
