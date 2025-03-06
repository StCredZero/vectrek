package main

import (
	"fmt"
	"github.com/StCredZero/vectrek/constants"
	"github.com/StCredZero/vectrek/ecs"
	"github.com/StCredZero/vectrek/ecstypes"
	"github.com/StCredZero/vectrek/geom"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerInstance is a minimal implementation of the game instance for server-side use
// that doesn't depend on Ebiten or any graphical capabilities
type ServerInstance struct {
	entities map[ecstypes.EntityID]struct{}
	name     string
	
	position     *ecs.SMSystem[ecs.Position]
	motion       *ecs.SMSystem[ecs.Motion]
	helm         *ecs.SMSystem[ecs.Helm]
	syncReceiver *ecs.SMSystem[ecs.SyncReceiver]
	syncSender   *ecs.SMSystem[ecs.SyncSender]
	
	counter    uint64
	parameters ecs.Parameters
	
	receiver ecstypes.Receiver
	sender   ecstypes.Sender
}

// NewServerInstance creates a new ServerInstance
func NewServerInstance() *ServerInstance {
	result := &ServerInstance{
		entities: make(map[ecstypes.EntityID]struct{}),
		name:     "Server",
		parameters: ecs.Parameters{
			ScreenWidth:  constants.ScreenWidth,
			ScreenHeight: constants.ScreenHeight,
		},
	}
	
	// Initialize systems
	result.position = ecs.NewSMSystem[ecs.Position](func(each *ecs.Position) error {
		// We're creating a minimal instance that doesn't use the full Instance struct
		// So we need to create a temporary Instance to pass to the Update method
		tempInstance := &ecs.Instance{
			Position: result.position,
			Motion:   result.motion,
			Helm:     result.helm,
			SyncReceiver: result.syncReceiver,
			SyncSender:   result.syncSender,
			Parameters:   result.parameters,
			Receiver:     result.receiver,
			Sender:       result.sender,
		}
		return each.Update(tempInstance)
	})
	
	result.motion = ecs.NewSMSystem[ecs.Motion](func(each *ecs.Motion) error {
		tempInstance := &ecs.Instance{
			Position: result.position,
			Motion:   result.motion,
			Helm:     result.helm,
			SyncReceiver: result.syncReceiver,
			SyncSender:   result.syncSender,
			Parameters:   result.parameters,
			Receiver:     result.receiver,
			Sender:       result.sender,
		}
		return each.Update(tempInstance)
	})
	
	result.helm = ecs.NewSMSystem[ecs.Helm](func(each *ecs.Helm) error {
		tempInstance := &ecs.Instance{
			Position: result.position,
			Motion:   result.motion,
			Helm:     result.helm,
			SyncReceiver: result.syncReceiver,
			SyncSender:   result.syncSender,
			Parameters:   result.parameters,
			Receiver:     result.receiver,
			Sender:       result.sender,
		}
		return each.Update(tempInstance)
	})
	
	result.syncReceiver = ecs.NewSMSystem[ecs.SyncReceiver](func(each *ecs.SyncReceiver) error {
		tempInstance := &ecs.Instance{
			Position: result.position,
			Motion:   result.motion,
			Helm:     result.helm,
			SyncReceiver: result.syncReceiver,
			SyncSender:   result.syncSender,
			Parameters:   result.parameters,
			Receiver:     result.receiver,
			Sender:       result.sender,
		}
		return each.Update(tempInstance)
	})
	
	result.syncSender = ecs.NewSMSystem[ecs.SyncSender](func(each *ecs.SyncSender) error {
		tempInstance := &ecs.Instance{
			Position: result.position,
			Motion:   result.motion,
			Helm:     result.helm,
			SyncReceiver: result.syncReceiver,
			SyncSender:   result.syncSender,
			Parameters:   result.parameters,
			Receiver:     result.receiver,
			Sender:       result.sender,
		}
		return each.Update(tempInstance)
	})
	
	// Initialize with dummy sender and receiver
	result.sender = &ecs.DummySender{}
	result.receiver = &ecs.DummyReceiver{}
	
	return result
}

// AddEntity adds an entity to the server instance
func (s *ServerInstance) AddEntity(entity ecstypes.EntityID, components ...ecstypes.Component) error {
	s.entities[entity] = struct{}{}
	
	// Create a temporary Instance to pass to the Init method
	tempInstance := &ecs.Instance{
		Position: s.position,
		Motion:   s.motion,
		Helm:     s.helm,
		SyncReceiver: s.syncReceiver,
		SyncSender:   s.syncSender,
		Parameters:   s.parameters,
		Receiver:     s.receiver,
		Sender:       s.sender,
	}
	
	for _, component := range components {
		if err := component.Init(tempInstance, entity); err != nil {
			return err
		}
	}
	
	return nil
}

// Update updates the server instance
func (s *ServerInstance) Update() error {
	s.counter++
	
	// Process incoming messages
	var hasMessage bool
	var msg ecstypes.ComponentMessage
	for {
		if msg, hasMessage = s.receiver.Receive(); !hasMessage {
			break
		}
		switch obj := msg.Payload.(type) {
		case ecs.HelmInput:
			if helm, ok := s.helm.GetComponent(msg.Entity); ok {
				helm.Input = obj
			}
		case ecs.SyncInput:
			if sync, ok := s.syncReceiver.GetComponent(msg.Entity); ok {
				sync.Input <- obj
			}
		default:
		}
	}
	
	// Update systems in the correct order
	s.helm.Iterate()
	s.motion.Iterate()
	s.syncSender.Iterate()
	s.syncReceiver.Iterate()
	
	return nil
}

// SetSender sets the sender for the server instance
func (s *ServerInstance) SetSender(sender ecstypes.Sender) {
	s.sender = sender
}

// SetReceiver sets the receiver for the server instance
func (s *ServerInstance) SetReceiver(receiver ecstypes.Receiver) {
	s.receiver = receiver
}

// GetSender returns the sender for the server instance
func (s *ServerInstance) GetSender() ecstypes.Sender {
	return s.sender
}

// GetReceiver returns the receiver for the server instance
func (s *ServerInstance) GetReceiver() ecstypes.Receiver {
	return s.receiver
}

// HeadlessServer wraps a ServerInstance and provides a game loop
type HeadlessServer struct {
	Instance *ServerInstance
	done     chan bool
	ticker   *time.Ticker
}

// NewHeadlessServer creates a new HeadlessServer
func NewHeadlessServer() *HeadlessServer {
	instance := NewServerInstance()
	
	// Add the player entity
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
		log.Fatalf("fatal error: %v", err)
	}
	
	return &HeadlessServer{
		Instance: instance,
		done:     make(chan bool, 10),
		ticker:   time.NewTicker(16667 * time.Microsecond), // ~60 FPS
	}
}

