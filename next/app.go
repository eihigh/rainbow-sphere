package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type empty struct{}

var _ ebiten.Game = (*app)(nil)

// またの名をAppContext
type app struct {
	ticks   int
	screen  *ebiten.Image
	windows map[window]empty

	// 「シーン」を跨ぐものは共通の祖先に置くことに注意 (lifting state up)
	// かつ、子contextはしばしば親contextを埋め込む

	// ここから下のフィールドは
	// コルーチン使えば要らないのになあってやつ
	//
	// titleが終了すると必然的にtitleWindowも消される
	// それが嫌ならappにtitleWindowを持たせる
	// でもそれは本当によろしいのか？
	// タイトル特有のリソースがないせいで却って混乱を招いている
	// 実際、多くの小さいアプリケーションでは多数のwindowを力技で入れ替えて
	// なんとかする構造が向いているのかもしれない...
	curr  string
	title lifecycle[*title]
	stage lifecycle[*stage]
}

type window interface {
	z() float64
	draw(*drawer)
}

// 他にもレイアウトのためのwidget interfaceとかがあると思うけど
// 今回は省略

func newApp() (*app, error) {
	a := &app{
		windows: make(map[window]empty),
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Rainbow Impact with Next-Gen Architecture")

	// 初期シーンはタイトル
	a.curr = "title"
	t := &title{app: a}
	a.title.begin(t)

	return a, nil
}

func (a *app) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	a.ticks++

	// contextのtickを呼び出す
	switch a.curr {
	case "title":
		title, _ := a.title.get()
		title.tick()

	case "stage":
		stage, _ := a.stage.get()
		stage.tick()
	}

	return nil
}

func (a *app) gotoStage() {
	a.title.finish()
	a.curr = "stage"
	s := &stage{}
	a.stage.begin(s)
}

func (a *app) Draw(screen *ebiten.Image) {
	a.screen = screen

	// 強い意志を持ってContextにはDrawメソッドを持たせず
	// 代わりにwindowに持たせる

	rootDrawer := &drawer{image: screen}

	// a.windowsをzでソートしてdrawする
	ws := make([]window, 0, len(a.windows))
	for w := range a.windows {
		ws = append(ws, w)
	}
	sort.Slice(ws, func(i, j int) bool {
		return ws[i].z() < ws[j].z()
	})
	for _, w := range ws {
		w.draw(rootDrawer)
	}
}

func (a *app) Layout(ww, wh int) (sw, sh int) {
	return ww, wh
}
