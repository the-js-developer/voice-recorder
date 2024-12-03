package service

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/the-js-developer/voice-recorder/app/storage"
)

type RecordingSession struct {
	ID        string
	recorder  *storage.Recorder
	StartTime time.Time
	endTime   time.Time
	mu        sync.Mutex
}

func NewRecordingSession() *RecordingSession {
	return &RecordingSession{
		ID:        uuid.New().String(),
		recorder:  storage.NewAudioRecorder().(*storage.Recorder),
		StartTime: time.Now(),
		endTime:   time.Now(),
	}
}

// AddChunk adds a new audio chunk to the session
func (rs *RecordingSession) AddChunk(chunk []byte) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.recorder.AppendChunk(chunk)
	rs.endTime = time.Now()
}

// TrySaveChunk attempts to save audio chunk if duration exceeded
func (rs *RecordingSession) TrySaveChunk() error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

    duration := time.Since(rs.StartTime)
	if duration >= 30*time.Second-time.Millisecond {
		rs.StartTime = time.Now()
		if err := rs.recorder.SaveChunk(); err != nil {
			return err
		}
	}

	if duration >= 30*time.Second {
        if err := rs.recorder.SaveChunk(); err != nil {
            return err
        }
        // Reset for next chunk
        rs.recorder.Reset()
        rs.StartTime = time.Now() // Reset the start time to the current time
    }
	return nil
}

// SaveFinalChunk saves any remaining audio data
func (rs *RecordingSession) SaveFinalChunk() error {
	return rs.recorder.SaveChunk()
}

// GetID get session id
func (r *RecordingSession) GetID() string {
	return r.ID
}
