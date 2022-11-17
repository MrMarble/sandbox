package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage = ebiten.NewImage(3, 3)
	subImage   = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func rect(image *ebiten.Image, x, y, w, h int, color color.Color) {
	whiteImage.Fill(color)
	var path vector.Path
	var vs []ebiten.Vertex
	var is []uint16

	for i := 0; i < 2; i++ {
		path.MoveTo(float32(x+i), float32(y+i))
		path.LineTo(float32(x+w-i), float32(y+i))
		path.LineTo(float32(x+w-i), float32(y+h-i))
		path.LineTo(float32(x+i), float32(y+h-i))
		path.LineTo(float32(x+i), float32(y+i))
	}

	vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	image.DrawTriangles(vs, is, subImage, &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd})
}
