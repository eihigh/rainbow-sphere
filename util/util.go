package util

import (
	"image/color"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/zu/typist"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func Col(i int) float64 { return asset.VW * float64(i) / 12 }
func Row(i int) float64 { return asset.VH * float64(i) / 12 }

func DrawText(dst *ebiten.Image, txt string, f font.Face, x, y int, clr color.Color) {
	y += f.Metrics().Ascent.Round() // dot position
	text.Draw(dst, txt, f, x, y, clr)
}

func DrawTextCenter(dst *ebiten.Image, txt string, size float64, x, y float64, clr color.Color) {
	f := asset.Mplus(size)
	_, lines := typist.Measure(f, txt, int(Col(10)))
	height := f.Metrics().Height.Round() * len(lines)

	y -= float64(height / 2)
	for _, line := range lines {
		xx := x - float64(line.Width/2)
		DrawText(dst, line.Text, f, int(xx), int(y), clr)
		y += float64(f.Metrics().Height.Round())
	}
}

func Center(g *ebiten.GeoM, i *ebiten.Image) {
	w, h := i.Size()
	g.Translate(-float64(w)/2, -float64(h)/2)
}
