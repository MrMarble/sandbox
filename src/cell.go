package main

import "image/color"

//go:generate stringer -type=CellType

type CellType int

const (
	AIR CellType = iota
	SAND
	WATR
	WALL
	STNE
	SMKE
)

func getColor(cType CellType) string {
	switch cType {
	case SAND:
		return "#c2b280"
	case WATR:
		return "#09c3db"
	case WALL:
		return "#252525"
	case STNE:
		return "#808080"
	case SMKE:
		return "#101010"
	default:
		return "#000000"
	}
}

type Cell struct {
	cType CellType

	color color.RGBA
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
	return &Cell{cType: cType, color: ParseHexColor(getColor(cType))}
}

func (c *Cell) WithVariation(variation byte) *Cell {
	if c.cType == WALL {
		return c
	}
	c.color.R += variation
	c.color.G += variation
	c.color.B += variation
	return c
}
