package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	a, err := newApp()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(a); err != nil {
		panic(err)
	}
}
