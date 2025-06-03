package commands

import (
	"context"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/modules/transcription/domain/services"
	"teammate/server/seedwork/domain"
	"teammate/server/seedwork/infrastructure/events"
)

// ProcessAudioChunkCommand represents a command to process an audio chunk
type ProcessAudioChunkCommand struct {
	TranscriptionID string                  `json:"transcription_id"`
	SessionID       string                  `json:"session_id"`
	AudioChunk      services.AudioChunk     `json:"audio_chunk"`
	Processor       services.AudioProcessor `json:"-"` // Not serialized
}

// ProcessAudioChunkResult represents the result of processing an audio chunk
type ProcessAudioChunkResult struct {
	TranscriptionID string                       `json:"transcription_id"`
	ChunkProcessed  time.Time                    `json:"chunk_processed"`
	Status          entities.TranscriptionStatus `json:"status"`
}

// ProcessAudioChunkHandler handles the process audio chunk command
type ProcessAudioChunkHandler struct {
	transcriptionRepo repositories.TranscriptionRepository
	eventBus          events.EventBus
}

// NewProcessAudioChunkHandler creates a new process audio chunk handler
func NewProcessAudioChunkHandler(
	transcriptionRepo repositories.TranscriptionRepository,
	eventBus events.EventBus,
) *ProcessAudioChunkHandler {
	return &ProcessAudioChunkHandler{
		transcriptionRepo: transcriptionRepo,
		eventBus:          eventBus,
	}
}

// Handle executes the process audio chunk command
func (h *ProcessAudioChunkHandler) Handle(ctx context.Context, cmd ProcessAudioChunkCommand) (*ProcessAudioChunkResult, error) {
	// Process the chunk with the audio processor
	err := cmd.Processor.ProcessChunk(ctx, cmd.SessionID, cmd.AudioChunk)
	if err != nil {
		return nil, domain.NewDomainError("PROCESS_CHUNK_FAILED", "Failed to process audio chunk", err)
	}

	// Load transcription aggregate
	transcription, err := h.transcriptionRepo.FindByID(ctx, cmd.TranscriptionID)
	if err != nil {
		return nil, domain.NewDomainError("TRANSCRIPTION_NOT_FOUND", "Transcription not found", err)
	}

	// Apply business rules if status needs to change
	var statusChanged bool
	if transcription.Status == entities.Pending {
		transcription.StartProcessing()
		statusChanged = true

		// Persist status change
		err = h.transcriptionRepo.Update(ctx, transcription)
		if err != nil {
			return nil, domain.NewDomainError("UPDATE_TRANSCRIPTION_FAILED", "Failed to update transcription status", err)
		}
	}

	// Publish processing event if status changed
	if statusChanged {
		event := &TranscriptionProcessingEvent{
			TranscriptionID: cmd.TranscriptionID,
			SessionID:       cmd.SessionID,
			ChunkProcessed:  time.Now(),
		}
		h.eventBus.Publish("transcription.processing", event)
	}

	return &ProcessAudioChunkResult{
		TranscriptionID: cmd.TranscriptionID,
		ChunkProcessed:  time.Now(),
		Status:          transcription.Status,
	}, nil
}

// TranscriptionProcessingEvent represents a domain event for transcription processing
type TranscriptionProcessingEvent struct {
	TranscriptionID string    `json:"transcription_id"`
	SessionID       string    `json:"session_id"`
	ChunkProcessed  time.Time `json:"chunk_processed"`
}
