package sandbox

import (
	"image/color"

	"pgregory.net/rand"
)

//go:generate stringer -type=CellType

type CellType int

const (
	SAND CellType = iota
	WATER
	WALL
	STONE
	SMOKE
	STEAM
	WOOD
	FIRE
	AIR // special type for empty cells. Always last for easy iteration
)

func (cType CellType) Color() color.RGBA {
	switch cType {
	case SAND:
		return color.RGBA{0xc2, 0xb2, 0x80, 0xff} //#c2b280
	case WATER:
		return color.RGBA{0x07, 0xa9, 0xbe, 0xff} //#07a9be
	case WALL:
		return color.RGBA{0x25, 0x25, 0x25, 0xff} //#252525
	case STONE:
		return color.RGBA{0x80, 0x80, 0x80, 0xff} //#808080
	case SMOKE:
		return color.RGBA{0x10, 0x10, 0x10, 0xff} //#101010
	case STEAM:
		return color.RGBA{0xad, 0xd8, 0xe6, 0xff} //#add8e6
	case WOOD:
		return color.RGBA{0xba, 0x8c, 0x63, 0xff} //#ba8c63
	case FIRE:
		return color.RGBA{0xf4, 0x4d, 0x2b, 0xff} //#f44d2b
	default:
		return color.RGBA{0x00, 0x00, 0x00, 0xff} //"#000000"
	}
}

type Cell struct {
	cType CellType

	colorOffset int

	temp       int
	extraData1 int
	extraData2 int
}

func NewCell(cType CellType) *Cell {
	cell := &Cell{
		cType:       cType,
		colorOffset: rand.Intn(20) + -10,
	}
	switch cType {
	case SMOKE:
		cell.extraData1 = 90 + (rand.Intn(40) + -20)
		cell.extraData2 = 90
	case STEAM:
		cell.temp = 100
	case FIRE:
		cell.extraData1 = rand.Intn(60)
		cell.temp = 130
	}
	return cell
}

func (c *Cell) ThermalConductivity() int {
	switch c.cType {
	case SAND:
		return 3
	case WATER:
		return 5
	case STONE, WOOD:
		return 1
	case FIRE:
		return 2
	case STEAM, SMOKE:
		return 6
	default:
		return 0
	}
}

func (c *Cell) IsFlamable() bool {
	switch c.cType {
	case WOOD:
		return true
	default:
		return false
	}
}

func (c *Cell) BaseColor() color.RGBA {
	switch c.cType {
	case SAND:
		if c.extraData1 > 0 {
			return color.RGBA{0xb1, 0x9d, 0x5e, 0xff} //#b19d5e
		}
		return c.cType.Color()
	default:
		return c.cType.Color()
	}
}
