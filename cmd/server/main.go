package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/StCredZero/vectrek/internal/game"
	"github.com/pion/webrtc/v3"
)

// GameServer handles the game state and WebRTC connections
type GameServer struct {
	world    *ecs.World
	peerConn *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
}

func NewGameServer() (*GameServer, error) {
	// Configure WebRTC with UDP semantics
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	server := &GameServer{
		world:    &ecs.World{},
		peerConn: peerConn,
	}

	// Create data channels for game state and input
	stateChannel, err := peerConn.CreateDataChannel("state", &webrtc.DataChannelInit{
		Ordered: webrtc.Bool(false), // false = unordered delivery (UDP-like)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create state channel: %v", err)
	}
	server.dataChannel = stateChannel

	inputChannel, err := peerConn.CreateDataChannel("input", &webrtc.DataChannelInit{
		Ordered: webrtc.Bool(false), // false = unordered delivery (UDP-like)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create input channel: %v", err)
	}

	// Handle input from client
	inputChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		var input struct{ Left, Right, Up bool }
		if err := json.Unmarshal(msg.Data, &input); err != nil {
			log.Printf("Failed to parse input: %v", err)
			return
		}

		// Apply input to ship control system
		if shipSystem := server.world.Systems()[0].(*game.ShipControlSystem); shipSystem != nil {
			shipSystem.HandleInput(input, 1.0/60.0) // Use fixed timestep for predictability
		}
	})

	// Set up the game world with server-side systems only
	server.world.AddSystem(&game.ShipControlSystem{})

	return server, nil
}

func (s *GameServer) Start() error {
	// Game loop
	ticker := time.NewTicker(time.Second / 60) // 60 FPS
	defer ticker.Stop()

	for range ticker.C {
		// Update game state
		s.world.Update(1.0 / 60.0)

		// Send position updates to client
		if s.dataChannel.ReadyState() == webrtc.DataChannelStateOpen {
			// Collect entity state data
			state := struct {
				Entities map[string]struct {
					X, Y     float32
					Rotation float32
					Type     string
				}
			}{
				Entities: make(map[string]struct {
					X, Y     float32
					Rotation float32
					Type     string
				}),
			}

			for _, e := range s.world.Entities() {
				if pos, ok := e.GetComponent(&game.PositionComponent{}).(*game.PositionComponent); ok {
					rot := float32(0)
					if rotComp, ok := e.GetComponent(&game.RotationComponent{}).(*game.RotationComponent); ok {
						rot = rotComp.Angle
					}

					// Determine entity type
					entityType := "rectangle"
					if e.ID() == "Ship" {
						entityType = "ship"
					}

					state.Entities[e.ID()] = struct {
						X, Y     float32
						Rotation float32
						Type     string
					}{
						X:        pos.X,
						Y:        pos.Y,
						Rotation: rot,
						Type:     entityType,
					}
				}
			}

			// Send state update
			data, err := json.Marshal(state)
			if err != nil {
				log.Printf("Failed to marshal state: %v", err)
				continue
			}

			if err := s.dataChannel.Send(data); err != nil {
				log.Printf("Failed to send state update: %v", err)
				// Continue even if send fails (UDP semantics)
			}
		}
	}

	return nil
}

func main() {
	server, err := NewGameServer()
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
