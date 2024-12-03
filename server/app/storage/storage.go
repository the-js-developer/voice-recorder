package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RecorderStorageInterface interface {
	AppendChunk(chunk []byte)
	SaveChunk() error
	Reset()
}

// AudioRecorder manages audio chunk storage
type Recorder struct {
	chunks []byte
	mu          sync.Mutex
}

// NewAudioRecorder creates a new audio recorder
func NewAudioRecorder() RecorderStorageInterface {
	return &Recorder{
		chunks: []byte{},
	}
}

// AppendChunk adds a new audio chunk
func (r *Recorder) AppendChunk(chunk []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chunks = append(r.chunks, chunk...)
}

func checkFileIntegrity(filename string) error {
    // Basic file size and read check
    fileInfo, err := os.Stat(filename)
    if err != nil {
        return fmt.Errorf("file stat error: %v", err)
    }

    if fileInfo.Size() == 0 {
        return fmt.Errorf("file is empty")
    }

    // Try to read the entire file
    _, err = os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("file read error: %v", err)
    }

    return nil
}

func diagnoseAudioFile(filename string) {
    // Check basic file integrity
    err := checkFileIntegrity(filename)
    if err != nil {
        log.Printf("File integrity check failed: %v", err)
    }

    // Additional diagnostics
    fileInfo, _ := os.Stat(filename)
    log.Printf("File Details: \n"+
        "Path: %s\n"+
        "Size: %d bytes\n"+
        "Modification Time: %v", 
        filename, 
        fileInfo.Size(), 
        fileInfo.ModTime(),
    )
}

// SaveChunk writes accumulated audio chunks to a file
func (r *Recorder) SaveChunk() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Skip saving if no chunks
	if len(r.chunks) == 0 {
		return nil
	}

	// Ensure recordings directory exists
	recordingsDir := "recordings"
	if err := os.MkdirAll(recordingsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create recordings directory: %v", err)
	}

	// Generate unique filename
	filename := filepath.Join(
		recordingsDir, 
		fmt.Sprintf("recording_%d.wav", time.Now().UnixNano()),
	)
	
	log.Printf("Chunk size before saving: %d bytes", len(r.chunks))

	// Create a copy of current chunks to save
    chunksCopy := make([]byte, len(r.chunks))
    copy(chunksCopy, r.chunks)

    err := os.WriteFile(filename, chunksCopy, 0644)

	diagnoseAudioFile(filename)
    if err != nil {
        return err
    }

    // Clear chunks after successful save
    r.chunks = []byte{}
    log.Println("Chunks reset after saving")

    return nil
}

// Reset clears accumulated audio chunks
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chunks = []byte{}
}