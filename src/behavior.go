package main

import "math/rand"

func (s *Worker) UpdateSteam(x, y int) {
	cell := s.GetCell(x, y)
	if cell.temp < 100 {
		s.SetCell(x, y, *NewCell(WATR))
	}
}

func (s *Worker) UpdateSmoke(x, y int) {
	cell := s.GetCell(x, y)
	if cell.extraData1 == 0 {
		cell.extraData2--
		if cell.extraData2 == 0 {
			s.SetCell(x, y, Cell{cType: AIR})
		}
	} else {
		cell.extraData1--
	}
}

func (s *Worker) UpdateWater(x, y int) {
	cell := s.GetCell(x, y)
	if cell.temp >= 100 {
		t := clamp(float64(cell.temp/150), 0, 1)
		chance := (1-t)*0.3 + t*0.7
		if rand.Intn(100) < int(100-chance) {
			s.SetCell(x, y, *NewCell(STEM))
		}
	}

	y2 := y + 1
	for {
		if !s.InBounds(x, y2) {
			return
		}
		cellBelow := s.GetCell(x, y2)
		if cellBelow.cType == SAND {
			if cellBelow.extraData1 == 0 {
				s.SetCell(x, y, Cell{cType: AIR})
				cellBelow.extraData1 = 1
				cellBelow.color = ParseHexColor("#b19d5e")
				return
			}
		} else {
			return
		}
		y2++
	}

}

func (s *Worker) UpdateSand(x, y int) {
	cell := s.GetCell(x, y)
	if cell.extraData1 == 1 && cell.temp >= 30 {
		cell.extraData1 = 0
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

func (s *Worker) MoveGas(x, y int, cell *Cell) {
	if s.IsEmpty(x, y-1) && rand.Intn(100) < 50 {
		s.MoveCell(x, y, x, y-1)
	} else {
		xn, yn := s.randomNeighbour(x, y, -1)
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
