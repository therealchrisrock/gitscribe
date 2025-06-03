package commands

import (
	"context"
	"strings"
	"time"

	meetingRepos "teammate/server/modules/meeting/domain/repositories"
	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/modules/transcription/domain/services"
	"teammate/server/seedwork/domain"
	"teammate/server/seedwork/infrastructure/events"
)

// CompleteTranscriptionCommand represents a command to complete a transcription session
type CompleteTranscriptionCommand struct {
	TranscriptionID string                  `json:"transcription_id"`
	SessionID       string                  `json:"session_id"`
	MeetingID       string                  `json:"meeting_id"`
	BotSessionID    *string                 `json:"bot_session_id,omitempty"`
	Processor       services.AudioProcessor `json:"-"` // Not serialized
}

// CompleteTranscriptionResult represents the result of completing a transcription
type CompleteTranscriptionResult struct {
	TranscriptionID string                       `json:"transcription_id"`
	MeetingID       string                       `json:"meeting_id"`
	Status          entities.TranscriptionStatus `json:"status"`
	Segments        []entities.TranscriptSegment `json:"segments"`
	AudioFilePath   string                       `json:"audio_file_path"`
	ProcessingMode  services.ProcessingMode      `json:"processing_mode"`
	Message         string                       `json:"message"`
	CompletedAt     time.Time                    `json:"completed_at"`
}

// CompleteTranscriptionHandler handles the complete transcription command
type CompleteTranscriptionHandler struct {
	transcriptionRepo repositories.TranscriptionRepository
	meetingRepo       meetingRepos.MeetingRepository
	eventBus          events.EventBus
}

// NewCompleteTranscriptionHandler creates a new complete transcription handler
func NewCompleteTranscriptionHandler(
	transcriptionRepo repositories.TranscriptionRepository,
	meetingRepo meetingRepos.MeetingRepository,
	eventBus events.EventBus,
) *CompleteTranscriptionHandler {
	return &CompleteTranscriptionHandler{
		transcriptionRepo: transcriptionRepo,
		meetingRepo:       meetingRepo,
		eventBus:          eventBus,
	}
}

// Handle executes the complete transcription command
func (h *CompleteTranscriptionHandler) Handle(ctx context.Context, cmd CompleteTranscriptionCommand) (*CompleteTranscriptionResult, error) {
	// End the audio processing session
	result, err := cmd.Processor.EndSession(ctx, cmd.SessionID)
	if err != nil {
		// Mark transcription as failed
		if transcription, getErr := h.transcriptionRepo.FindByID(ctx, cmd.TranscriptionID); getErr == nil {
			transcription.FailTranscription()
			h.transcriptionRepo.Update(ctx, transcription)
		}

		// Update bot session status if exists
		if cmd.BotSessionID != nil {
			h.updateBotSessionStatus(ctx, *cmd.BotSessionID, meetingRepos.BotSessionStatusFailed)
		}

		return nil, domain.NewDomainError("END_SESSION_FAILED", "Failed to end audio session", err)
	}

	// Load transcription aggregate
	transcription, err := h.transcriptionRepo.FindByID(ctx, cmd.TranscriptionID)
	if err != nil {
		return nil, domain.NewDomainError("TRANSCRIPTION_NOT_FOUND", "Transcription not found", err)
	}

	// Apply business rules through aggregate methods
	content := h.segmentsToText(result.Segments)
	confidence := h.calculateAverageConfidence(result.Segments)
	transcription.CompleteTranscription(content, confidence, result.Segments)
	transcription.AudioFilePath = result.FirebaseURL

	// Persist transcription changes
	err = h.transcriptionRepo.Update(ctx, transcription)
	if err != nil {
		return nil, domain.NewDomainError("UPDATE_TRANSCRIPTION_FAILED", "Failed to update transcription", err)
	}

	// Save transcript segments
	if len(result.Segments) > 0 {
		err = h.transcriptionRepo.SaveSegments(ctx, cmd.TranscriptionID, result.Segments)
		if err != nil {
			return nil, domain.NewDomainError("SAVE_SEGMENTS_FAILED", "Failed to save transcript segments", err)
		}
	}

	// Update meeting status to completed
	if meeting, err := h.meetingRepo.FindMeetingByID(ctx, cmd.MeetingID); err == nil {
		meeting.Status = meetingRepos.MeetingStatusCompleted
		now := time.Now()
		meeting.EndTime = &now
		h.meetingRepo.UpdateMeeting(ctx, meeting)
	}

	// Update bot session status if exists
	if cmd.BotSessionID != nil {
		h.updateBotSessionStatus(ctx, *cmd.BotSessionID, meetingRepos.BotSessionStatusCompleted)
	}

	// Publish completion domain event
	event := &TranscriptionCompletedEvent{
		TranscriptionID: cmd.TranscriptionID,
		MeetingID:       cmd.MeetingID,
		SessionID:       cmd.SessionID,
		Status:          result.Status,
		SegmentCount:    len(result.Segments),
		ProcessingMode:  result.ProcessingMode,
		CompletedAt:     time.Now(),
	}
	h.eventBus.Publish("transcription.completed", event)

	return &CompleteTranscriptionResult{
		TranscriptionID: cmd.TranscriptionID,
		MeetingID:       cmd.MeetingID,
		Status:          result.Status,
		Segments:        result.Segments,
		AudioFilePath:   result.FirebaseURL,
		ProcessingMode:  result.ProcessingMode,
		Message:         result.Message,
		CompletedAt:     time.Now(),
	}, nil
}

// Helper methods - these should ideally be domain services
func (h *CompleteTranscriptionHandler) updateBotSessionStatus(ctx context.Context, sessionID, status string) error {
	session, err := h.meetingRepo.FindBotSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Status = status
	if status == meetingRepos.BotSessionStatusCompleted || status == meetingRepos.BotSessionStatusFailed {
		now := time.Now()
		session.LeftAt = &now
	}

	return h.meetingRepo.UpdateBotSession(ctx, session)
}

func (h *CompleteTranscriptionHandler) segmentsToText(segments []entities.TranscriptSegment) string {
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

func (h *CompleteTranscriptionHandler) calculateAverageConfidence(segments []entities.TranscriptSegment) float64 {
	if len(segments) == 0 {
		return 0.0
	}

	var total float64
	for _, segment := range segments {
		total += segment.Confidence
	}

	return total / float64(len(segments))
}

// TranscriptionCompletedEvent represents a domain event for transcription completed
type TranscriptionCompletedEvent struct {
	TranscriptionID string                       `json:"transcription_id"`
	MeetingID       string                       `json:"meeting_id"`
	SessionID       string                       `json:"session_id"`
	Status          entities.TranscriptionStatus `json:"status"`
	SegmentCount    int                          `json:"segment_count"`
	ProcessingMode  services.ProcessingMode      `json:"processing_mode"`
	CompletedAt     time.Time                    `json:"completed_at"`
}
