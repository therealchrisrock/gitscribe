package queries

import (
	"context"

	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/seedwork/domain"
)

// GetTranscriptionStatsQuery represents a query to get transcription statistics
type GetTranscriptionStatsQuery struct {
	MeetingID string `json:"meeting_id"`
}

// GetTranscriptionStatsResult represents the result of the stats query
type GetTranscriptionStatsResult struct {
	Stats *repositories.TranscriptionStats `json:"stats"`
}

// GetTranscriptionStatsHandler handles the get transcription stats query
type GetTranscriptionStatsHandler struct {
	transcriptionRepo repositories.TranscriptionRepository
}

// NewGetTranscriptionStatsHandler creates a new get transcription stats handler
func NewGetTranscriptionStatsHandler(
	transcriptionRepo repositories.TranscriptionRepository,
) *GetTranscriptionStatsHandler {
	return &GetTranscriptionStatsHandler{
		transcriptionRepo: transcriptionRepo,
	}
}

// Handle executes the get transcription stats query
func (h *GetTranscriptionStatsHandler) Handle(ctx context.Context, query GetTranscriptionStatsQuery) (*GetTranscriptionStatsResult, error) {
	// Get stats for the meeting
	stats, err := h.transcriptionRepo.GetTranscriptionStats(ctx, query.MeetingID)
	if err != nil {
		return nil, domain.NewDomainError("GET_STATS_FAILED", "Failed to get transcription stats", err)
	}

	return &GetTranscriptionStatsResult{
		Stats: stats,
	}, nil
}
