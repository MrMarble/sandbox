package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mrmarble/sandbox/pkg/game"
)

func main() {
	game := game.New()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
