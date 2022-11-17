package main

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
	pingX := 0
	pingY := 0
	//fmt.Println("Moving cell from", x, y, "to", dx, dy, "in chunk", w.chunk.x, w.chunk.y)
	if x == w.chunk.x*w.chunk.width {
		pingX = -1
	}
	if x == w.chunk.x*w.chunk.width+w.chunk.width-1 {
		pingX = 1
	}
	if y == w.chunk.y*w.chunk.height {
		pingY = -1
	}
	if y == w.chunk.y*w.chunk.height+w.chunk.height-1 {
		pingY = 1
	}

	if pingX != 0 {
		w.sandbox.KeepAlive(x+pingX, y)
	}
	if pingY != 0 {
		w.sandbox.KeepAlive(x, y+pingY)
	}
	if pingX != 0 && pingY != 0 {
		w.sandbox.KeepAlive(x+pingX, y+pingY)
	}

	if w.chunk.InBounds(x, y) && w.chunk.InBounds(dx, dy) {
		w.chunk.MoveCell(w.chunk, x, y, dx, dy)
	} else {
		w.sandbox.MoveCell(x, y, dx, dy)
	}
}

func (w *Worker) UpdateChunk() {
	for x := w.chunk.minX; x < w.chunk.maxX; x++ {
		for y := w.chunk.minY; y < w.chunk.maxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.width)
			px := x + w.chunk.x*w.chunk.width
			py := y + w.chunk.y*w.chunk.height
			if c.cType == AIR {
				continue
			}
			switch c.cType {
			case SAND:
				w.MovePowder(px, py, c)
			case WATR:
				w.MoveLiquid(px, py, c)
			case STNE:
				w.MoveSolid(px, py, c)
			case SMKE, STEM:
				w.MoveGas(px, py, c)
			}
		}
	}
}

func (w *Worker) UpdateChunkState() {
	for x := w.chunk.minX; x < w.chunk.maxX; x++ {
		for y := w.chunk.minY; y < w.chunk.maxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.width)
			px := x + w.chunk.x*w.chunk.width
			py := y + w.chunk.y*w.chunk.height
			if c.cType == AIR {
				continue
			}
			switch c.cType {
			case SMKE:
				w.UpdateSmoke(px, py)
			case STEM:
				w.UpdateSteam(px, py)
			case WATR:
				w.UpdateWater(px, py)
			case SAND:
				w.UpdateSand(px, py)
			}
		}
	}
}

func (w *Worker) UpdateChunkTemp() {
	for x := w.chunk.minX; x < w.chunk.maxX; x++ {
		for y := w.chunk.minY; y < w.chunk.maxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.width)
			px := x + w.chunk.x*w.chunk.width
			py := y + w.chunk.y*w.chunk.height
			if c.cType == AIR || c.temp == 0 {
				continue
			}
			temp := c.temp
			conductivity := c.ThermalConductivity()

			if w.InBounds(px, py+1) {
				other := w.GetCell(px, py+1)
				if other != nil && other.temp != 0 {
					tc := conductivity + other.ThermalConductivity()
					t := temp / tc
					c.temp -= t
					other.temp += t
				}
			}
			if w.InBounds(px, py-1) {
				other := w.GetCell(px, py-1)
				if other != nil && other.temp != 0 {
					tc := conductivity + other.ThermalConductivity()
					t := temp / tc
					c.temp -= t
					other.temp += t
				}
			}
			if w.InBounds(px+1, py) {
				other := w.GetCell(px+1, py)
				if other != nil && other.temp != 0 {
					tc := conductivity + other.ThermalConductivity()
					t := temp / tc
					c.temp -= t
					other.temp += t
				}
			}
			if w.InBounds(px-1, py) {
				other := w.GetCell(px-1, py)
				if other != nil && other.temp != 0 {
					tc := conductivity + other.ThermalConductivity()
					t := temp / tc
					c.temp -= t
					other.temp += t
				}
			}
		}
	}
}
