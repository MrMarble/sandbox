package sandbox

import (
	"strconv"
	"sync"

	"github.com/mrmarble/sandbox/pkg/misc"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Sandbox struct {
	width, height   int
	cWidth, cHeight int

	Chunks      []*Chunk
	chunkLookup cmap.ConcurrentMap[string, *Chunk]

	chunkMutex sync.Mutex
}

func hash(x, y int) string {
	return strconv.Itoa(x*0xf1f1f1f1 ^ y)
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
		Chunks:      []*Chunk{},
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
	s.Chunks = append(s.Chunks, chunk)
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

func (s *Sandbox) SetCell(x, y int, cell *Cell) {
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
	for i := 0; i < len(s.Chunks); i++ {
		chunk := s.Chunks[i]
		if chunk.filledCells == 0 {
			s.chunkLookup.Remove(hash(chunk.X, chunk.Y))
			s.chunkMutex.Lock()
			s.Chunks = append(s.Chunks[:i], s.Chunks[i+1:]...)
			s.chunkMutex.Unlock()
			i--

			chunk = nil
		} else {
			if chunk != nil {
				if !s.chunkLookup.Has(hash(chunk.X, chunk.Y)) {
					s.chunkMutex.Lock()
					s.Chunks = append(s.Chunks[:i], s.Chunks[i+1:]...)
					s.chunkMutex.Unlock()
				}
			}
		}

	}
}

func (s *Sandbox) MoveUpdate() {
	var wg sync.WaitGroup
	for _, chunk := range s.Chunks {
		wg.Add(1)
		go func(s *Sandbox, c *Chunk) {
			NewWorker(s, c).UpdateChunk()
			wg.Done()
		}(s, chunk)
	}
	wg.Wait()

	for _, chunk := range s.Chunks {
		wg.Add(1)
		go func(c *Chunk) {
			c.ApplyChanges()
			wg.Done()
		}(chunk)
	}
	wg.Wait()

	for _, chunk := range s.Chunks {
		chunk.UpdateRect()
	}
}

func (s *Sandbox) TempUpdate() {
	var wg sync.WaitGroup
	for _, chunk := range s.Chunks {
		wg.Add(1)
		go func(s *Sandbox, c *Chunk) {
			NewWorker(s, c).UpdateChunkTemp()
			wg.Done()
		}(s, chunk)
	}
	wg.Wait()
}

func (s *Sandbox) StateUpdate() {
	var wg sync.WaitGroup
	for _, chunk := range s.Chunks {
		wg.Add(1)
		go func(s *Sandbox, c *Chunk) {
			NewWorker(s, c).UpdateChunkState()
			wg.Done()
		}(s, chunk)
	}
	wg.Wait()
}

func (s *Sandbox) Update() {
	s.RemoveEmptyChunks()
	s.MoveUpdate()
	s.TempUpdate()
	s.StateUpdate()
}

func (s *Sandbox) KeepAlive(x, y int) {
	chunk := s.GetChunk(x, y)
	if chunk != nil {
		chunk.KeepAlive(x, y)
	}
}

func (s *Sandbox) Draw(pix []byte, screenWidth int) {
	for _, c := range s.Chunks {
		for i, cell := range c.cells {

			x := i%c.Width + c.X*c.Width
			y := i/c.Width + c.Y*c.Height
			idx := (x + y*screenWidth)
			if isEmpty(cell) {
				pix[idx*4] = 0
				pix[idx*4+1] = 0
				pix[idx*4+2] = 0
				pix[idx*4+3] = 0
				continue
			}

			r := 0
			g := 0
			b := 0
			if cell.temp < 0 {
				b = -cell.temp
				g = -cell.temp / 30
			} else {
				r = cell.temp
			}

			if cell.cType == FIRE {
				g += cell.extraData1
				r -= cell.extraData2 / 3
				g -= cell.extraData2 / 3
				b -= cell.extraData2 / 3
			}
			cR := cell.BaseColor().R
			cG := cell.BaseColor().G
			cB := cell.BaseColor().B
			cA := cell.BaseColor().A
			r = int(cR) + r + cell.colorOffset
			g = int(cG) + g + cell.colorOffset
			b = int(cB) + b + cell.colorOffset

			pix[idx*4] = uint8(misc.Clamp(r, 0, 255))
			pix[idx*4+1] = uint8(misc.Clamp(g, 0, 255))
			pix[idx*4+2] = uint8(misc.Clamp(b, 0, 255))
			pix[idx*4+3] = cA
		}
	}
}
