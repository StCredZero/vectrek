package main

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecs"
	"github.com/StCredZero/vectrek/game"
	"github.com/StCredZero/vectrek/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ents := make([]*ecs.Entity, 0)
	ents = append(ents, ecs.NewEntity(constants.ScreenWidth/2, constants.ScreenHeight/2))
	ents[0].ID = ecs.EntityID(0)

	if false {
		instance := &game.Game{
			Counter:  0,
			Entities: ents,
		}
		ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
		ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
		if err := ebiten.RunGame(instance); err != nil {
			log.Fatal(err)
		}
	}
	instance := ecs.NewInstance(ecs.Parameters{
		ScreenWidth:  constants.ScreenWidth,
		ScreenHeight: constants.ScreenHeight,
	})
	err := instance.AddEntity(
		ecs.EntityID(0),
		&ecs.Position{
			Vector: geom.Vector{
				X: constants.ScreenWidth / 2,
				Y: constants.ScreenHeight / 2,
			},
		},
		new(ecs.Motion),
		new(ecs.Helm),
		new(ecs.Sprite),
	)
	fmt.Println(err)
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err := ebiten.RunGame(instance); err != nil {
		log.Fatal(err)
	}
}
