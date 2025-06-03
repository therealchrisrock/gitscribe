package providers

import (
	"context"
	"fmt"
	"log"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/services"
)

// MockAssemblyAIProvider implements the AudioProcessor interface with mock functionality
type MockAssemblyAIProvider struct {
	firebaseUploader FirebaseUploader
	sessions         map[string]*MockAssemblyAISession
}

// MockAssemblyAISession tracks an active mock AssemblyAI processing session
type MockAssemblyAISession struct {
	SessionID    string
	TranscriptID string
	Metadata     services.AudioStreamMetadata
	Options      services.AudioProcessingOptions
	AudioChunks  []services.AudioChunk
	Status       entities.TranscriptionStatus
	CreatedAt    time.Time
	FirebaseURL  string
}

// NewMockAssemblyAIProvider creates a new mock AssemblyAI provider
func NewMockAssemblyAIProvider(firebaseUploader FirebaseUploader) *MockAssemblyAIProvider {
	return &MockAssemblyAIProvider{
		firebaseUploader: firebaseUploader,
		sessions:         make(map[string]*MockAssemblyAISession),
	}
}

// GetSupportedModes returns the processing modes supported by AssemblyAI
func (p *MockAssemblyAIProvider) GetSupportedModes() []services.ProcessingMode {
	return []services.ProcessingMode{
		services.RealTimeMode,
		services.BatchMode,
	}
}

// StartSession initializes a new AssemblyAI processing session
func (p *MockAssemblyAIProvider) StartSession(ctx context.Context, metadata services.AudioStreamMetadata, options services.AudioProcessingOptions) (string, error) {
	sessionID := metadata.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("mock_assemblyai_%d", time.Now().UnixNano())
	}

	session := &MockAssemblyAISession{
		SessionID:   sessionID,
		Metadata:    metadata,
		Options:     options,
		AudioChunks: make([]services.AudioChunk, 0),
		Status:      entities.Pending,
		CreatedAt:   time.Now(),
	}

	p.sessions[sessionID] = session

	log.Printf("Mock AssemblyAI session started: %s (mode: %s, diarization: %t)",
		sessionID, options.Mode, options.SpeakerDiarization)

	return sessionID, nil
}

// ProcessChunk processes an individual audio chunk
func (p *MockAssemblyAIProvider) ProcessChunk(ctx context.Context, sessionID string, chunk services.AudioChunk) error {
	session, exists := p.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Add chunk to session
	session.AudioChunks = append(session.AudioChunks, chunk)
	session.Status = entities.Processing

	log.Printf("Mock AssemblyAI: Added chunk #%d to session %s (%d bytes)",
		len(session.AudioChunks), sessionID, len(chunk.Data))

	return nil
}

// EndSession finalizes the session and processes the complete audio
func (p *MockAssemblyAIProvider) EndSession(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
	session, exists := p.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	// Concatenate all audio chunks
	audioData := p.concatenateAudioChunks(session.AudioChunks)

	// Upload to Firebase Storage (if uploader is available)
	if p.firebaseUploader != nil {
		firebaseURL, err := p.firebaseUploader.UploadAudio(ctx, audioData, session.Metadata.MeetingID, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to upload audio to Firebase: %w", err)
		}
		session.FirebaseURL = firebaseURL
	} else {
		// For testing without Firebase
		session.FirebaseURL = fmt.Sprintf("mock://test-bucket/meetings/%s/audio/%s.wav", session.Metadata.MeetingID, sessionID)
		log.Printf("Mock: Skipping Firebase upload (no uploader configured)")
	}

	// Mock transcript ID
	session.TranscriptID = fmt.Sprintf("mock_transcript_%d", time.Now().UnixNano())

	// Simulate processing delay
	time.Sleep(2 * time.Second)

	// Generate mock transcript segments with speaker diarization
	segments := p.generateMockTranscriptSegments(session)
	session.Status = entities.Completed

	result := &services.AudioProcessingResult{
		TranscriptionID: session.TranscriptID,
		Status:          session.Status,
		Segments:        segments,
		ProcessingMode:  session.Options.Mode,
		Message:         fmt.Sprintf("Mock transcription completed with %d segments", len(segments)),
		FirebaseURL:     session.FirebaseURL,
	}

	// Clean up session after delay
	go func() {
		time.Sleep(5 * time.Minute)
		delete(p.sessions, sessionID)
	}()

	return result, nil
}

