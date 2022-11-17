package main

import "math/rand"

type Worker struct {
	chunk   *Chunk
	sandbox *Sandbox
}

func NewWorker(sandbox *Sandbox, chunk *Chunk) *Worker {
	return &Worker{
		sandbox: sandbox,
		chunk:   chunk,
	}
}

func (w *Worker) InBounds(x, y int) bool {
	return w.chunk.InBounds(x, y) || w.sandbox.InBounds(x, y)
}

func (w *Worker) IsEmpty(x, y int) bool {
	if w.chunk.InBounds(x, y) {
		return w.chunk.IsEmpty(x, y)
	}
	return w.sandbox.IsEmpty(x, y)
}

func (w *Worker) GetCell(x, y int) *Cell {
	if w.chunk.InBounds(x, y) {
		return w.chunk.GetCell(x, y)
	}
	return w.sandbox.GetCell(x, y)
}

func (w *Worker) SetCell(x, y int, cell Cell) {
	if w.chunk.InBounds(x, y) {
		w.chunk.SetCell(x, y, cell)
	} else {
		w.sandbox.SetCell(x, y, cell)
	}
}

func (w *Worker) MoveCell(x, y, dx, dy int) {
	if w.chunk.InBounds(x, y) && w.chunk.InBounds(dx, dy) {
		w.chunk.MoveCell(w.chunk, x, y, dx, dy)
	} else {
		w.sandbox.MoveCell(x, y, dx, dy)
	}
}

func (w *Worker) UpdateChunk() {
	for y := 0; y < w.chunk.height; y++ {
		for x := 0; x < w.chunk.width; x++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.width)
			px := x + w.chunk.x*w.chunk.width
			py := y + w.chunk.y*w.chunk.height
			if c.cType == EMPTY {
				continue
			}
			switch c.cType {
			case SAND:
				w.MovePowder(px, py, c)
			case WATER:
				w.MoveLiquid(px, py, c)
			case STONE:
				w.MoveSolid(px, py, c)
			}
		}
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
