package queries

import (
	"context"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/seedwork/domain"
)

// GetTranscriptionHistoryQuery represents a query to get transcription history
type GetTranscriptionHistoryQuery struct {
	MeetingID string `json:"meeting_id"`
}

// TranscriptionHistoryItem represents a transcription history item
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

// GetTranscriptionHistoryResult represents the result of the query
type GetTranscriptionHistoryResult struct {
	Items []TranscriptionHistoryItem `json:"items"`
	Total int                        `json:"total"`
}

// GetTranscriptionHistoryHandler handles the get transcription history query
type GetTranscriptionHistoryHandler struct {
	transcriptionRepo repositories.TranscriptionRepository
}

// NewGetTranscriptionHistoryHandler creates a new get transcription history handler
func NewGetTranscriptionHistoryHandler(
	transcriptionRepo repositories.TranscriptionRepository,
) *GetTranscriptionHistoryHandler {
	return &GetTranscriptionHistoryHandler{
		transcriptionRepo: transcriptionRepo,
	}
}

// Handle executes the get transcription history query
func (h *GetTranscriptionHistoryHandler) Handle(ctx context.Context, query GetTranscriptionHistoryQuery) (*GetTranscriptionHistoryResult, error) {
	// Find all transcriptions for the meeting
	transcriptions, err := h.transcriptionRepo.FindByMeetingID(ctx, query.MeetingID)
	if err != nil {
		return nil, domain.NewDomainError("FIND_TRANSCRIPTIONS_FAILED", "Failed to find transcriptions", err)
	}

	items := make([]TranscriptionHistoryItem, 0, len(transcriptions))

	for _, transcription := range transcriptions {
		// Get segments for this transcription
		segments, err := h.transcriptionRepo.FindSegmentsByTranscriptionID(ctx, transcription.GetID())
		if err != nil {
			// Log error but continue processing other transcriptions
			segments = []entities.TranscriptSegment{}
		}

		// Get stats for this transcription
		stats, err := h.transcriptionRepo.GetTranscriptionStats(ctx, transcription.GetID())
		if err != nil {
			// Log error but continue processing
			stats = nil
		}

		item := TranscriptionHistoryItem{
			ID:            transcription.GetID(),
			MeetingID:     transcription.MeetingID,
			Status:        transcription.Status,
			Provider:      transcription.Provider,
			Content:       transcription.Content,
			Confidence:    transcription.Confidence,
			AudioFilePath: transcription.AudioFilePath,
			Segments:      segments,
			Stats:         stats,
			CreatedAt:     transcription.CreatedAt,
			UpdatedAt:     transcription.UpdatedAt,
		}

		items = append(items, item)
	}

	return &GetTranscriptionHistoryResult{
		Items: items,
		Total: len(items),
	}, nil
}
