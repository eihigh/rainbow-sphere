package model

import (
	"fmt"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/mathg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (s *Stage) DrawObjects(t T) {
	{
		// area
		i := s.areaImage
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(1, 1, 1, 0.4)
		op.GeoM.Translate(s.area.Min.X, s.area.Min.Y)
		asset.Screen.DrawImage(i, op)
	}

	{
		// sphere
		for _, sp := range s.spheres {
			i := asset.SphereImages[sp.index]
			op := &ebiten.DrawImageOptions{}
			util.Center(&op.GeoM, i)
			op.GeoM.Scale(sp.scale, sp.scale)

			switch sp.state.State() {
			case "":

			case "damage":
				// ダメージを受けてから24fの間点滅する
				sp.state.Tick.RepeatFor(0, 6, 4, func(n int, t T) {
					t.Span(0, 3, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.6)
					}).Span(0, 3, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.9)
					})
				})

			case "break":
				// 拡大しつつフェードアウト
				sp.state.Tick.Span(0, 18, func(t T) {
					scale := mathg.Lerp(1, 3, t.Rate())
					alpha := mathg.Lerp(1, 0, t.Rate())
					op.GeoM.Scale(scale, scale)
					op.ColorM.Scale(1, 1, 1, alpha)
				}).From(0, func(T) {
					op.ColorM.Scale(0, 0, 0, 0) // 見た目上消える
				})
			}

			op.GeoM.Scale(0.5, 0.5)
			dir := sp.vel.Angle()
			op.GeoM.Rotate(dir)
			op.GeoM.Translate(sp.pos.X, sp.pos.Y)
			asset.Screen.DrawImage(i, op)
		}
	}

	// debug
	if asset.Debug {
		d := ""
		for _, sp := range s.spheres {
			d += fmt.Sprintf("%d: %d\n", sp.index, sp.hp)
		}
		ebitenutil.DebugPrint(asset.Screen, d)
	}

	{
		// shoot
		i := asset.ShootImage
		for _, sh := range s.shoots {
			op := &ebiten.DrawImageOptions{}
			util.Center(&op.GeoM, i)
			op.GeoM.Scale(0.25, 0.25)
			dir := sh.vel.Angle()
			op.GeoM.Rotate(dir)
			op.GeoM.Translate(sh.pos.X, sh.pos.Y)
			asset.Screen.DrawImage(i, op)
		}
	}

	{
		// player
		i := asset.PlayerImage
		op := &ebiten.DrawImageOptions{}
		util.Center(&op.GeoM, i)
		if s.player.left {
			op.GeoM.Scale(-1, 1)
		}

		switch s.player.state.State() {
		case "":
			// 開幕無敵
			s.player.state.Tick.Span(1, startInvincibleTick, func(t T) {
				t.Repeat(0, 8, func(n int, t T) {
					t.Span(0, 4, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.2)
					}).Span(0, 4, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.7)
					})
				})
			})

		case "damage":
			// ダメージを受けてから点滅する
			s.player.state.Tick.Span(0, 24, func(t T) {
				t.Repeat(0, 6, func(n int, t T) {
					t.Span(0, 3, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.2)
					}).Span(0, 3, func(T) {
						op.ColorM.Scale(1, 1, 1, 0.7)
					})
				})
			})

		case "dead":
			// 拡大しつつフェードアウト
			s.player.state.Tick.Span(0, 18, func(t T) {
				scale := 1.0 + 2.0*t.ElapsedRate()
				op.GeoM.Scale(scale, scale)
				alpha := 1.0 - t.ElapsedRate()
				op.ColorM.Scale(1, 1, 1, alpha)
			}).From(0, func(T) {
				op.ColorM.Scale(0, 0, 0, 0) // 見た目上消える
			})
		}

		op.GeoM.Translate(s.player.pos.X, s.player.pos.Y)
		asset.Screen.DrawImage(i, op)
	}

}

func (s *Stage) DrawUI(t T) {
	// draw hearts
	scale := 1.0 / 8
	w := 16.0 * float64(s.player.hp)
	x := s.player.pos.X - w/2
	y := s.player.pos.Y - 24 - 16
	for i := 0; i < s.player.hp; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		asset.Screen.DrawImage(asset.HeartImage, op)
		x += 16
	}
}