// GetSessionStatus returns the current status of a session
func (p *MockAssemblyAIProvider) GetSessionStatus(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
	session, exists := p.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	result := &services.AudioProcessingResult{
		TranscriptionID: session.TranscriptID,
		Status:          session.Status,
		ProcessingMode:  session.Options.Mode,
		Message:         fmt.Sprintf("Session %s status: %s", sessionID, session.Status),
	}

	return result, nil
}

// AbortSession cancels an ongoing session
func (p *MockAssemblyAIProvider) AbortSession(ctx context.Context, sessionID string) error {
	session, exists := p.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.Status = entities.Failed
	delete(p.sessions, sessionID)

	log.Printf("Mock AssemblyAI session aborted: %s", sessionID)
	return nil
}

// IsSessionActive checks if a session is currently active
func (p *MockAssemblyAIProvider) IsSessionActive(ctx context.Context, sessionID string) bool {
	session, exists := p.sessions[sessionID]
	if !exists {
		return false
	}
	return session.Status == entities.Processing || session.Status == entities.Pending
}

// Helper method to concatenate audio chunks
func (p *MockAssemblyAIProvider) concatenateAudioChunks(chunks []services.AudioChunk) []byte {
	var totalSize int
	for _, chunk := range chunks {
		totalSize += len(chunk.Data)
	}

	result := make([]byte, 0, totalSize)
	for _, chunk := range chunks {
		result = append(result, chunk.Data...)
	}

	return result
}

// Helper method to generate mock transcript segments with speaker diarization
func (p *MockAssemblyAIProvider) generateMockTranscriptSegments(session *MockAssemblyAISession) []entities.TranscriptSegment {
	segments := make([]entities.TranscriptSegment, 0)

	// Generate realistic mock data based on the number of chunks
	numSegments := len(session.AudioChunks)
	if numSegments == 0 {
		numSegments = 1
	}

	mockTexts := []string{
		"Welcome everyone to today's meeting.",
		"Thank you for joining us. Let's get started with the agenda.",
		"I'd like to discuss the project timeline first.",
		"That sounds like a great idea. What are your thoughts?",
		"I agree with that approach. We should move forward.",
		"Let me share my screen to show the current progress.",
		"These results look promising. Should we proceed to the next phase?",
		"I think we need to consider the budget implications.",
		"Good point. Let's schedule a follow-up meeting to discuss this further.",
		"Thank you everyone for your time. Have a great day!",
	}

	speakers := []string{"Speaker A", "Speaker B", "Speaker C"}
	if session.Options.SpeakerDiarization {
		log.Printf("Mock: Generating diarized transcript for session %s", session.SessionID)
	}

	for i := 0; i < numSegments && i < len(mockTexts); i++ {
		speaker := "Speaker Unknown"
		if session.Options.SpeakerDiarization {
			speaker = speakers[i%len(speakers)]
		}

		startTime := float64(i) * 3.0          // 3 seconds per segment
		endTime := float64(i+1) * 3.0          // End of current segment
		confidence := 0.85 + float64(i%3)*0.05 // Varying confidence

		segment := entities.NewTranscriptSegment(
			session.SessionID,
			speaker,
			mockTexts[i],
			startTime,
			endTime,
			confidence,
			i+1,
		)
		segments = append(segments, segment)
	}

	log.Printf("Mock: Generated %d transcript segments for session %s", len(segments), session.SessionID)
	return segments
}