// Start starts the server
func (s *HeadlessServer) Start() {
	go func() {
		defer s.ticker.Stop()
		for {
			select {
			case <-s.ticker.C:
				s.Instance.Update()
			case <-s.done:
				return
			}
		}
	}()
}

// Stop stops the server
func (s *HeadlessServer) Stop() {
	s.done <- true
}

func main() {
	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create and start the headless server
	headlessServer := NewHeadlessServer()
	headlessServer.Start()

	// Create and start the websocket server
	server := ecs.NewServer(
		&ecs.Instance{
			Name:         "Server",
			Position:     headlessServer.Instance.position,
			Motion:       headlessServer.Instance.motion,
			Helm:         headlessServer.Instance.helm,
			SyncReceiver: headlessServer.Instance.syncReceiver,
			SyncSender:   headlessServer.Instance.syncSender,
			Parameters:   headlessServer.Instance.parameters,
			Receiver:     headlessServer.Instance.receiver,
			Sender:       headlessServer.Instance.sender,
		},
	)
	
	// Start the websocket server in a goroutine
	addr := fmt.Sprintf(":%s", ecs.WebsocketPort)
	log.Printf("Starting server on %s", addr)
	
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(addr)
	}()

	// Wait for either a signal or an error
	select {
	case <-sigs:
		log.Println("Received shutdown signal")
	case err := <-errChan:
		log.Fatalf("Error starting server: %v", err)
	}

	// Clean up
	headlessServer.Stop()
	log.Println("Server shutdown complete")
}
