package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/zu/tick"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// gt - グローバルな時の流れを表すtick。これをいろんなところで使う
	gt tick.Tick

	firstStage = 1
)

type app struct {
	child interface{}
}

func newApp() (*app, error) {
	rand.Seed(time.Now().UnixNano())
	a := &app{
		child: newTitle(),
		// child: newStage(1),
	}
	return a, nil
}

func (a *app) Update() error {
	defer gt.Advance(1)

	if asset.Debug && ebiten.IsKeyPressed(ebiten.KeyQ) {
		return io.EOF
	}

	// 文字列のかわりに型で分岐
	// updateやdrawのシグネチャがどうなるか読めなかったためinterface定義はしてない
	// あとシーンは少ないのでdynamic castしてもそう問題にならない
	switch c := a.child.(type) {
	case *title:
		switch c.update() {
		case "stage":
			a.child = newStage(firstStage)
		}

	case *stage:
		switch c.update() {
		case "title":
			a.child = newTitle()
		}
	}

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	asset.Screen = screen

	switch c := a.child.(type) {
	case *title:
		c.draw()

	case *stage:
		c.draw()
	}
}

func (a *app) Layout(ow, oh int) (int, int) {
	return asset.VW, asset.VH
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
