package providers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/services"

	assemblyai "github.com/therealchrisrock/assemblyai-go"
)

// AssemblyAIProvider implements the AudioProcessor interface using AssemblyAI
type AssemblyAIProvider struct {
	client           *assemblyai.Client
	firebaseUploader FirebaseUploader
	sessions         map[string]*AssemblyAISession
}

// AssemblyAISession tracks an active AssemblyAI processing session
type AssemblyAISession struct {
	SessionID    string
	TranscriptID string
	Metadata     services.AudioStreamMetadata
	Options      services.AudioProcessingOptions
	AudioChunks  []services.AudioChunk
	Status       entities.TranscriptionStatus
	CreatedAt    time.Time
	FirebaseURL  string
}

// FirebaseUploader interface for uploading audio files to Firebase Storage
type FirebaseUploader interface {
	UploadAudio(ctx context.Context, audioData []byte, meetingID, sessionID string) (string, error)
	UploadAudioStream(ctx context.Context, reader io.Reader, meetingID, sessionID string) (string, error)
}

// NewAssemblyAIProvider creates a new AssemblyAI provider
func NewAssemblyAIProvider(apiKey string, firebaseUploader FirebaseUploader) *AssemblyAIProvider {
	client := assemblyai.NewClient(apiKey)
	return &AssemblyAIProvider{
		client:           client,
		firebaseUploader: firebaseUploader,
		sessions:         make(map[string]*AssemblyAISession),
	}
}

// GetSupportedModes returns the processing modes supported by AssemblyAI
func (p *AssemblyAIProvider) GetSupportedModes() []services.ProcessingMode {
	return []services.ProcessingMode{
		services.RealTimeMode,
		services.BatchMode,
	}
}

// StartSession initializes a new AssemblyAI processing session
func (p *AssemblyAIProvider) StartSession(ctx context.Context, metadata services.AudioStreamMetadata, options services.AudioProcessingOptions) (string, error) {
	sessionID := metadata.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("assemblyai_%d", time.Now().UnixNano())
	}

	session := &AssemblyAISession{
		SessionID:   sessionID,
		Metadata:    metadata,
		Options:     options,
		AudioChunks: make([]services.AudioChunk, 0),
		Status:      entities.Pending,
		CreatedAt:   time.Now(),
	}

	p.sessions[sessionID] = session

	log.Printf("AssemblyAI session started: %s (mode: %s, diarization: %t)",
		sessionID, options.Mode, options.SpeakerDiarization)

	return sessionID, nil
}

// ProcessChunk processes an individual audio chunk
func (p *AssemblyAIProvider) ProcessChunk(ctx context.Context, sessionID string, chunk services.AudioChunk) error {
	session, exists := p.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Add chunk to session
	session.AudioChunks = append(session.AudioChunks, chunk)
	session.Status = entities.Processing

	log.Printf("AssemblyAI: Added chunk #%d to session %s (%d bytes)",
		len(session.AudioChunks), sessionID, len(chunk.Data))

	return nil
}

