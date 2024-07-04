package main

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

// またの名をTitleContext
type title struct {
	*app
	window lifecycle[*titleWindow]

	// titleWindowが依存するようなタイトル特有のリソースがなくて
	// 却って設計上の混乱を招いている
	// 本当はそれらを解放するタイミングがtitleのfinishになるんだけど

	// リソースにはアセットも含む
}

func (t *title) begin() {
	w := &titleWindow{}
	t.window.begin(w)
}

func (t *title) finish() {
	t.window.finish()
}

func (t *title) tick() {
	// contextはあくまでライフタイムの管理であり、
	// フェードとかはやらないことを思い出す
	// ただし、兄弟contextが共存することはありうる。今回だとstage
}

// not a context
type titleWindow struct {
	*app
}

func (t *titleWindow) begin() {
	t.windows[t] = struct{}{}
}

func (t *titleWindow) finish() {
	delete(t.windows, t)
}

func (t *titleWindow) draw(d *drawer) {
	// titleTextやhowToPlayはwindow型でもいいし、でなくてもいい
	// 今回はwindowにはしない
}

func (t *titleWindow) z() float64 {
	return 0
}
