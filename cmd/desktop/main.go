package main

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecs"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/geom"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net/url"
)

func newClientInstance(sender ecstypes.Sender, receiver ecstypes.Receiver) *ecs.Instance {
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
		log.Fatalf("fatal error: %v", err)
	}
	instance.SetSender(sender)
	instance.SetReceiver(receiver)
	return instance
}

func main() {
	// Connect to the websocket server
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", ecs.WebsocketPort), Path: "/ws"}
	log.Printf("Connecting to %s", u.String())
	
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// Create the sender and receiver
	sender := ecs.NewWebsocketSender(conn)
	receiver := ecs.NewWebsocketReceiver(conn)

	// Create the client instance
	clientInstance := newClientInstance(sender, receiver)

	// Run the game
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err = ebiten.RunGame(clientInstance); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}
