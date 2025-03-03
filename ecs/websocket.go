package ecs

import (
	"encoding/json"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	WebsocketPort = "8080"
)

// WebsocketMessage wraps ComponentMessage for JSON serialization
type WebsocketMessage struct {
	Entity  ecstypes.EntityID `json:"entity"`
	Payload json.RawMessage   `json:"payload"`
	Type    string            `json:"type"`
}

// WebsocketSender implements the ecstypes.Sender interface for websocket communication
type WebsocketSender struct {
	conn      *websocket.Conn
	writeMu   sync.Mutex
	connected bool
}

// NewWebsocketSender creates a new WebsocketSender
func NewWebsocketSender(conn *websocket.Conn) *WebsocketSender {
	return &WebsocketSender{
		conn:      conn,
		connected: conn != nil,
	}
}

// Send sends a message over the websocket
func (ws *WebsocketSender) Send(msg ecstypes.ComponentMessage) {
	if !ws.connected {
		return
	}

	var payloadType string
	var payloadBytes []byte
	var err error

	switch payload := msg.Payload.(type) {
	case HelmInput:
		payloadType = "HelmInput"
		payloadBytes, err = json.Marshal(payload)
	case SyncInput:
		payloadType = "SyncInput"
		payloadBytes, err = json.Marshal(payload)
	default:
		log.Printf("Unknown payload type: %T", msg.Payload)
		return
	}

	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	wsMsg := WebsocketMessage{
		Entity:  msg.Entity,
		Payload: payloadBytes,
		Type:    payloadType,
	}

	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	if err := ws.conn.WriteJSON(wsMsg); err != nil {
		log.Printf("Error sending message: %v", err)
		ws.connected = false
	}
}

// WebsocketReceiver implements the ecstypes.Receiver interface for websocket communication
type WebsocketReceiver struct {
	conn       *websocket.Conn
	messages   []ecstypes.ComponentMessage
	messagesMu sync.Mutex
	connected  bool
}

// NewWebsocketReceiver creates a new WebsocketReceiver
func NewWebsocketReceiver(conn *websocket.Conn) *WebsocketReceiver {
	receiver := &WebsocketReceiver{
		conn:      conn,
		messages:  make([]ecstypes.ComponentMessage, 0),
		connected: conn != nil,
	}

	if conn != nil {
		go receiver.readMessages()
	}

	return receiver
}

// readMessages reads messages from the websocket and stores them
func (wr *WebsocketReceiver) readMessages() {
	for {
		var wsMsg WebsocketMessage
		if err := wr.conn.ReadJSON(&wsMsg); err != nil {
			log.Printf("Error reading message: %v", err)
			wr.connected = false
			return
		}

		var msg ecstypes.ComponentMessage
		msg.Entity = wsMsg.Entity

		switch wsMsg.Type {
		case "HelmInput":
			var payload HelmInput
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling HelmInput: %v", err)
				continue
			}
			msg.Payload = payload
		case "SyncInput":
			var payload SyncInput
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling SyncInput: %v", err)
				continue
			}
			msg.Payload = payload
		default:
			log.Printf("Unknown message type: %s", wsMsg.Type)
			continue
		}

		wr.messagesMu.Lock()
		wr.messages = append(wr.messages, msg)
		wr.messagesMu.Unlock()
	}
}

// Receive returns the next message from the websocket
func (wr *WebsocketReceiver) Receive() (ecstypes.ComponentMessage, bool) {
	wr.messagesMu.Lock()
	defer wr.messagesMu.Unlock()

	if len(wr.messages) == 0 {
		return ecstypes.ComponentMessage{}, false
	}

	msg := wr.messages[0]
	wr.messages = wr.messages[1:]
	return msg, true
}

// Server wraps an Instance and provides websocket functionality
type Server struct {
	Instance *Instance
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	clientsMu sync.Mutex
}

// NewServer creates a new Server
func NewServer(instance *Instance) *Server {
	return &Server{
		Instance: instance,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections in development
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start starts the server
func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebsocket)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleWebsocket handles websocket connections
func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	s.clientsMu.Lock()
	s.clients[conn] = true
	s.clientsMu.Unlock()

	// Set up the sender and receiver for this client
	sender := NewWebsocketSender(conn)
	receiver := NewWebsocketReceiver(conn)

	// Update the instance to use this client's sender and receiver
	s.Instance.SetSender(sender)
	s.Instance.SetReceiver(receiver)

	// Keep the connection open until it's closed
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
	}

	s.clientsMu.Lock()
	delete(s.clients, conn)
	s.clientsMu.Unlock()
}
