package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/zu/hsm"
	"github.com/hajimehoshi/ebiten/v2"
)

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

const crtText = `
//go:build ignore

//kage:unit pixels

// Based on https://github.com/XorDev/1PassBlur
// https://www.shadertoy.com/view/DtscRf

package main

const radius = 16.0
const samples = 32.0
const base = 0.5
const glow = 1.5

func hash2(p vec2) vec2 {
  return normalize(
    fract(cos(p * mat2(195,174,286,183)) * 742) - 0.5)
}

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
  blur := vec4(0)
  weights := 0.0
  scale := radius / sqrt(samples)
  offset := hash2(srcPos) * scale
  rot := mat2(-0.7373688, -0.6754904, 0.6754904, -0.7373688)
  for i := 0.0; i < samples; i += 1 {
    // rotate by golden angle
    offset *= rot
    dist := sqrt(i)
    pos := srcPos + offset*dist
    color := imageSrc0At(pos)

    weight := 1.0 / (1+dist)
    blur += color * weight
    weights += weight
  }
  blur /= weights
  clr := mix(blur*glow, imageSrc0At(srcPos), base)

  rgb := clr.rgb
  rgb = clamp(mix(rgb, rgb*rgb, 0.4), 0, 1)

  // vignette
  uv := (srcPos - imageSrc0Origin()) / imageSrc0Size()
  vig := 40*uv.x*uv.y*(1-uv.x)*(1-uv.y)
  rgb *= vec3(pow(vig, 0.3))
  rgb *= vec3(0.95, 1.05, 0.95)

  n := floor(imageDstSize().y / 480) + 1
  rgb *= 1.0 - mod(dstPos.y, n)*0.3

  rgb *= 1.4
  return vec4(rgb, clr.a)
}
`
