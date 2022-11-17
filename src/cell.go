package main

import "image/color"

//go:generate stringer -type=CellType

type CellType int

const (
	EMPTY CellType = iota
	SAND
	WATER
	WALL
	STONE
)

type Cell struct {
	cType CellType

	color color.RGBA
}

func NewCell(cType CellType) *Cell {
	clr := color.RGBA{0, 0, 0, 0}
	switch cType {
	case SAND:
		clr = color.RGBA{0xc2, 0xb2, 0x80, 255}
	case WATER:
		clr = color.RGBA{0x09, 0xc3, 0xdb, 255}
	case WALL:
		clr = color.RGBA{0x10, 0x10, 0x10, 255}
	case STONE:
		clr = color.RGBA{0x80, 0x80, 0x80, 255}
	}
	return &Cell{cType: cType, color: clr}
}

func (c *Cell) WithVariation(variation byte) *Cell {
	c.color.R += variation
	c.color.G += variation
	c.color.B += variation
	return c
}
