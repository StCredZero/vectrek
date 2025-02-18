package main

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecs"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func newServerInstance(inputPipe ecstypes.Receiver, outputPipe ecstypes.Sender) *ecs.Instance {
	instance := ecs.NewInstance(ecs.Parameters{
		ScreenWidth:  constants.ScreenWidth,
		ScreenHeight: constants.ScreenHeight,
	})
	instance.Name = "Server"
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
		new(ecs.SyncSender),
	)
	if err != nil {
		log.Fatalf("fatal error: %w", err)
	}
	instance.SetReceiver(inputPipe)
	instance.SetSender(outputPipe)
	return instance
}

func newClientInstance(inputPipe ecstypes.Receiver, outputPipe ecstypes.Sender) *ecs.Instance {
	instance := ecs.NewInstance(ecs.Parameters{
		ScreenWidth:  constants.ScreenWidth,
		ScreenHeight: constants.ScreenHeight,
	})
	instance.Name = "Client"
	err := instance.AddEntity(
		ecstypes.EntityID(0),
		&ecs.Position{
			Vector: geom.Vector{
				X: constants.ScreenWidth / 2,
				Y: constants.ScreenHeight / 2,
			},
		},
		new(ecs.Motion),
		new(ecs.Sprite),
		new(ecs.Player),
		new(ecs.SyncReceiver),
	)
	if err != nil {
		log.Fatalf("fatal error: %w", err)
	}
	instance.SetReceiver(inputPipe)
	instance.SetSender(outputPipe)
	return instance
}

func main() {
	var err error
	var serverReceiver = ecs.NewPipe()
	var serverSender = ecs.NewPipe()
	serverInstance := newServerInstance(serverReceiver, serverSender)
	clientInstance := newClientInstance(serverSender, serverReceiver)

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	fmt.Println("about to run server")
	done := make(chan bool, 10)
	go serverInstance.RunServer(done)
	fmt.Println("about to run game")
	if err = ebiten.RunGame(clientInstance); err != nil {
		log.Fatalf("fatal error: %w", err)
	}
}
