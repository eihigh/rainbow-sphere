package main

import (
	_ "embed"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/zu/hsm"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed crt.go
var crtText []byte

var (
	firstStage = 1

	scenes = []*hsm.State{
		{
			Name:   "/title",
			Enter:  enterTitle,
			Update: updateTitle,
			Draw:   drawTitle,
		},
		{
			Name:   "/title/howtoplay",
			Update: updateTitleHowtoplay,
			Draw:   drawTitleHowtoplay,
		},
		{
			Name:  "/stage",
			Enter: enterStage,
		},
		{
			Name:   "/stage/open",
			Enter:  enterStageOpen,
			Update: updateStageOpen,
			Draw:   drawStageOpen,
		},
		{
			Name:   "/stage/main",
			Update: updateStageMain,
			Draw:   drawStageMain,
		},
		{
			Name:   "/stage/main/pause",
			Update: updateStageMainPause,
			Draw:   drawStageMainPause,
		},
		{
			Name:   "/stage/clear",
			Update: updateStageClear,
			Draw:   drawStageClear,
		},
		{
			Name:   "/stage/failed",
			Update: updateStageFailed,
			Draw:   drawStageFailed,
		},
	}

	scene = hsm.NewHSM(scenes, "/title", 1)

	crtShader *ebiten.Shader
)

type app struct{}

func newApp() (*app, error) {
	rand.Seed(time.Now().UnixNano())
	a := &app{}
	return a, nil
}

func (a *app) Update() error {
	if asset.Debug && ebiten.IsKeyPressed(ebiten.KeyQ) {
		return io.EOF
	}
	scene.Update()
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	asset.Screen = screen
	scene.Draw()
}

func (a *app) Layout(ow, oh int) (int, int) {
	return asset.VW, asset.VH
}

func (a *app) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	if crtShader == nil {
		var err error
		crtShader, err = ebiten.NewShader([]byte(crtText))
		if err != nil {
			panic(err)
		}
	}
	b := offscreen.Bounds()
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = offscreen
	opts.GeoM = geoM
	screen.DrawRectShader(b.Dx(), b.Dy(), crtShader, opts)
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "debug" {
		asset.Debug = true
	}
	if err := asset.Load(); err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(asset.VW, asset.VH)
	ebiten.SetWindowTitle("Rainbow Impact")
	app, err := newApp()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}
