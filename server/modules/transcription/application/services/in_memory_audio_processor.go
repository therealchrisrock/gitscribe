package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/services"
)

// inMemoryAudioSession represents an active audio processing session in memory
type inMemoryAudioSession struct {
	ID          string
	Metadata    services.AudioStreamMetadata
	Options     services.AudioProcessingOptions
	StartTime   int64
	Chunks      []services.AudioChunk
	Status      entities.TranscriptionStatus
	TotalBytes  int
	ChunkCount  int
	LastChunkAt int64
	mu          sync.RWMutex
}

// InMemoryAudioProcessor implements the AudioProcessor interface with in-memory session management
// This is a lightweight implementation suitable for testing, development, and scenarios without database persistence
type InMemoryAudioProcessor struct {
	sessions map[string]*inMemoryAudioSession
	mu       sync.RWMutex
}

// NewInMemoryAudioProcessor creates a new in-memory audio processor
func NewInMemoryAudioProcessor() *InMemoryAudioProcessor {
	return &InMemoryAudioProcessor{
		sessions: make(map[string]*inMemoryAudioSession),
	}
}

// GetSupportedModes returns the processing modes supported by this processor
func (s *InMemoryAudioProcessor) GetSupportedModes() []services.ProcessingMode {
	return []services.ProcessingMode{
		services.RealTimeMode,
		services.BatchMode,
	}
}

// StartSession initializes a new audio processing session
func (s *InMemoryAudioProcessor) StartSession(ctx context.Context, metadata services.AudioStreamMetadata, options services.AudioProcessingOptions) (string, error) {
	sessionID := metadata.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	// Set default processing mode if not specified
	if metadata.Mode == "" {
		metadata.Mode = services.RealTimeMode
	}

	// Set default options based on processing mode
	if options.Mode == "" {
		options.Mode = metadata.Mode
	}

	session := &inMemoryAudioSession{
		ID:         sessionID,
		Metadata:   metadata,
		Options:    options,
		StartTime:  time.Now().Unix(),
		Chunks:     make([]services.AudioChunk, 0),
		Status:     entities.Processing,
		TotalBytes: 0,
		ChunkCount: 0,
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	log.Printf("Started audio session %s with mode: %s, provider: %s", sessionID, metadata.Mode, options.Provider)

	// TODO: For batch mode, we might want to initialize differently
	// For now, both modes use the same initialization
	if options.Mode == services.BatchMode {
		log.Printf("Session %s initialized for batch processing", sessionID)
	} else {
		log.Printf("Session %s initialized for real-time processing", sessionID)
	}

	return sessionID, nil
}

// ProcessChunk processes an individual audio chunk
func (s *InMemoryAudioProcessor) ProcessChunk(ctx context.Context, sessionID string, chunk services.AudioChunk) error {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	// Update session status to processing on first chunk
	if session.Status == entities.Pending {
		session.Status = entities.Processing
	}

	// Set chunk metadata
	chunk.SequenceNum = session.ChunkCount + 1
	chunk.Size = len(chunk.Data)
	chunk.Timestamp = time.Now().Unix()

	// Add chunk to session
	session.Chunks = append(session.Chunks, chunk)
	session.ChunkCount++
	session.TotalBytes += chunk.Size
	session.LastChunkAt = chunk.Timestamp

	log.Printf("Processed chunk #%d for session %s: %d bytes (total: %d bytes)",
		chunk.SequenceNum, sessionID, chunk.Size, session.TotalBytes)

	// TODO: Here you would integrate with actual transcription providers
	// For example:
	// - Send chunk to AssemblyAI real-time API
	// - Buffer chunks and send to Whisper API
	// - Process with local Whisper model

	return nil
}

// EndSession finalizes the audio processing session and returns the result
func (s *InMemoryAudioProcessor) EndSession(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
	s.mu.Lock()
	session, exists := s.sessions[sessionID]
	if !exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	// Mark session as completed
	session.Status = entities.Completed
	session.LastChunkAt = time.Now().Unix()
	s.mu.Unlock()

	// Create transcript segments based on processing mode
	var segments []entities.TranscriptSegment

	if session.Options.Mode == services.BatchMode {
		// For batch mode, we would typically submit the complete audio for processing
		// and wait for the result. For now, we'll simulate this.
		log.Printf("Processing session %s in batch mode - would submit complete audio for processing", sessionID)

		// TODO: Implement actual batch processing logic here
		// This is where you would:
		// 1. Concatenate all audio chunks
		// 2. Submit to batch transcription service
		// 3. Wait for or poll for results

		segment := entities.NewTranscriptSegment(
			sessionID,
			"speaker_1",
			fmt.Sprintf("Batch processed transcript for session %s with %d chunks", sessionID, session.ChunkCount),
			0.0,
			float64(session.ChunkCount)*1.0, // Estimate based on chunks
			0.95,
			1,
		)
		segments = append(segments, segment)
	} else {
		// Real-time mode - create segments based on chunks received
		for i := range session.Chunks {
			segment := entities.NewTranscriptSegment(
				sessionID,
				"speaker_1",
				fmt.Sprintf("Real-time transcript segment %d", i+1),
				float64(i)*1.0,
				float64(i+1)*1.0,
				0.85,
				i+1,
			)
			segments = append(segments, segment)
		}
	}

	result := &services.AudioProcessingResult{
		TranscriptionID: sessionID,
		Status:          session.Status,
		Segments:        segments,
		ProcessingMode:  session.Options.Mode,
		Message:         fmt.Sprintf("Session completed with %d chunks processed using %s mode", session.ChunkCount, session.Options.Mode),
	}

	// Clean up session after a delay (in a real implementation, you might want to keep it longer)
	go func() {
		time.Sleep(1 * time.Minute)
		s.mu.Lock()
		delete(s.sessions, sessionID)
		s.mu.Unlock()
		log.Printf("Cleaned up session %s", sessionID)
	}()

	log.Printf("Ended session %s with %d chunks processed (mode: %s)", sessionID, session.ChunkCount, session.Options.Mode)
	return result, nil
}

// GetSessionStatus returns the current status of a session
func (s *InMemoryAudioProcessor) GetSessionStatus(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	result := &services.AudioProcessingResult{
		TranscriptionID: "", // No transcription ID until completed
		Status:          session.Status,
		Message:         fmt.Sprintf("Session active: %d chunks processed (%d bytes)", session.ChunkCount, session.TotalBytes),
	}

	return result, nil
}

// AbortSession cancels an ongoing session and cleans up resources
func (s *InMemoryAudioProcessor) AbortSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.mu.Lock()
	session.Status = entities.Failed
	session.mu.Unlock()

	log.Printf("Aborted audio processing session %s", sessionID)

	// Clean up session
	delete(s.sessions, sessionID)

	return nil
}

// IsSessionActive checks if a session is currently active
func (s *InMemoryAudioProcessor) IsSessionActive(ctx context.Context, sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return false
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	return session.Status == entities.Pending || session.Status == entities.Processing
}

// GetActiveSessionsCount returns the number of currently active sessions
func (s *InMemoryAudioProcessor) GetActiveSessionsCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, session := range s.sessions {
		session.mu.RLock()
		if session.Status == entities.Pending || session.Status == entities.Processing {
			count++
		}
		session.mu.RUnlock()
	}

	return count
}
