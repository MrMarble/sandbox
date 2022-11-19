package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage = ebiten.NewImage(3, 3)
	subImage   = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func Rect(image *ebiten.Image, x, y, w, h int, color color.Color, filled bool) {
	whiteImage.Fill(color)
	var path vector.Path
	var vs []ebiten.Vertex
	var is []uint16
	r := 2
	if filled {
		r = 1
	}
	for i := 0; i < r; i++ {
		path.MoveTo(float32(x+i), float32(y+i))
		path.LineTo(float32(x+w-i), float32(y+i))
		path.LineTo(float32(x+w-i), float32(y+h-i))
		path.LineTo(float32(x+i), float32(y+h-i))
		path.LineTo(float32(x+i), float32(y+i))
	}

	vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	image.DrawTriangles(vs, is, subImage, &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd})
}

func Button(image *ebiten.Image, str string, x, y int, col color.Color, active bool) {
	width := 25
	height := 8
	padding := 3
	inset := (width - len(str)*5) / 2
	Rect(image, x, y, width+padding*2, height+padding*2, col, true)
	text.Draw(image, str, bitmapfont.Gothic10r, x+padding+inset+4, y+height+padding, color.White)
	if active {
		Rect(image, x, y, width+padding*2, height+padding*2, color.White, false)
	}
}
