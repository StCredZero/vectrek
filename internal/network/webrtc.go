package network

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pion/webrtc/v3"
)

// GameState represents the state to be sent over the network
type GameState struct {
	Entities map[string]EntityState `json:"entities"`
}

// EntityState represents an entity's network state
type EntityState struct {
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Rotation float32 `json:"rotation"`
	Type     string  `json:"type"`
}

// Server handles WebRTC connections and state broadcasting
type Server struct {
	peerConns   map[string]*webrtc.PeerConnection
	dataChannels map[string]*webrtc.DataChannel
}

// NewServer creates a new WebRTC server
func NewServer() (*Server, error) {
	return &Server{
		peerConns:    make(map[string]*webrtc.PeerConnection),
		dataChannels: make(map[string]*webrtc.DataChannel),
	}, nil
}

// CreatePeerConnection creates a new WebRTC peer connection with UDP semantics
func (s *Server) CreatePeerConnection(peerID string) (*webrtc.PeerConnection, error) {
	// Configure WebRTC with STUN server for NAT traversal
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create peer connection
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	// Create data channel with UDP-like semantics
	dataChannel, err := peerConn.CreateDataChannel("game", &webrtc.DataChannelInit{
		Ordered:        webrtc.Bool(false), // Disable ordering for UDP-like behavior
		MaxRetransmits: webrtc.Uint16(0),   // No retransmissions
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data channel: %v", err)
	}

	// Store connections
	s.peerConns[peerID] = peerConn
	s.dataChannels[peerID] = dataChannel

	return peerConn, nil
}

// BroadcastState sends game state to all connected clients
func (s *Server) BroadcastState(state *GameState) {
	// Convert state to JSON
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Failed to marshal game state: %v", err)
		return
	}

	// Send to all connected clients
	for peerID, dc := range s.dataChannels {
		if dc.ReadyState() == webrtc.DataChannelStateOpen {
			if err := dc.Send(data); err != nil {
				log.Printf("Failed to send state to peer %s: %v", peerID, err)
				// Continue even if send fails (UDP semantics)
			}
		}
	}
}

// Client handles WebRTC connection to the server
type Client struct {
	peerConn    *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
	OnState     func(*GameState)
}

// NewClient creates a new WebRTC client
func NewClient() (*Client, error) {
	// Configure WebRTC with STUN server
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create peer connection
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	client := &Client{
		peerConn: peerConn,
	}

	// Handle incoming data channels
	peerConn.OnDataChannel(func(dc *webrtc.DataChannel) {
		client.dataChannel = dc
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			if client.OnState != nil {
				var state GameState
				if err := json.Unmarshal(msg.Data, &state); err != nil {
					log.Printf("Failed to unmarshal game state: %v", err)
					return
				}
				client.OnState(&state)
			}
		})
	})

	return client, nil
}
