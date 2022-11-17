package main

import (
	"math/rand"
	"sort"
	"sync"
)

type Change struct {
	dst, src int
	chunk    *Chunk
}

type Chunk struct {
	width, height int
	x, y          int

	// Dirty rect
	minX, minY int
	maxX, maxY int
	// Working dirty rect
	minXw, minYw int
	maxXw, maxYw int

	filledCells int
	cells       []Cell
	changes     []Change

	filledCellsMutex sync.Mutex
	changesMutex     sync.Mutex
	dirtyRectMutex   sync.Mutex
}

func NewChunk(width, height, x, y int) *Chunk {
	return &Chunk{
		width:  width,
		height: height,
		x:      x,
		y:      y,
		cells:  make([]Cell, width*height),
	}
}

func (c *Chunk) KeepAlive(x, y int) {
	c.KeepAliveAt(c.GetIndex(x, y))
}
func (c *Chunk) KeepAliveAt(i int) {
	x := i % c.width
	y := i / c.width

	c.dirtyRectMutex.Lock()
	c.minXw = clamp(min(x-2, c.minXw), 0, c.width)
	c.minYw = clamp(min(y-2, c.minYw), 0, c.height)
	c.maxXw = clamp(max(x+2, c.maxXw), 0, c.width)
	c.maxYw = clamp(max(y+2, c.maxYw), 0, c.height)
	c.dirtyRectMutex.Unlock()
}

func (c *Chunk) UpdateRect() {
	c.minX = c.minXw
	c.minY = c.minYw
	c.maxX = c.maxXw
	c.maxY = c.maxYw

	c.minXw = c.width
	c.minYw = c.height
	c.maxXw = -1
	c.maxYw = -1
}

// GetIndex returns the index of the cell at the given coordinates.
func (c *Chunk) GetIndex(x, y int) int {
	return (x - c.x*c.width) + (y-c.y*c.height)*c.width
}

func (c *Chunk) InBounds(x, y int) bool {
	left := c.x * c.width
	right := left + c.width
	top := c.y * c.height
	bottom := top + c.height

	return x >= left && x < right &&
		y >= top && y < bottom
}

func (c *Chunk) IsEmpty(x, y int) bool {
	return c.InBounds(x, y) && c.IsEmptyAt(c.GetIndex(x, y))
}

func (c *Chunk) IsEmptyAt(i int) bool {
	return c.cells[i].cType == AIR
}

func (c *Chunk) GetCell(x, y int) *Cell {
	return c.GetCellAt(c.GetIndex(x, y))
}

func (c *Chunk) GetCellAt(i int) *Cell {
	return &c.cells[i]
}

func (c *Chunk) SetCell(x, y int, cell Cell) {
	c.SetCellAt(c.GetIndex(x, y), cell)
}

func (c *Chunk) SetCellAt(i int, cell Cell) {
	if c.cells[i].cType == AIR && cell.cType != AIR {
		c.filledCellsMutex.Lock()
		c.filledCells++
		c.filledCellsMutex.Unlock()
	} else if c.cells[i].cType != AIR && cell.cType == AIR {
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

			c.SetCellAt(dst, *chunk.GetCellAt(src))
			chunk.SetCellAt(src, Cell{})

			iPrev = i + 1
		}
	}
	c.changes = c.changes[:0]
}
