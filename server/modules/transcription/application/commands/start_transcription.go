package commands

import (
	"context"
	"fmt"
	"time"

	meetingRepos "teammate/server/modules/meeting/domain/repositories"
	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/modules/transcription/domain/services"
	"teammate/server/seedwork/domain"
	"teammate/server/seedwork/infrastructure/events"
)

// StartTranscriptionCommand represents a command to start a new transcription session
type StartTranscriptionCommand struct {
	MeetingID           string                          `json:"meeting_id"`
	AudioStreamMetadata services.AudioStreamMetadata    `json:"metadata"`
	ProcessingOptions   services.AudioProcessingOptions `json:"options"`
	CreateBotSession    bool                            `json:"create_bot_session"`
}

// StartTranscriptionResult represents the result of starting a transcription
type StartTranscriptionResult struct {
	TranscriptionID string    `json:"transcription_id"`
	SessionID       string    `json:"session_id"`
	MeetingID       string    `json:"meeting_id"`
	BotSessionID    *string   `json:"bot_session_id,omitempty"`
	StartedAt       time.Time `json:"started_at"`
}

// StartTranscriptionHandler handles the start transcription command
type StartTranscriptionHandler struct {
	transcriptionRepo repositories.TranscriptionRepository
	meetingRepo       meetingRepos.MeetingRepository
	audioFactory      services.AudioProcessorFactory
	eventBus          events.EventBus
}

// NewStartTranscriptionHandler creates a new start transcription handler
func NewStartTranscriptionHandler(
	transcriptionRepo repositories.TranscriptionRepository,
	meetingRepo meetingRepos.MeetingRepository,
	audioFactory services.AudioProcessorFactory,
	eventBus events.EventBus,
) *StartTranscriptionHandler {
	return &StartTranscriptionHandler{
		transcriptionRepo: transcriptionRepo,
		meetingRepo:       meetingRepo,
		audioFactory:      audioFactory,
		eventBus:          eventBus,
	}
}

// Handle executes the start transcription command
func (h *StartTranscriptionHandler) Handle(ctx context.Context, cmd StartTranscriptionCommand) (*StartTranscriptionResult, error) {
	// Validate the meeting exists
	meeting, err := h.meetingRepo.FindMeetingByID(ctx, cmd.MeetingID)
	if err != nil {
		return nil, domain.NewDomainError("MEETING_NOT_FOUND", "Meeting not found", err)
	}

	// Create transcription aggregate using domain factory method
	transcription := entities.NewTranscription(cmd.MeetingID, "", cmd.ProcessingOptions.Provider)

	// Apply business rules through aggregate methods
	transcription.StartProcessing()

	// Save transcription to database
	err = h.transcriptionRepo.Save(ctx, &transcription)
	if err != nil {
		return nil, domain.NewDomainError("SAVE_TRANSCRIPTION_FAILED", "Failed to save transcription", err)
	}

	// Create audio processor
	processor, err := h.audioFactory.CreateProcessor(cmd.ProcessingOptions.Mode, cmd.ProcessingOptions)
	if err != nil {
		// Mark transcription as failed
		transcription.FailTranscription()
		h.transcriptionRepo.Update(ctx, &transcription)
		return nil, domain.NewDomainError("CREATE_PROCESSOR_FAILED", "Failed to create audio processor", err)
	}

	// Start audio processing session
	sessionID, err := processor.StartSession(ctx, cmd.AudioStreamMetadata, cmd.ProcessingOptions)
	if err != nil {
		// Mark transcription as failed
		transcription.FailTranscription()
		h.transcriptionRepo.Update(ctx, &transcription)
		return nil, domain.NewDomainError("START_SESSION_FAILED", "Failed to start audio session", err)
	}

	var botSessionID *string

	// Create bot session if requested
	if cmd.CreateBotSession {
		botSession := &meetingRepos.BotSession{
			MeetingID: cmd.MeetingID,
			SessionID: sessionID,
			Status:    meetingRepos.BotSessionStatusActive,
			JoinedAt:  time.Now(),
			Metadata: map[string]interface{}{
				"transcription_id": transcription.GetID(),
				"provider":         cmd.ProcessingOptions.Provider,
				"mode":             string(cmd.ProcessingOptions.Mode),
				"diarization":      cmd.ProcessingOptions.SpeakerDiarization,
			},
		}

		err = h.meetingRepo.SaveBotSession(ctx, botSession)
		if err != nil {
			return nil, domain.NewDomainError("SAVE_BOT_SESSION_FAILED", "Failed to save bot session", err)
		}
		botSessionID = &botSession.ID
	}

	// Update meeting status to in-progress
	meeting.Status = meetingRepos.MeetingStatusInProgress
	err = h.meetingRepo.UpdateMeeting(ctx, meeting)
	if err != nil {
		// Log warning but don't fail the command
		fmt.Printf("Warning: failed to update meeting status: %v\n", err)
	}

	// Publish domain event
	event := &TranscriptionStartedEvent{
		TranscriptionID: transcription.GetID(),
		MeetingID:       cmd.MeetingID,
		SessionID:       sessionID,
		Provider:        cmd.ProcessingOptions.Provider,
		StartedAt:       time.Now(),
	}
	h.eventBus.Publish("transcription.started", event)

	return &StartTranscriptionResult{
		TranscriptionID: transcription.GetID(),
		SessionID:       sessionID,
		MeetingID:       cmd.MeetingID,
		BotSessionID:    botSessionID,
		StartedAt:       time.Now(),
	}, nil
}

// TranscriptionStartedEvent represents a domain event for transcription started
type TranscriptionStartedEvent struct {
	TranscriptionID string    `json:"transcription_id"`
	MeetingID       string    `json:"meeting_id"`
	SessionID       string    `json:"session_id"`
	Provider        string    `json:"provider"`
	StartedAt       time.Time `json:"started_at"`
}
