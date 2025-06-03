package services

import (
	"context"
	"strings"
	"time"

	meetingRepos "teammate/server/modules/meeting/domain/repositories"
	"teammate/server/modules/transcription/application/commands"
	"teammate/server/modules/transcription/application/queries"
	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/modules/transcription/domain/services"
	"teammate/server/seedwork/infrastructure/events"
)

// EnhancedTranscriptionService provides advanced transcription functionality using DDD commands
type EnhancedTranscriptionService struct {
	// Command handlers
	startHandler    *commands.StartTranscriptionHandler
	processHandler  *commands.ProcessAudioChunkHandler
	completeHandler *commands.CompleteTranscriptionHandler

	// Query handlers
	historyHandler *queries.GetTranscriptionHistoryHandler
	statsHandler   *queries.GetTranscriptionStatsHandler

	// Dependencies for non-command operations
	transcriptionRepo repositories.TranscriptionRepository
	meetingRepo       meetingRepos.MeetingRepository
	audioFactory      services.AudioProcessorFactory // Use domain interface
	concreteFactory   *AudioProcessorFactory         // Keep concrete for compatibility

	// TODO: Remove session management from application service - violates DDD
	// This should be handled by a proper SessionRepository or moved to infrastructure
	activeSessions map[string]*TranscriptionSession
}

// NewEnhancedTranscriptionService creates a new enhanced transcription service
func NewEnhancedTranscriptionService(
	transcriptionRepo repositories.TranscriptionRepository,
	meetingRepo meetingRepos.MeetingRepository,
	audioFactory *AudioProcessorFactory,
	eventBus events.EventBus,
) *EnhancedTranscriptionService {
	return &EnhancedTranscriptionService{
		startHandler:      commands.NewStartTranscriptionHandler(transcriptionRepo, meetingRepo, audioFactory, eventBus),
		processHandler:    commands.NewProcessAudioChunkHandler(transcriptionRepo, eventBus),
		completeHandler:   commands.NewCompleteTranscriptionHandler(transcriptionRepo, meetingRepo, eventBus),
		historyHandler:    queries.NewGetTranscriptionHistoryHandler(transcriptionRepo),
		statsHandler:      queries.NewGetTranscriptionStatsHandler(transcriptionRepo),
		transcriptionRepo: transcriptionRepo,
		meetingRepo:       meetingRepo,
		audioFactory:      audioFactory,
		concreteFactory:   audioFactory,
		activeSessions:    make(map[string]*TranscriptionSession),
	}
}

// StartTranscriptionSession starts a new transcription session using command pattern
func (s *EnhancedTranscriptionService) StartTranscriptionSession(ctx context.Context, req *StartTranscriptionRequest) (*TranscriptionSession, error) {
	// Create and execute command
	cmd := commands.StartTranscriptionCommand{
		MeetingID:           req.MeetingID,
		AudioStreamMetadata: req.Metadata,
		ProcessingOptions:   req.Options,
		CreateBotSession:    req.CreateBotSession,
	}

	result, err := s.startHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Create audio processor for session management using concrete factory
	processor, err := s.concreteFactory.CreateProcessor(req.Options.Mode, req.Options)
	if err != nil {
		return nil, err
	}

	// Create session object for compatibility
	session := &TranscriptionSession{
		ID:           result.TranscriptionID,
		SessionID:    result.SessionID,
		MeetingID:    result.MeetingID,
		Status:       entities.Processing,
		Processor:    processor,
		BotSessionID: result.BotSessionID,
		StartedAt:    result.StartedAt,
		Options:      req.Options,
	}

	// Store active session
	s.activeSessions[session.ID] = session

	return session, nil
}

// ProcessAudioChunk processes an audio chunk using command pattern
func (s *EnhancedTranscriptionService) ProcessAudioChunk(ctx context.Context, session *TranscriptionSession, chunk services.AudioChunk) error {
	// Create and execute command
	cmd := commands.ProcessAudioChunkCommand{
		TranscriptionID: session.ID,
		SessionID:       session.SessionID,
		AudioChunk:      chunk,
		Processor:       session.Processor,
	}

	result, err := s.processHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	// Update session status
	session.Status = result.Status

	return nil
}

// EndTranscriptionSession ends a transcription session using command pattern
func (s *EnhancedTranscriptionService) EndTranscriptionSession(ctx context.Context, session *TranscriptionSession) (*TranscriptionResult, error) {
	// Create and execute command
	cmd := commands.CompleteTranscriptionCommand{
		TranscriptionID: session.ID,
		SessionID:       session.SessionID,
		MeetingID:       session.MeetingID,
		BotSessionID:    session.BotSessionID,
		Processor:       session.Processor,
	}

	result, err := s.completeHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// TODO: Remove this session management - violates DDD aggregate boundaries
	// Session lifecycle should be managed by domain aggregates or infrastructure
	delete(s.activeSessions, session.ID)

	// Convert to expected result format
	return &TranscriptionResult{
		TranscriptionID: result.TranscriptionID,
		MeetingID:       result.MeetingID,
		Status:          result.Status,
		Segments:        result.Segments,
		AudioFilePath:   result.AudioFilePath,
		ProcessingMode:  result.ProcessingMode,
		Message:         result.Message,
		CompletedAt:     result.CompletedAt,
	}, nil
}

