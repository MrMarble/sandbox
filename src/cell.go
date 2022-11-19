package main

import (
	"image/color"

	"pgregory.net/rand"
)

//go:generate stringer -type=CellType

type CellType int

const (
	AIR CellType = iota
	SAND
	WATR
	WALL
	STNE
	SMKE
	STEM
	WOOD
	FIRE
)

func getColor(cType CellType) string {
	switch cType {
	case SAND:
		return "#c2b280"
	case WATR:
		return "#07a9be"
	case WALL:
		return "#252525"
	case STNE:
		return "#808080"
	case SMKE:
		return "#101010"
	case STEM:
		return "#ADD8E6"
	case WOOD:
		return "#BA8C63"
	case FIRE:
		return "#F44D2B"
	default:
		return "#000000"
	}
}

type Cell struct {
	cType CellType

	color       color.RGBA
	colorOffset int

	temp       int
	extraData1 int
	extraData2 int
}

func ParseHexColor(s string) (c color.RGBA) {
	c.A = 0xff

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
	}
	return
}

func NewCell(cType CellType) *Cell {
	cell := &Cell{
		cType:       cType,
		color:       ParseHexColor(getColor(cType)),
		colorOffset: rand.Intn(20) + -10,
	}
	switch cType {
	case SMKE:
		cell.extraData1 = 90 + (rand.Intn(40) + -20)
		cell.extraData2 = 90
	case STEM:
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
	case WATR:
		return 5
	case STNE, WOOD:
		return 1
	case FIRE:
		return 2
	case STEM, SMKE:
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
