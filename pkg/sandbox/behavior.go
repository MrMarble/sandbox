package sandbox

import (
	"github.com/mrmarble/sandbox/pkg/misc"
	"pgregory.net/rand"
)

func (s *Worker) UpdateSteam(x, y int) {
	cell := s.GetCell(x, y)
	if cell.temp < 100 {
		s.SetCell(x, y, NewCell(WATER))
	}
}

func (w *Worker) UpdateFire(x, y int) {
	cell := w.GetCell(x, y)
	if cell.temp < 40 || cell.extraData2 > 60 {
		w.SetCell(x, y, nil)
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

func (s *Worker) UpdateSmoke(x, y int) {
	cell := s.GetCell(x, y)

	if cell.extraData1 == 0 {
		cell.extraData2--
		if cell.extraData2 == 0 {
			s.SetCell(x, y, nil)
		}
	} else {
		cell.extraData1--
	}
}

func (s *Worker) UpdateWater(x, y int) {
	cell := s.GetCell(x, y)

	if cell.temp >= 100 {
		t := misc.Clamp(float64(cell.temp/150), 0, 1)
		chance := (1-t)*0.3 + t*0.7
		if rand.Intn(100) > int(chance) {
			s.SetCell(x, y, NewCell(STEAM))
		}
	}

	y2 := y + 1
	for {
		if !s.InBounds(x, y2) {
			return
		}
		cellBelow := s.GetCell(x, y2)
		if !isEmpty(cellBelow) && cellBelow.CType == SAND {
			if cellBelow.extraData1 == 0 {
				s.SetCell(x, y, nil)
				cellBelow.extraData1 = 1
				return
			}
		} else {
			return
		}
		y2++
	}

}

func (s *Worker) MoveFire(x, y int, cell *Cell) {
	nx, ny := s.MoveGas(x, y, cell)
	if other := s.GetCell(nx, ny); !isEmpty(other) {
		extraData2 := other.extraData2
		if ny == y {
			extraData2++
		} else {
			extraData2--
		}
		cell.extraData2 = extraData2
	}
}

func (s *Worker) UpdateSand(x, y int) {
	cell := s.GetCell(x, y)
	if !isEmpty(cell) && cell.extraData1 == 1 && cell.temp >= 30 {
		cell.extraData1 = 0
	}

	if cell.temp >= 120 {
		s.SetCell(x, y, NewCell(GLASS))
	}
}

func (s *Worker) MovePowder(x, y int, cell *Cell) {
	if s.IsEmpty(x, y+1) {
		s.MoveCell(x, y, x, y+1)
	} else {
		xn, yn := s.randomNeighbour(x, y, 1)
		if xn != -1 && yn != -1 {
			s.MoveCell(x, y, xn, yn)
		}
	}
}

func (s *Worker) MoveLiquid(x, y int, cell *Cell) {
	if s.IsEmpty(x, y+1) {
		s.MoveCell(x, y, x, y+1)
	} else {
		xn, yn := s.randomNeighbour(x, y, 1)
		if xn != -1 && yn != -1 {
			s.MoveCell(x, y, xn, yn)
		} else {
			xn, yn = s.randomNeighbour(x, y, 0)
			if xn != -1 && yn != -1 {
				s.MoveCell(x, y, xn, yn)
			}
		}
	}
}

// TODO: refactor this to return position
func (s *Worker) MoveGas(x, y int, cell *Cell) (int, int) {
	if s.IsEmpty(x, y-1) && rand.Intn(100) < 50 {
		s.MoveCell(x, y, x, y-1)
		return x, y - 1
	} else {
		xn, yn := s.randomNeighbour(x, y, -1)
		if xn != -1 && yn != -1 {
			s.MoveCell(x, y, xn, yn)
			return xn, yn
		} else {
			xn, yn = s.randomNeighbour(x, y, 0)
			if xn != -1 && yn != -1 {
				s.MoveCell(x, y, xn, yn)
				return xn, yn
			}
		}
	}
	return x, y
}

func (s *Worker) MoveSolid(x, y int, cell *Cell) {
	if s.IsEmpty(x, y+1) {
		s.MoveCell(x, y, x, y+1)
	}
}

func (s *Worker) randomNeighbour(x, y, yOffset int) (int, int) {

	// check if left is air
	leftFree := false
	if s.InBounds(x-1, y+yOffset) && s.IsEmpty(x-1, y) && s.IsEmpty(x-1, y+yOffset) {
		leftFree = true
	}

	// check if right is air
	rightFree := false
	if s.InBounds(x+1, y+yOffset); s.IsEmpty(x+1, y) && s.IsEmpty(x+1, y+yOffset) {
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
