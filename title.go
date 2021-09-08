package main

import (
	"image/color"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/colorf"
	"github.com/eihigh/zu/hsm"
	"github.com/eihigh/zu/mathg"
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

var title struct {
	// bgm *audio.Player
}

func enterTitle(h *hsm.HSM) {
	// title.bgm.Play()
}

func updateTitle(h *hsm.HSM) {
	if h.Tick().Elapsed() > 60 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.Change("/title/howtoplay")
	}
}

func updateTitleHowtoplay(h *hsm.HSM) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.Change("/stage/open")
	}
}

func drawTitle(h *hsm.HSM) {
	// 落ちながらフェードイン
	t := h.Tick()
	r := mathg.Clamp(t.Elapsedf() / 60)

	e := ease.OutQuad(r)
	alpha := e
	x := util.Col(6)
	y := mathg.Lerp(util.Row(2), util.Row(4), e)

	/*
		v := tween.Key(0, tween.Values{"y": Row(2), "a": 0.0}).
			Key(60, ease.OutQuad, tween.Values{"y": Row(4), "a": 1.1}).
			Get(scene.Elapsed())
	*/

	msg := "恐怖！虹色の陰陽玉\n～ Rainbow Impact"
	clr := colorf.Alpha{A: alpha}
	util.DrawTextCenter(asset.Screen, msg, 48, x, y, clr)

	// 次に進むテキスト
	t.From(60, func(tick.Tick) {
		msg = "CLICK TO START"
		x, y = util.Col(6), util.Row(7)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
	})
}

func drawTitleHowtoplay(h *hsm.HSM) {
	msg := howToPlay
	x, y := util.Col(6), util.Row(6)
	util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
}
