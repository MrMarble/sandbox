package sandbox

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

func (w *Worker) SetCell(x, y int, cell *Cell) {
	if w.chunk.InBounds(x, y) {
		w.chunk.SetCell(x, y, cell)
	} else {
		w.sandbox.SetCell(x, y, cell)
	}
}

func (w *Worker) MoveCell(x, y, dx, dy int) {
	pingX := 0
	pingY := 0

	if x == w.chunk.X*w.chunk.Width {
		pingX = -1
	}
	if x == w.chunk.X*w.chunk.Width+w.chunk.Width-1 {
		pingX = 1
	}
	if y == w.chunk.Y*w.chunk.Height {
		pingY = -1
	}
	if y == w.chunk.Y*w.chunk.Height+w.chunk.Height-1 {
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
	for x := w.chunk.MinX; x < w.chunk.MaxX; x++ {
		for y := w.chunk.MinY; y < w.chunk.MaxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.Width)
			if isEmpty(c) {
				continue
			}
			px := x + w.chunk.X*w.chunk.Width
			py := y + w.chunk.Y*w.chunk.Height

			switch c.CType {
			case SAND:
				if c.extraData1 == 0 {
					w.MovePowder(px, py, c)
				} else {
					w.MoveSolid(px, py, c)
				}
			case WATER:
				if c.temp > -80 {
					w.MoveLiquid(px, py, c)
				} else {
					w.MoveSolid(px, py, c)
				}
			case GLASS:
				if c.temp >= 30 {
					w.MoveLiquid(px, py, c)
				} else {
					w.MoveSolid(px, py, c)
				}
			case STONE:
				w.MoveSolid(px, py, c)
			case SMOKE, STEAM:
				w.MoveGas(px, py, c)
			case FIRE:
				w.MoveFire(px, py, c)
			case PLANT:
				if c.extraData2 == 0 {
					w.MovePowder(px, py, c)
				}
			}
		}
	}
}

func (w *Worker) UpdateChunkState() {
	for x := w.chunk.MinX; x < w.chunk.MaxX; x++ {
		for y := w.chunk.MinY; y < w.chunk.MaxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.Width)
			if isEmpty(c) {
				continue
			}
			px := x + w.chunk.X*w.chunk.Width
			py := y + w.chunk.Y*w.chunk.Height
			if c.CType == AIR {
				continue
			}

			switch c.CType {
			case SMOKE:
				w.UpdateSmoke(px, py)
			case STEAM:
				w.UpdateSteam(px, py)
			case WATER:
				w.UpdateWater(px, py)
			case SAND:
				w.UpdateSand(px, py)
			case FIRE:
				w.UpdateFire(px, py)
			case CLNE:
				w.UpdateReplicator(px, py)
			case PLANT:
				w.UpdatePlant(px, py)
			}
		}
	}
}

func (w *Worker) UpdateChunkTemp() {

	directions := [][]int{
		{0, -1}, // bottom
		{0, 1},  // up
		{1, 0},  // right
		{-1, 0}, // left
	}

	for x := w.chunk.MinX; x < w.chunk.MaxX; x++ {
		for y := w.chunk.MinY; y < w.chunk.MaxY; y++ {
			c := w.chunk.GetCellAt(x + y*w.chunk.Width)
			if isEmpty(c) {
				continue
			}
			px := x + w.chunk.X*w.chunk.Width
			py := y + w.chunk.Y*w.chunk.Height
			if c.CType == AIR || c.temp == 0 {
				continue
			}
			temp := c.temp
			conductivity := c.ThermalConductivity()

			for _, dir := range directions {
				w.updateTemp(dir[0]+px, dir[1]+py, temp, conductivity, c.temp)
			}

		}
	}
}

func (w *Worker) updateTemp(px, py, temp, conductivity, cTemp int) {
	if w.InBounds(px, py) {
		other := w.GetCell(px, py)
		if !isEmpty(other) {
			tc := conductivity + other.ThermalConductivity()
			if tc > 0 {
				t := temp / tc
				cTemp -= t
				other.temp += t
				w.chunk.KeepAlive(px, py)
			}
		}
	}
}
