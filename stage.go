package main

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/model"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/hsm"
	"github.com/eihigh/zu/mathg"
	"github.com/eihigh/zu/tick"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var stage struct {
	id    int
	stage *model.Stage
}

func enterStage(h *hsm.HSM) {
	stage.id = firstStage
}

func enterStageOpen(h *hsm.HSM) {
	stage.stage = newStageModel(stage.id)
}

func updateStageOpen(h *hsm.HSM) {
	if h.Tick().Elapsed() > 100 {
		h.Change("/stage/main")
	}
}

func updateStageMain(h *hsm.HSM) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		h.Change("/stage/main/pause")
	}

	i := &model.Input{}
	x, y := ebiten.CursorPosition()
	i.Cursor = mathg.Pt(float64(x), float64(y))
	i.Shoot = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	i.Left = ebiten.IsKeyPressed(ebiten.KeyI) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	i.Up = ebiten.IsKeyPressed(ebiten.KeyP) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
	i.Right = ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	i.Down = ebiten.IsKeyPressed(ebiten.KeyN) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)

	switch stage.stage.UpdateMain(i) {
	case "clear":
		h.Change("/stage/clear")
	case "failed":
		h.Change("/stage/failed")
	}
}

func updateStageMainPause(h *hsm.HSM) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.Change("/stage/main")
	}
}

func updateStageClear(h *hsm.HSM) {
	i := model.Input{}
	stage.stage.UpdateMain(&i)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		stage.id++
		h.Change("/stage/open")
	}
}

func updateStageFailed(h *hsm.HSM) {
	// clickでタイトルに戻る
	i := &model.Input{}
	stage.stage.UpdateMain(i)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		h.Change("/title")
	}
}

func drawStageOpen(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	// STAGE N を描画
	msg := "STAGE " + strconv.Itoa(stage.id)
	x, y := util.Col(6), util.Row(5)
	util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

	t.From(50, func(tick.Tick) {
		msg := "START!"
		x, y := util.Col(6), util.Row(7)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
	})
}

func drawStageMain(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	msg := "PRESS SPACE TO PAUSE"
	x, y := util.Col(6), util.Row(12)-6
	util.DrawTextCenter(asset.Screen, msg, 12, x, y, color.White)
}

func drawStageMainPause(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	msg := "PRESS SPACE TO PAUSE"
	x, y := util.Col(6), util.Row(12)-6
	util.DrawTextCenter(asset.Screen, msg, 12, x, y, color.White)
}

func drawStageClear(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	msg := "STAGE CLEAR"
	x, y := util.Col(6), util.Row(5)
	util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

	msg = "CLICK TO GO TO THE NEXT STAGE"
	x, y = util.Col(6), util.Row(7)
	util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
}

func drawStageFailed(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	msg := "GAME OVER"
	x, y := util.Col(6), util.Row(5)
	util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

	msg = "最終到達ステージ " + strconv.Itoa(stage.id)
	x, y = util.Col(6), util.Row(6)
	util.DrawTextCenter(asset.Screen, msg, 36, x, y, color.White)

	msg = "CLICK TO BACK TO TITLE"
	x, y = util.Col(6), util.Row(7)
	util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
}

func drawStagePause(h *hsm.HSM) {
	t := h.Tick()
	stage.stage.DrawObjects(t)
	stage.stage.DrawUI(t)

	msg := "PAUSE"
	x, y := util.Col(6), util.Row(5)
	util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

	msg = "CLICK TO RESUME"
	x, y = util.Col(6), util.Row(7)
	util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
}

func newStageModel(id int) *model.Stage {
	id-- // 1 => 0
	var c *model.Config
	if id >= len(configs) {

		// 一定ステージクリア後
		// endlessConfigからランダムに選択
		// HPはステージごとに1減る
		hp := 4 - (id-len(configs))/4
		if hp <= 0 {
			hp = 1
		}
		i := rand.Intn(len(endlessConfigs))
		c = endlessConfigs[i]
		c.HP = hp

	} else {
		c = configs[id]
	}

	return model.NewStage(c)
}
