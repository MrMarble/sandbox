package game

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mrmarble/sandbox/pkg/sandbox"
	"github.com/mrmarble/sandbox/pkg/ui"
)

var (
	screenWidth  = 640
	screenHeight = 480
	margin       = 50

	offscrenOptions = &ebiten.DrawImageOptions{}
)

type Game struct {
	pixels  []byte
	sandbox *sandbox.Sandbox
	menu    *ui.Menu

	pause       bool
	debug       bool
	tempOverlay bool

	brushSize int

	prevPos   [2]int
	cursorPos [2]int

	cellQueue [][2][2]int

	offscreen *ebiten.Image
}

func New() *Game {
	ebiten.SetWindowTitle("Sandbox")
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	return &Game{
		sandbox:     sandbox.NewSandbox(screenWidth-margin, screenHeight-margin),
		brushSize:   10,
		tempOverlay: true,
		offscreen:   ebiten.NewImage(screenWidth-margin, screenHeight-margin),
		menu:        ui.NewMenu(margin/2, screenHeight-margin/2+5),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	return screenWidth * int(s), screenHeight * int(s)
}

func (g *Game) Update() error {
	g.updateCursor()
	g.handleInput()
	g.menu.Update()
	g.placeQueueParticles()

	if !g.pause || inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		g.sandbox.Update(g.tempOverlay)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		log.Println("init")
		g.pixels = make([]byte, g.offscreen.Bounds().Dx()*g.offscreen.Bounds().Dy()*4)
		offscrenOptions.GeoM.Translate(float64(margin/2), float64(margin/2))
	}

	g.sandbox.Draw(g.pixels, g.offscreen.Bounds().Dx(), g.tempOverlay)
	g.offscreen.WritePixels(g.pixels)

	// Brush size
	offscreenX, offscreenY := offscreenCursor(g.cursorPos[0], g.cursorPos[1])
	ui.Rect(g.offscreen, offscreenX-g.brushSize/2, offscreenY-g.brushSize/2, g.brushSize, g.brushSize, color.White, false)
	// Border
	ui.Rect(g.offscreen, 0, 0, screenWidth-margin, screenHeight-margin, color.RGBA{20, 20, 20, 100}, false)

	g.menu.Draw(screen)
	g.debugInfo(screen)
	screen.DrawImage(g.offscreen, offscrenOptions)
}

func offscreenCursor(x, y int) (int, int) {
	return x - margin/2, y - margin/2
}
