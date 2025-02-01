package main

import (
	"encoding/json"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/StCredZero/vectrek/internal/game"
	"github.com/pion/webrtc/v3"
)

type mainScene struct {
	world      *ecs.World
	peerConn   *webrtc.PeerConnection
	entityMap  map[string]*ecs.Entity
}

func (s *mainScene) Type() string { return "VecTrek Scene" }

func (s *mainScene) Preload() {}

func (s *mainScene) Setup(u engo.Updater) {
	var err error
	s.world, _ = u.(*ecs.World)
	s.entityMap = make(map[string]*ecs.Entity)

	// Configure WebRTC
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create peer connection
	s.peerConn, err = webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatalf("Failed to create peer connection: %v", err)
	}

	// Set up data channel handler
	s.peerConn.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			// Parse position updates
			var positions map[string]struct {
				X, Y     float32
				Rotation float32
			}
			if err := json.Unmarshal(msg.Data, &positions); err != nil {
				log.Printf("Failed to parse position update: %v", err)
				return
			}

			// Update entity positions
			for id, pos := range positions {
				entity, exists := s.entityMap[id]
				if !exists {
					// Create new entity if it doesn't exist
					entity = ecs.NewEntity(id)
					entity.AddComponent(&game.PositionComponent{X: pos.X, Y: pos.Y})
					entity.AddComponent(&game.RotationComponent{Angle: pos.Rotation})
					entity.AddComponent(&game.RenderComponent{
						Points: []float32{
							0, -20, // nose
							15, 20,  // right
							-15, 20, // left
						},
						Color: struct{ R, G, B float32 }{1, 1, 1},
					})
					s.world.AddEntity(entity)
					s.entityMap[id] = entity
				} else {
					// Update existing entity
					if posComp, ok := entity.GetComponent(&game.PositionComponent{}).(*game.PositionComponent); ok {
						posComp.X = pos.X
						posComp.Y = pos.Y
					}
					if rotComp, ok := entity.GetComponent(&game.RotationComponent{}).(*game.RotationComponent); ok {
						rotComp.Angle = pos.Rotation
					}
				}
			}
		})
	})

	// Add render system for graphics
	s.world.AddSystem(&game.RenderSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:  "VecTrek Client",
		Width:  800,
		Height: 600,
	}

	scene := &mainScene{}
	engo.Run(opts, scene)
}
