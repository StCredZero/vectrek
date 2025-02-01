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

	// Set up data channel handler for state updates
	s.peerConn.OnDataChannel(func(d *webrtc.DataChannel) {
		if d.Label() == "state" {
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				// Parse state update
				var state struct {
					Entities map[string]struct {
						X, Y     float32
						Rotation float32
						Type     string
					}
				}
				if err := json.Unmarshal(msg.Data, &state); err != nil {
					log.Printf("Failed to parse state update: %v", err)
					return
				}

				// Update or create entities based on server state
				for id, entityState := range state.Entities {
					entity, exists := s.entityMap[id]
					if !exists {
						// Create new entity
						entity = ecs.NewEntity(id)
						entity.AddComponent(&game.PositionComponent{X: entityState.X, Y: entityState.Y})
						entity.AddComponent(&game.RotationComponent{Angle: entityState.Rotation})

						// Add render component based on entity type
						var renderComp *game.RenderComponent
						switch entityState.Type {
						case "ship":
							renderComp = game.CreateRenderComponent(game.EntityTypeShip)
						case "rectangle":
							renderComp = game.CreateRenderComponent(game.EntityTypeRectangle)
						default:
							log.Printf("Unknown entity type: %s", entityState.Type)
							continue
						}
						entity.AddComponent(renderComp)

						s.world.AddEntity(entity)
						s.entityMap[id] = entity
					} else {
						// Update existing entity
						if posComp, ok := entity.GetComponent(&game.PositionComponent{}).(*game.PositionComponent); ok {
							posComp.X = entityState.X
							posComp.Y = entityState.Y
						}
						if rotComp, ok := entity.GetComponent(&game.RotationComponent{}).(*game.RotationComponent); ok {
							rotComp.Angle = entityState.Rotation
						}
					}
				}

				// Remove entities that no longer exist in server state
				for id, entity := range s.entityMap {
					if _, exists := state.Entities[id]; !exists {
						s.world.RemoveEntity(*entity)
						delete(s.entityMap, id)
					}
				}
			})
		}
	})

	// Register input controls
	engo.Input.RegisterButton("ArrowLeft", engo.KeyArrowLeft)
	engo.Input.RegisterButton("ArrowRight", engo.KeyArrowRight)
	engo.Input.RegisterButton("ArrowUp", engo.KeyArrowUp)

	// Add render system for graphics (client-side only)
	s.world.AddSystem(&game.RenderSystem{})
}

// Update handles input and sends it to the server
func (s *mainScene) Update(dt float32) {
	// Check input state
	input := struct{ Left, Right, Up bool }{
		Left:  engo.Input.Button("ArrowLeft").Down(),
		Right: engo.Input.Button("ArrowRight").Down(),
		Up:    engo.Input.Button("ArrowUp").Down(),
	}

	// Send input to server if connection is ready
	if s.peerConn != nil && s.peerConn.ConnectionState() == webrtc.PeerConnectionStateConnected {
		// Find input channel
		for _, dc := range s.peerConn.DataChannels() {
			if dc.Label() == "input" && dc.ReadyState() == webrtc.DataChannelStateOpen {
				data, err := json.Marshal(input)
				if err != nil {
					log.Printf("Failed to marshal input: %v", err)
					return
				}
				
				if err := dc.Send(data); err != nil {
					log.Printf("Failed to send input: %v", err)
				}
				break
			}
		}
	}
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
