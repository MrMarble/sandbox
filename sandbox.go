package main

import (
	"math/rand"
)

type Sandbox struct {
	width, height   int
	cWidth, cHeight int

	chunks      []*Chunk
	chunkLookup map[int]*Chunk
}

func NewSandbox(width, height int) *Sandbox {
	cWidth := width / 4
	cHeight := height / 4
	return &Sandbox{
		width:       width,
		height:      height,
		cWidth:      cWidth,
		cHeight:     cHeight,
		chunks:      []*Chunk{},
		chunkLookup: map[int]*Chunk{},
	}
}

func (s *Sandbox) GetChunk(x, y int) *Chunk {
	cx, cy := s.GetChunkLocation(x, y)
	if chunk, ok := s.chunkLookup[cy*s.cWidth+cx]; ok {
		return chunk
	}
	return s.CreateChunk(cx, cy)
}

const MaxChunks = 4

func (s *Sandbox) CreateChunk(x, y int) *Chunk {
	if x < 0 || y < 0 || x >= MaxChunks || y >= MaxChunks {
		return nil
	}
	chunk := NewChunk(s.cWidth, s.cHeight, x, y)
	s.chunks = append(s.chunks, chunk)
	s.chunkLookup[y*s.cWidth+x] = chunk
	return chunk
}

func (s *Sandbox) GetChunkLocation(x, y int) (int, int) {
	return x / s.cWidth, y / s.cHeight
}

func (s *Sandbox) InBounds(x, y int) bool {
	chunk := s.GetChunk(x, y)
	if chunk != nil {
		return chunk.InBounds(x, y)
	}
	return false
}

func (s *Sandbox) IsEmpty(x, y int) bool {
	return s.InBounds(x, y) && s.GetChunk(x, y).IsEmpty(x, y)
}

func (s *Sandbox) GetCell(x, y int) *Cell {
	chunk := s.GetChunk(x, y)
	if chunk != nil {
		return chunk.GetCell(x, y)
	}
	return nil
}

func (s *Sandbox) SetCell(x, y int, cell Cell) {
	chunk := s.GetChunk(x, y)
	if chunk != nil {
		chunk.SetCell(x, y, cell)
	}
}

func (s *Sandbox) MoveCell(x, y, xn, yn int) {
	src := s.GetChunk(x, y)
	dst := s.GetChunk(xn, yn)
	if src != nil && dst != nil {
		dst.MoveCell(src, x, y, xn, yn)
	}
}

func (s *Sandbox) Update() {
	for _, chunk := range s.chunks {
		for y := 0; y < chunk.height; y++ {
			for x := 0; x < chunk.width; x++ {
				c := chunk.GetCellAt(x + y*chunk.width)
				px := x + chunk.x*chunk.width
				py := y + chunk.y*chunk.height
				if c.cType == EMPTY {
					continue
				}
				switch c.cType {
				case SAND:
					s.MovePowder(px, py, c)
				case WATER:
					s.MoveLiquid(px, py, c)
				case STONE:
					s.MoveSolid(px, py, c)
				}
			}
		}
	}

	for _, chunk := range s.chunks {
		chunk.ApplyChanges()
	}
}

func (s *Sandbox) Draw(pix []byte) {
	for _, c := range s.chunks {
		for i, cell := range c.cells {
			x := i%c.width + c.x*c.width
			y := i/c.width + c.y*c.height
			idx := (x + y*screenWidth)
			if cell.cType == EMPTY {
				pix[idx*4] = 0
				pix[idx*4+1] = 0
				pix[idx*4+2] = 0
				pix[idx*4+3] = 0
				continue
			}
			pix[idx*4] = cell.color.R
			pix[idx*4+1] = cell.color.G
			pix[idx*4+2] = cell.color.B
			pix[idx*4+3] = cell.color.A
		}
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
