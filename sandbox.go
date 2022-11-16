package main

import (
	"math/rand"
	"sort"
)

type Sandbox struct {
	width, height int
	cells         []Cell
	changes       [][2]int // dst, src
}

func NewSandbox(width, height int) *Sandbox {
	return &Sandbox{
		width:  width,
		height: height,
		cells:  make([]Cell, width*height),
	}
}

// GetIndex returns the index of the cell at the given coordinates.
func (s *Sandbox) GetIndex(x, y int) int {
	return x + y*s.width
}

func (s *Sandbox) GetCell(x, y int) *Cell {
	return s.GetCellAt(s.GetIndex(x, y))
}

func (s *Sandbox) GetCellAt(i int) *Cell {
	return &s.cells[i]
}

func (s *Sandbox) InBounds(x, y int) bool {
	return x >= 0 && x < s.width && y >= 0 && y < s.height
}

func (s *Sandbox) IsEmpty(x, y int) bool {
	return s.InBounds(x, y) && s.GetCell(x, y).cType == EMPTY
}

func (s *Sandbox) SetCell(x, y int, c Cell) {
	s.cells[s.GetIndex(x, y)] = c
}

func (s *Sandbox) MoveCell(x, y, dx, dy int) {
	s.changes = append(s.changes, [2]int{s.GetIndex(dx, dy), s.GetIndex(x, y)})
}

func (s *Sandbox) ApplyChanges() {
	// remove changes that have the destination cell occupied
	for i := 0; i < len(s.changes); i++ {
		if s.cells[s.changes[i][0]].cType != EMPTY {
			s.changes = append(s.changes[:i], s.changes[i+1:]...)
			i--
		}
	}

	// sort changes by destination index
	sort.Slice(s.changes, func(i, j int) bool {
		return s.changes[i][0] < s.changes[j][0]
	})

	// pick random source for each destination
	iPrev := 0
	s.changes = append(s.changes, [2]int{-1, -1}) // catch the last one
	for i := 0; i < len(s.changes)-1; i++ {
		if s.changes[i+1][0] != s.changes[i][0] {
			rng := rand.Intn(i-iPrev+1) + iPrev

			dst := s.changes[rng][0]
			src := s.changes[rng][1]

			s.cells[dst] = s.cells[src]
			s.cells[src] = Cell{}

			iPrev = i + 1
		}
	}
	s.changes = [][2]int{}
}

func (s *Sandbox) Update() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			c := s.GetCell(x, y)
			if c.cType == EMPTY {
				continue
			}
			switch c.cType {
			case SAND:
				s.MovePowder(x, y, c)
			case WATER:
				s.MoveLiquid(x, y, c)
			case STONE:
				s.MoveSolid(x, y, c)
			}
		}
	}
	s.ApplyChanges()
}

func (s *Sandbox) Draw(pix []byte) {
	for i, c := range s.cells {
		if c.cType == EMPTY {
			pix[i*4] = 0
			pix[i*4+1] = 0
			pix[i*4+2] = 0
			pix[i*4+3] = 0
			continue
		}
		pix[i*4] = c.color.R
		pix[i*4+1] = c.color.G
		pix[i*4+2] = c.color.B
		pix[i*4+3] = c.color.A
	}
}

func (s *Sandbox) MovePowder(x, y int, cell *Cell) {
	if s.IsEmpty(x, y+1) {
		s.MoveCell(x, y, x, y+1)
	} else {
		xn, yn := s.randomNeighbour(x, y, 1)
		if xn != -1 && yn != -1 {
			s.MoveCell(x, y, xn, yn)
		}
	}
}

func (s *Sandbox) MoveLiquid(x, y int, cell *Cell) {
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

func (s *Sandbox) MoveSolid(x, y int, cell *Cell) {
	if s.IsEmpty(x, y+1) {
		s.MoveCell(x, y, x, y+1)
	}
}

func (s *Sandbox) randomNeighbour(x, y, yOffset int) (int, int) {

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
