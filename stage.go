package main

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/model"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/geom"
	"github.com/eihigh/zu/tick"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type stage struct {
	id    int
	start int
	mode  string
	stage *model.Stage
}

func newStage(id int) *stage {
	return &stage{
		id:    id,
		start: gt.Elapsed(),
		mode:  "open",
		stage: newStageModel(id),
	}
}

func (s *stage) update() (next string) {
	t := gt.SubTick(s.start)

	switch s.mode {
	case "open":
		if t.Elapsed() > 100 {
			s.mode = "main"
			s.start = t.Elapsed()
		}

	case "main":
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			s.mode = "pause" // カウントリセットしない
			break
		}

		i := &model.Input{}
		x, y := ebiten.CursorPosition()
		i.Cursor = geom.Pt(float64(x), float64(y))
		i.Shoot = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		i.Left = ebiten.IsKeyPressed(ebiten.KeyI) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
		i.Up = ebiten.IsKeyPressed(ebiten.KeyP) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
		i.Right = ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
		i.Down = ebiten.IsKeyPressed(ebiten.KeyN) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)

		switch s.stage.UpdateMain(i) {
		case "clear":
			s.mode = "clear"
			s.start = gt.Elapsed()
		case "failed":
			s.mode = "failed"
			s.start = gt.Elapsed()
		}

	case "clear":
		// clickで次へ
		i := &model.Input{}
		s.stage.UpdateMain(i)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.id++
			s.stage = newStageModel(s.id)
			s.mode = "open"
			s.start = gt.Elapsed()
		}

	case "failed":
		// clickでタイトルに戻る
		i := &model.Input{}
		s.stage.UpdateMain(i)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return "title"
		}

	case "pause":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.mode = "main" // カウントリセットしない
			break
		}
	}

	return ""
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

func (s *stage) draw() {
	t := gt.SubTick(s.start)

	switch s.mode {
	case "open":
		s.stage.DrawObjects(t)
		s.stage.DrawUI(t)

		// STAGE N を描画
		msg := "STAGE " + strconv.Itoa(s.id)
		x, y := util.Col(6), util.Row(5)
		util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

		t.From(50, func(tick.Tick) {
			msg := "START!"
			x, y := util.Col(6), util.Row(7)
			util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
		})

	case "main":
		s.stage.DrawObjects(t)
		s.stage.DrawUI(t)

		msg := "PRESS SPACE TO PAUSE"
		x, y := util.Col(6), util.Row(12)-6
		util.DrawTextCenter(asset.Screen, msg, 12, x, y, color.White)

	case "clear":
		s.stage.DrawObjects(t)
		s.stage.DrawUI(t)

		msg := "STAGE CLEAR"
		x, y := util.Col(6), util.Row(5)
		util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

		msg = "CLICK TO GO TO THE NEXT STAGE"
		x, y = util.Col(6), util.Row(7)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)

	case "failed":
		s.stage.DrawObjects(t)
		s.stage.DrawUI(t)

		msg := "GAME OVER"
		x, y := util.Col(6), util.Row(5)
		util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

		msg = "最終到達ステージ " + strconv.Itoa(s.id)
		x, y = util.Col(6), util.Row(6)
		util.DrawTextCenter(asset.Screen, msg, 36, x, y, color.White)

		msg = "CLICK TO BACK TO TITLE"
		x, y = util.Col(6), util.Row(7)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)

	case "pause":
		s.stage.DrawObjects(t)
		s.stage.DrawUI(t)

		msg := "PAUSE"
		x, y := util.Col(6), util.Row(5)
		util.DrawTextCenter(asset.Screen, msg, 48, x, y, color.White)

		msg = "CLICK TO RESUME"
		x, y = util.Col(6), util.Row(7)
		util.DrawTextCenter(asset.Screen, msg, 24, x, y, color.White)
	}
}
