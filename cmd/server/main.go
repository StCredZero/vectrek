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
)

// HeadlessServer is a server that doesn't require a display
type HeadlessServer struct {
	Instance *ecs.Instance
	done     chan bool
}

// NewHeadlessServer creates a new HeadlessServer
func NewHeadlessServer() *HeadlessServer {
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
		log.Fatalf("fatal error: %v", err)
	}
	
	return &HeadlessServer{
		Instance: instance,
		done:     make(chan bool, 10),
	}
}

// Start starts the server
func (s *HeadlessServer) Start() {
	go s.Instance.RunServer(s.done)
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
	server := ecs.NewServer(headlessServer.Instance)
	
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
