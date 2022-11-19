package sandbox

import (
	"sort"
	"sync"

	. "github.com/mrmarble/sandbox/pkg/misc"
	"pgregory.net/rand"
)

type Change struct {
	dst, src int
	chunk    *Chunk
}

type Chunk struct {
	Width, Height int
	X, Y          int

	// Dirty rect
	MinX, MinY int
	MaxX, MaxY int
	// Working dirty rect
	minXw, minYw int
	maxXw, maxYw int

	filledCells int
	cells       []*Cell
	changes     []Change

	filledCellsMutex sync.Mutex
	changesMutex     sync.Mutex
	dirtyRectMutex   sync.Mutex
}

func NewChunk(width, height, x, y int) *Chunk {
	return &Chunk{
		Width:  width,
		Height: height,
		X:      x,
		Y:      y,
		cells:  make([]*Cell, width*height),
	}
}

func (c *Chunk) KeepAlive(x, y int) {
	c.KeepAliveAt(c.GetIndex(x, y))
}
func (c *Chunk) KeepAliveAt(i int) {
	x := i % c.Width
	y := i / c.Width

	c.dirtyRectMutex.Lock()
	c.minXw = Clamp(Min(x-2, c.minXw), 0, c.Width)
	c.minYw = Clamp(Min(y-2, c.minYw), 0, c.Height)
	c.maxXw = Clamp(Max(x+2, c.maxXw), 0, c.Width)
	c.maxYw = Clamp(Max(y+2, c.maxYw), 0, c.Height)
	c.dirtyRectMutex.Unlock()
}

func (c *Chunk) UpdateRect() {
	c.MinX = c.minXw
	c.MinY = c.minYw
	c.MaxX = c.maxXw
	c.MaxY = c.maxYw

	c.minXw = c.Width
	c.minYw = c.Height
	c.maxXw = -1
	c.maxYw = -1
}

// GetIndex returns the index of the cell at the given coordinates.
func (c *Chunk) GetIndex(x, y int) int {
	return (x - c.X*c.Width) + (y-c.Y*c.Height)*c.Width
}

func (c *Chunk) InBounds(x, y int) bool {
	left := c.X * c.Width
	right := left + c.Width
	top := c.Y * c.Height
	bottom := top + c.Height
	return x >= left && x < right &&
		y >= top && y < bottom
}

func (c *Chunk) IsEmpty(x, y int) bool {
	return c.InBounds(x, y) && c.IsEmptyAt(c.GetIndex(x, y))
}

func (c *Chunk) IsEmptyAt(i int) bool {
	return isEmpty(c.cells[i])
}

func (c *Chunk) GetCell(x, y int) *Cell {
	return c.GetCellAt(c.GetIndex(x, y))
}

func (c *Chunk) GetCellAt(i int) *Cell {
	return c.cells[i]
}

func (c *Chunk) SetCell(x, y int, cell *Cell) {
	c.SetCellAt(c.GetIndex(x, y), cell)
}

func isEmpty(cell *Cell) bool {
	return cell == nil || cell.CType == AIR
}

func (c *Chunk) SetCellAt(i int, cell *Cell) {
	if isEmpty(c.cells[i]) && !isEmpty(cell) {
		c.filledCellsMutex.Lock()
		c.filledCells++
		c.filledCellsMutex.Unlock()
	} else if !isEmpty(c.cells[i]) && isEmpty(cell) {
		c.filledCellsMutex.Lock()
		c.filledCells--
		c.filledCellsMutex.Unlock()
	}
	c.cells[i] = cell
	c.KeepAliveAt(i)
}

func (c *Chunk) MoveCell(src *Chunk, x, y, dx, dy int) {
	c.changesMutex.Lock()
	c.changes = append(c.changes, Change{dst: c.GetIndex(dx, dy), src: src.GetIndex(x, y), chunk: src})
	c.changesMutex.Unlock()
}

func (c *Chunk) ApplyChanges() {
	// remove changes that have the destination cell occupied
	for i := 0; i < len(c.changes); i++ {
		if !c.IsEmptyAt(c.changes[i].dst) {
			c.changes = append(c.changes[:i], c.changes[i+1:]...)
			i--
		}
	}

	// sort changes by destination index
	sort.Slice(c.changes, func(i, j int) bool {
		return c.changes[i].dst < c.changes[j].dst
	})

	// pick random source for each destination
	iPrev := 0
	c.changes = append(c.changes, Change{-1, -1, nil}) // catch the last one
	for i := 0; i < len(c.changes)-1; i++ {
		if c.changes[i+1].dst != c.changes[i].dst {
			rng := rand.Intn(i-iPrev+1) + iPrev

			dst := c.changes[rng].dst
			src := c.changes[rng].src
			chunk := c.changes[rng].chunk

			c.SetCellAt(dst, chunk.GetCellAt(src))
			chunk.SetCellAt(src, nil)

			iPrev = i + 1
		}
	}
	c.changes = c.changes[:0]
}
