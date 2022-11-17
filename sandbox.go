package main

type Sandbox struct {
	width, height   int
	cWidth, cHeight int

	chunks      []*Chunk
	chunkLookup map[int]*Chunk
}

func hash(x, y int) int {
	return x*31 ^ y
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
	if chunk, ok := s.chunkLookup[hash(cx, cy)]; ok {
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
	s.chunkLookup[hash(x, y)] = chunk
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

func (s *Sandbox) RemoveEmptyChunks() {
	for i := 0; i < len(s.chunks); i++ {
		chunk := s.chunks[i]
		if chunk.filledCells == 0 {
			delete(s.chunkLookup, hash(chunk.x, chunk.y))
			s.chunks = append(s.chunks[:i], s.chunks[i+1:]...)
			i--

			chunk = nil
		}
	}
}

func (s *Sandbox) Update() {
	for _, chunk := range s.chunks {
		NewWorker(s, chunk).UpdateChunk()
	}

	for _, chunk := range s.chunks {
		chunk.ApplyChanges()
	}
	for _, chunk := range s.chunks {
		chunk.UpdateRect()
	}
}

func (s *Sandbox) KeepAlive(x, y int) {
	chunk := s.GetChunk(x, y)
	if chunk != nil {
		chunk.KeepAlive(x, y)
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
	// Remove chunks here to avoid leaving dangling particles
	s.RemoveEmptyChunks()

}
