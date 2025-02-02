package main

import (
	"github.com/StCredZero/vectrek/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	g := &game.Game{
		Counter: 0,
		Ship:    game.NewShip(),
	}

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
