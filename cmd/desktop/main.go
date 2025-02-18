package main

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecs"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"time"
)

func main() {
	instance := ecs.NewInstance(ecs.Parameters{
		ScreenWidth:  constants.ScreenWidth,
		ScreenHeight: constants.ScreenHeight,
	})
	err := instance.AddEntity(
		ecstypes.EntityID(0),
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
	if err != nil {
		log.Fatalf("fatal error: %w", err)
	}
	pipe := ecs.NewPipe()
	instance.SetPipe(pipe)

	var sender ecs.Sender = pipe
	var currentInput ecs.HelmInput
	go func() {
		for {
			var shipInput ecs.HelmInput
			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
				shipInput.Left = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
				shipInput.Right = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				shipInput.Thrust = true
			}

			if currentInput != shipInput {
				currentInput = shipInput
				sender.Send(ecs.ComponentMessage{
					Entity:  ecstypes.EntityID(0),
					Payload: currentInput,
				})
			}
			time.Sleep(1 / 60 * time.Second)
		}
	}()

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	fmt.Println("about to run game")
	if err = ebiten.RunGame(instance); err != nil {
		log.Fatalf("fatal error: %w", err)
	}
}