// GetTranscriptionHistory retrieves transcription history for a meeting
func (s *EnhancedTranscriptionService) GetTranscriptionHistory(ctx context.Context, meetingID string) ([]*TranscriptionHistoryItem, error) {
	// Create and execute query
	query := queries.GetTranscriptionHistoryQuery{
		MeetingID: meetingID,
	}

	result, err := s.historyHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert to legacy format for compatibility
	items := make([]*TranscriptionHistoryItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = &TranscriptionHistoryItem{
			ID:            item.ID,
			MeetingID:     item.MeetingID,
			Status:        item.Status,
			Provider:      item.Provider,
			Content:       item.Content,
			Confidence:    item.Confidence,
			AudioFilePath: item.AudioFilePath,
			Segments:      item.Segments,
			Stats:         item.Stats,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
	}

	return items, nil
}

// GetTranscriptionStats retrieves analytics for a meeting's transcriptions
func (s *EnhancedTranscriptionService) GetTranscriptionStats(ctx context.Context, meetingID string) (*repositories.TranscriptionStats, error) {
	// Create and execute query
	query := queries.GetTranscriptionStatsQuery{
		MeetingID: meetingID,
	}

	result, err := s.statsHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	return result.Stats, nil
}

// Request/Response Types

type StartTranscriptionRequest struct {
	MeetingID        string                          `json:"meeting_id"`
	Metadata         services.AudioStreamMetadata    `json:"metadata"`
	Options          services.AudioProcessingOptions `json:"options"`
	CreateBotSession bool                            `json:"create_bot_session,omitempty"` // Optional flag for bot session creation
}

type TranscriptionSession struct {
	ID           string                          `json:"id"`
	SessionID    string                          `json:"session_id"`
	MeetingID    string                          `json:"meeting_id"`
	Status       entities.TranscriptionStatus    `json:"status"`
	Processor    services.AudioProcessor         `json:"-"`
	BotSessionID *string                         `json:"bot_session_id,omitempty"` // Made optional
	StartedAt    time.Time                       `json:"started_at"`
	Options      services.AudioProcessingOptions `json:"options"`
}

type TranscriptionResult struct {
	TranscriptionID string                       `json:"transcription_id"`
	MeetingID       string                       `json:"meeting_id"`
	Status          entities.TranscriptionStatus `json:"status"`
	Segments        []entities.TranscriptSegment `json:"segments"`
	AudioFilePath   string                       `json:"audio_file_path"`
	ProcessingMode  services.ProcessingMode      `json:"processing_mode"`
	Message         string                       `json:"message"`
	CompletedAt     time.Time                    `json:"completed_at"`
}

type TranscriptionHistoryItem struct {
	ID            string                           `json:"id"`
	MeetingID     string                           `json:"meeting_id"`
	Status        entities.TranscriptionStatus     `json:"status"`
	Provider      string                           `json:"provider"`
	Content       string                           `json:"content"`
	Confidence    float64                          `json:"confidence"`
	AudioFilePath string                           `json:"audio_file_path"`
	Segments      []entities.TranscriptSegment     `json:"segments"`
	Stats         *repositories.TranscriptionStats `json:"stats"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     time.Time                        `json:"updated_at"`
}

// Event Types for real-time updates

type TranscriptionSessionStartedEvent struct {
	TranscriptionID string    `json:"transcription_id"`
	MeetingID       string    `json:"meeting_id"`
	SessionID       string    `json:"session_id"`
	Provider        string    `json:"provider"`
	StartedAt       time.Time `json:"started_at"`
}

type TranscriptionProcessingEvent struct {
	TranscriptionID string    `json:"transcription_id"`
	MeetingID       string    `json:"meeting_id"`
	SessionID       string    `json:"session_id"`
	ChunkProcessed  time.Time `json:"chunk_processed"`
}

type TranscriptionCompletedEvent struct {
	TranscriptionID string                       `json:"transcription_id"`
	MeetingID       string                       `json:"meeting_id"`
	SessionID       string                       `json:"session_id"`
	Status          entities.TranscriptionStatus `json:"status"`
	SegmentCount    int                          `json:"segment_count"`
	ProcessingMode  services.ProcessingMode      `json:"processing_mode"`
	CompletedAt     time.Time                    `json:"completed_at"`
}

// TODO: Move these helper methods to domain aggregates to fix DDD violations
// These should be methods on the Transcription aggregate root

func (s *EnhancedTranscriptionService) segmentsToText(segments []entities.TranscriptSegment) string {
	// TODO: This should be transcription.GenerateContentFromSegments()
	var text string
	for _, segment := range segments {
		// Skip segments with empty text to avoid creating malformed content
		if strings.TrimSpace(segment.Text) == "" {
			continue
		}

		// Handle speaker identification properly - only normalize truly unknown speakers
		if segment.Speaker != "" && segment.Speaker != "speaker_unknown" {
			text += segment.Speaker + ": " + segment.Text + "\n"
		} else {
			text += segment.Text + "\n"
		}
	}
	return text
}

func (s *EnhancedTranscriptionService) calculateAverageConfidence(segments []entities.TranscriptSegment) float64 {
	// TODO: This should be transcription.CalculateConfidence()
	if len(segments) == 0 {
		return 0.0
	}

	var total float64
	for _, segment := range segments {
		total += segment.Confidence
	}

	return total / float64(len(segments))
}
