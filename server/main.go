package main

import (
	"log"
	"net/http"

	"github.com/the-js-developer/voice-recorder/app/handler"
)

func main() {
	// Initialize WebSocket handler
	wsHandler := handler.NewHandler()

	// Setup routes
	http.HandleFunc("/stream", wsHandler.Handle)

	port := ":8080"
	log.Printf("WebSocket server starting on %s", port)

	// Start server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