// EndSession finalizes the session and processes the complete audio
func (p *AssemblyAIProvider) EndSession(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
	session, exists := p.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	// Concatenate all audio chunks
	audioData := p.concatenateAudioChunks(session.AudioChunks)

	// Upload to Firebase Storage
	firebaseURL, err := p.firebaseUploader.UploadAudio(ctx, audioData, session.Metadata.MeetingID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio to Firebase: %w", err)
	}
	session.FirebaseURL = firebaseURL
	log.Printf("AssemblyAI: Audio uploaded to Firebase for session %s: %s", sessionID, firebaseURL)

	// Upload to AssemblyAI
	uploadResp, err := p.client.UploadFile(ctx, io.NopCloser(bytes.NewReader(audioData)))
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio to AssemblyAI: %w", err)
	}

	// Create transcription request with diarization
	request := p.buildTranscriptRequest(uploadResp.UploadURL, session.Options)

	// Submit for transcription
	transcript, err := p.client.CreateTranscript(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create AssemblyAI transcript: %w", err)
	}

	session.TranscriptID = transcript.ID

	// Poll for completion
	transcript, err = p.pollForCompletion(ctx, transcript.ID)
	if err != nil {
		session.Status = entities.Failed
		return nil, fmt.Errorf("failed to get transcript result: %w", err)
	}

	// Process results
	segments := p.convertToTranscriptSegments(transcript, sessionID)
	session.Status = entities.Completed

	result := &services.AudioProcessingResult{
		TranscriptionID: transcript.ID,
		Status:          session.Status,
		Segments:        segments,
		ProcessingMode:  session.Options.Mode,
		Message:         fmt.Sprintf("Transcription completed with %d segments", len(segments)),
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
func (p *AssemblyAIProvider) GetSessionStatus(ctx context.Context, sessionID string) (*services.AudioProcessingResult, error) {
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
func (p *AssemblyAIProvider) AbortSession(ctx context.Context, sessionID string) error {
	session, exists := p.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.Status = entities.Failed
	delete(p.sessions, sessionID)

	log.Printf("AssemblyAI session aborted: %s", sessionID)
	return nil
}

// IsSessionActive checks if a session is currently active
func (p *AssemblyAIProvider) IsSessionActive(ctx context.Context, sessionID string) bool {
	session, exists := p.sessions[sessionID]
	if !exists {
		return false
	}
	return session.Status == entities.Processing || session.Status == entities.Pending
}

// Helper method to concatenate audio chunks
func (p *AssemblyAIProvider) concatenateAudioChunks(chunks []services.AudioChunk) []byte {
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

// Helper method to build AssemblyAI transcript request
func (p *AssemblyAIProvider) buildTranscriptRequest(audioURL string, options services.AudioProcessingOptions) *assemblyai.TranscriptRequest {
	request := &assemblyai.TranscriptRequest{
		AudioURL:          audioURL,
		SpeakerLabels:     assemblyai.Bool(options.SpeakerDiarization),
		Punctuate:         assemblyai.Bool(true),
		FormatText:        assemblyai.Bool(true),
		SentimentAnalysis: assemblyai.Bool(true),
		EntityDetection:   assemblyai.Bool(true),
		AutoHighlights:    assemblyai.Bool(true),
		Disfluencies:      assemblyai.Bool(false),
	}

	if options.Language != "" {
		request.LanguageCode = &options.Language
	}

	if options.ConfidenceThreshold > 0 {
		threshold := options.ConfidenceThreshold
		request.SpeechThreshold = &threshold
	}

	return request
}

// Helper method to poll for transcript completion
func (p *AssemblyAIProvider) pollForCompletion(ctx context.Context, transcriptID string) (*assemblyai.Transcript, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(30 * time.Minute) // 30 minute timeout
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout.C:
			return nil, fmt.Errorf("transcript polling timeout for ID: %s", transcriptID)
		case <-ticker.C:
			transcript, err := p.client.GetTranscript(ctx, transcriptID)
			if err != nil {
				log.Printf("Error polling transcript %s: %v", transcriptID, err)
				continue
			}

			switch transcript.Status {
			case assemblyai.StatusCompleted:
				return transcript, nil
			case assemblyai.StatusError:
				errorMsg := "unknown error"
				if transcript.Error != nil {
					errorMsg = *transcript.Error
				}
				return nil, fmt.Errorf("transcript failed: %s", errorMsg)
			case assemblyai.StatusQueued, assemblyai.StatusProcessing:
				log.Printf("Transcript %s still processing...", transcriptID)
				continue
			default:
				return nil, fmt.Errorf("unexpected transcript status: %s", transcript.Status)
			}
		}
	}
}

// Helper method to convert AssemblyAI transcript to domain segments
func (p *AssemblyAIProvider) convertToTranscriptSegments(transcript *assemblyai.Transcript, sessionID string) []entities.TranscriptSegment {
	segments := make([]entities.TranscriptSegment, 0)

	// Use utterances for speaker-diarized content
	if len(transcript.Utterances) > 0 {
		for i, utterance := range transcript.Utterances {
			// Skip utterances with empty or whitespace-only text
			if strings.TrimSpace(utterance.Text) == "" {
				log.Printf("Skipping empty utterance from speaker %s at %d-%d ms", utterance.Speaker, utterance.Start, utterance.End)
				continue
			}

			// Normalize speaker labels consistently
			speaker := p.normalizeSpeakerLabel(utterance.Speaker)

			segment := entities.NewTranscriptSegment(
				sessionID,
				speaker,
				utterance.Text,
				float64(utterance.Start)/1000.0, // Convert milliseconds to seconds
				float64(utterance.End)/1000.0,
				utterance.Confidence,
				i+1,
			)
			segments = append(segments, segment)
		}
	} else if transcript.Text != nil {
		// Fallback to full text if no utterances
		segment := entities.NewTranscriptSegment(
			sessionID,
			"Speaker Unknown",
			*transcript.Text,
			0.0,
			float64(*transcript.AudioDuration),
			*transcript.Confidence,
			1,
		)
		segments = append(segments, segment)
	}

	return segments
}

// Helper method to normalize speaker labels from AssemblyAI
func (p *AssemblyAIProvider) normalizeSpeakerLabel(speaker string) string {
	speaker = strings.TrimSpace(speaker)

	// Only normalize the truly "unknown" cases - preserve speaker differentiation
	if speaker == "" || speaker == "speaker_unknown" || speaker == "unknown" {
		return "Speaker Unknown"
	}

	// Keep all other labels as-is to preserve speaker differentiation
	// This includes: "Speaker A", "Speaker B", "speaker_0", "speaker_1", etc.
	return speaker
}
