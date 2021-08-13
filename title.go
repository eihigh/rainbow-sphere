package main

import (
	"image/color"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/tick"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var howToPlay = `【操作説明】
結界内から暴走陰陽玉を破壊せよ！
※ライフはステージクリアでリセットされますが、
初期ライフはステージが進むごとにどんどん減ります

Pause: Space
Move: WASD or ↑←↓→
Aim: Mouse Cursor
Shoot: Mouse Click

CLICK TO START
`

type title struct {
	start int // 開始tick
	mode  string
}

func newTitle() *title {
	return &title{
		start: gt.Elapsed(),
		mode:  "",
	}
}

func (tt *title) update() (next string) {

	t := gt.SubTick(tt.start)
	switch tt.mode {
	case "":
		if t.Elapsed() > 60 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			tt.mode = "howToPlay"
			tt.start = gt.Elapsed()
		}

	case "howToPlay":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return "stage"
		}
	}

	return ""
}

func (tt *title) draw() {
	t := gt.SubTick(tt.start)
	switch tt.mode {
	case "":
		// 落ちながらフェードイン
		alpha := 1.0
		x, y := util.Col(6), util.Row(4)
		t.Span(0, 60, func(t tick.Tick) {
			e := ease.OutQuad(t.ElapsedRate())
			y = util.Row(2) + util.Row(2)*e
			alpha = e
		})

		msg := "恐怖！虹色の陰陽玉\n～ Rainbow Impact"
		clr := color.Alpha{uint8(alpha * 255)}
		util.DrawTextCenter(asset.Screen, msg, 48, x, y, clr)

		// 次に進むテキスト
		t.From(60, func(tick.Tick) {
			msg = "CLICK TO START"
			x, y = util.Col(6), util.Row(7)
			util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
		})

	case "howToPlay":
		msg := howToPlay
		x, y := util.Col(6), util.Row(6)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
	}
}
