package main

import "github.com/hajimehoshi/ebiten/v2"

// 描画をCPUサイドで合成できるようにするための構造体

type drawer struct {
	image *ebiten.Image
	geoM  ebiten.GeoM
}

func (d *drawer) draw(src *ebiten.Image, opt *ebiten.DrawImageOptions) {
	opt.GeoM.Concat(d.geoM)
	d.image.DrawImage(src, opt)
}
