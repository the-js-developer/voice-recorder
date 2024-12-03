package handler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/the-js-developer/voice-recorder/app/service"
)

// Handler manages WebSocket connections
type Handler struct {
	upgrader websocket.Upgrader
	sessions map[string]*service.RecordingSession
	mu       sync.Mutex
}

// NewHandler creates a new WebSocket handler
func NewHandler() *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		sessions: make(map[string]*service.RecordingSession),
	}
}

// Handle manages incoming WebSocket connections
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Create new recording session
	session := service.NewRecordingSession()

	// Register session
	h.mu.Lock()
	h.sessions[session.GetID()] = session
	h.mu.Unlock()

	// Create a channel to control chunk saving
    saveChunkChan := make(chan struct{}, 1)

	// recorded chunk saving ticker
	// recordedChunkTimer := time.NewTicker(30 * time.Second)
	// defer recordedChunkTimer.Stop()

	dynamicTicker := time.NewTimer(func() time.Duration {
		// Calculate remaining time until next 30-second interval
		elapsed := time.Since(session.StartTime)
		nextInterval := (elapsed/time.Second + 1) * time.Second
		remainingTime := nextInterval - elapsed
		
		if remainingTime < 0 {
			remainingTime = 30 * time.Second
		}
		
		return remainingTime
	}())

	defer dynamicTicker.Stop()

	// Chunk saving goroutine
    go func() {
        for {
            select {
            case <-dynamicTicker.C:
				log.Println("session", time.Since(session.StartTime))
                saveChunkChan <- struct{}{}
            }
        }
    }()

	// Chunk saving processor
    go func() {
        for range saveChunkChan {
            log.Println("Attempting to save chunk:")
            err := session.TrySaveChunk()
            if err != nil {
                log.Println("Error saving chunk:", err)
            }

			dynamicTicker.Reset(30 * time.Second)
        }
    }()

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// Add incoming audio chunk
		session.AddChunk(message)
	}

	// Save final chunk
	if err := session.SaveFinalChunk(); err != nil {
		log.Println("Error saving final chunk:", err)
	}

	// Cleanup session
	h.mu.Lock()
	delete(h.sessions, session.ID)
	h.mu.Unlock()
}
