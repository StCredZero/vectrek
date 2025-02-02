package main

import (
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/entities"
	"github.com/StCredZero/vectrek/game"
	"github.com/StCredZero/vectrek/sparse"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ents := make([]*entities.Entity, 0)
	ents = append(ents, entities.NewEntity(constants.ScreenWidth/2, constants.ScreenHeight/2))
	ents[0].Key = sparse.Key(0)

	g := &game.Game{
		Counter:  0,
		Entities: ents,
	}

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
