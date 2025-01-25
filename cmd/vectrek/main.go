package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/StCredZero/vectrek/internal/game"
)

type mainScene struct{}

func (*mainScene) Type() string { return "VecTrek Scene" }

func (*mainScene) Preload() {}

func (*mainScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	
	world.AddSystem(&game.ShipControlSystem{})
	world.AddSystem(&game.RenderSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:  "VecTrek",
		Width:  800,
		Height: 600,
	}
	
	engo.Run(opts, &mainScene{})
}
