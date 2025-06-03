package repositories

import (
	"context"
	"teammate/server/modules/transcription/domain/entities"
)

// TranscriptionRepository defines the interface for transcription persistence
type TranscriptionRepository interface {
	// Transcription CRUD operations
	Save(ctx context.Context, transcription *entities.Transcription) error
	FindByID(ctx context.Context, id string) (*entities.Transcription, error)
	FindByMeetingID(ctx context.Context, meetingID string) ([]*entities.Transcription, error)
	Update(ctx context.Context, transcription *entities.Transcription) error
	Delete(ctx context.Context, id string) error

	// Transcript segment operations
	SaveSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error
	FindSegmentsByTranscriptionID(ctx context.Context, transcriptionID string) ([]entities.TranscriptSegment, error)
	UpdateSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error
	UpdateSegment(ctx context.Context, segment *entities.TranscriptSegment) error

	// Query operations
	FindByStatus(ctx context.Context, status entities.TranscriptionStatus) ([]*entities.Transcription, error)
	FindByProvider(ctx context.Context, provider string) ([]*entities.Transcription, error)
	FindPendingTranscriptions(ctx context.Context) ([]*entities.Transcription, error)

	// Statistics and analytics
	GetTranscriptionStats(ctx context.Context, meetingID string) (*TranscriptionStats, error)
	GetSegmentCount(ctx context.Context, transcriptionID string) (int, error)
}

// TranscriptionStats provides analytics data for transcriptions
type TranscriptionStats struct {
	TotalDuration     float64 `json:"total_duration"`
	SegmentCount      int     `json:"segment_count"`
	SpeakerCount      int     `json:"speaker_count"`
	AverageConfidence float64 `json:"average_confidence"`
	WordCount         int     `json:"word_count"`
	ProcessingTime    float64 `json:"processing_time"`
}
