package main

import (
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/entities"
	"github.com/StCredZero/vectrek/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	g := &game.Game{
		Counter:  0,
		Entities: []*entities.Entity{entities.NewEntity(constants.ScreenWidth/2, constants.ScreenHeight/2)},
	}

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
