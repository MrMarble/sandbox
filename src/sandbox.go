package main

import (
	"strconv"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Sandbox struct {
	width, height   int
	cWidth, cHeight int

	chunks      []*Chunk
	chunkLookup cmap.ConcurrentMap[string, *Chunk]

	chunkMutex sync.Mutex
}

func hash(x, y int) string {
	return strconv.Itoa(x*31 ^ y)
}

const MaxChunks = 8

func NewSandbox(width, height int) *Sandbox {
	cWidth := width / MaxChunks
	cHeight := height / MaxChunks
	return &Sandbox{
		width:       width,
		height:      height,
		cWidth:      cWidth,
		cHeight:     cHeight,
		chunks:      []*Chunk{},
		chunkLookup: cmap.New[*Chunk](),
	}
}

func (s *Sandbox) GetChunk(x, y int) *Chunk {
	cx, cy := s.GetChunkLocation(x, y)
	if chunk, ok := s.chunkLookup.Get(hash(cx, cy)); ok {
		return chunk
	}
	return s.CreateChunk(cx, cy)
}

func (s *Sandbox) CreateChunk(x, y int) *Chunk {
	if x < 0 || y < 0 || x >= MaxChunks || y >= MaxChunks {
		return nil
	}
	chunk := NewChunk(s.cWidth, s.cHeight, x, y)
	s.chunkMutex.Lock()
	s.chunks = append(s.chunks, chunk)
	s.chunkMutex.Unlock()
	s.chunkLookup.Set(hash(x, y), chunk)
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
			s.chunkLookup.Remove(hash(chunk.x, chunk.y))
			s.chunks = append(s.chunks[:i], s.chunks[i+1:]...)
			i--

			chunk = nil
		}
	}
}

func (s *Sandbox) MoveUpdate() {
	var wg sync.WaitGroup
	for _, chunk := range s.chunks {
		wg.Add(1)
		go func(s *Sandbox, c *Chunk) {
			NewWorker(s, c).UpdateChunk()
			wg.Done()
		}(s, chunk)
	}
	wg.Wait()

	for _, chunk := range s.chunks {
		wg.Add(1)
		go func(c *Chunk) {
			c.ApplyChanges()
			wg.Done()
		}(chunk)
	}
	wg.Wait()

	for _, chunk := range s.chunks {
		chunk.UpdateRect()
	}
}

func (s *Sandbox) StateUpdate() {
	var wg sync.WaitGroup
	for _, chunk := range s.chunks {
		wg.Add(1)
		go func(s *Sandbox, c *Chunk) {
			NewWorker(s, c).UpdateChunkState()
			wg.Done()
		}(s, chunk)
	}
	wg.Wait()
}

func (s *Sandbox) Update() {
	s.MoveUpdate()
	s.StateUpdate()
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
			if cell.cType == AIR {
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
