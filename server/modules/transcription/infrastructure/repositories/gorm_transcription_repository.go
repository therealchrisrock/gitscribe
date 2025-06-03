package repositories

import (
	"context"
	"fmt"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
	"teammate/server/seedwork/infrastructure/database"

	"gorm.io/gorm"
)

// GormTranscriptionRepository implements TranscriptionRepository using GORM
type GormTranscriptionRepository struct {
	db *gorm.DB
}

// NewGormTranscriptionRepository creates a new GORM transcription repository
func NewGormTranscriptionRepository() *GormTranscriptionRepository {
	return &GormTranscriptionRepository{db: database.GetDB()}
}

// Save stores a transcription in the database
func (r *GormTranscriptionRepository) Save(ctx context.Context, transcription *entities.Transcription) error {
	return r.db.WithContext(ctx).Save(transcription).Error
}

// FindByID retrieves a transcription by its ID
func (r *GormTranscriptionRepository) FindByID(ctx context.Context, id string) (*entities.Transcription, error) {
	var transcription entities.Transcription
	err := r.db.WithContext(ctx).First(&transcription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transcription, nil
}

// FindByMeetingID retrieves all transcriptions for a meeting
func (r *GormTranscriptionRepository) FindByMeetingID(ctx context.Context, meetingID string) ([]*entities.Transcription, error) {
	var transcriptions []*entities.Transcription
	err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Order("created_at DESC").Find(&transcriptions).Error
	return transcriptions, err
}

// Update updates an existing transcription
func (r *GormTranscriptionRepository) Update(ctx context.Context, transcription *entities.Transcription) error {
	result := r.db.WithContext(ctx).Save(transcription)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("transcription not found: %s", transcription.ID)
	}
	return nil
}

// Delete removes a transcription from the database
func (r *GormTranscriptionRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entities.Transcription{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("transcription not found: %s", id)
	}
	return nil
}

// SaveSegments stores transcript segments
func (r *GormTranscriptionRepository) SaveSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error {
	// First, delete existing segments for this transcription
	err := r.db.WithContext(ctx).Where("transcription_id = ?", transcriptionID).Delete(&entities.TranscriptSegment{}).Error
	if err != nil {
		return err
	}

	// Insert new segments
	for i := range segments {
		segments[i].TranscriptionID = transcriptionID
	}

	return r.db.WithContext(ctx).Create(&segments).Error
}

// FindSegmentsByTranscriptionID retrieves segments for a transcription
func (r *GormTranscriptionRepository) FindSegmentsByTranscriptionID(ctx context.Context, transcriptionID string) ([]entities.TranscriptSegment, error) {
	var segments []entities.TranscriptSegment
	err := r.db.WithContext(ctx).Where("transcription_id = ?", transcriptionID).Order("start_time ASC").Find(&segments).Error
	return segments, err
}

// UpdateSegments updates transcript segments
func (r *GormTranscriptionRepository) UpdateSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error {
	return r.SaveSegments(ctx, transcriptionID, segments)
}

// UpdateSegment updates a single transcript segment
func (r *GormTranscriptionRepository) UpdateSegment(ctx context.Context, segment *entities.TranscriptSegment) error {
	if segment.GetID() == "" {
		return fmt.Errorf("segment ID is required for update")
	}

	result := r.db.WithContext(ctx).Save(segment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("segment not found: %s", segment.GetID())
	}
	return nil
}

// FindByStatus retrieves transcriptions by status
func (r *GormTranscriptionRepository) FindByStatus(ctx context.Context, status entities.TranscriptionStatus) ([]*entities.Transcription, error) {
	var transcriptions []*entities.Transcription
	err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&transcriptions).Error
	return transcriptions, err
}

// FindByProvider retrieves transcriptions by provider
func (r *GormTranscriptionRepository) FindByProvider(ctx context.Context, provider string) ([]*entities.Transcription, error) {
	var transcriptions []*entities.Transcription
	err := r.db.WithContext(ctx).Where("provider = ?", provider).Find(&transcriptions).Error
	return transcriptions, err
}

// FindPendingTranscriptions retrieves pending transcriptions
func (r *GormTranscriptionRepository) FindPendingTranscriptions(ctx context.Context) ([]*entities.Transcription, error) {
	return r.FindByStatus(ctx, entities.Pending)
}

// GetTranscriptionStats retrieves analytics for a meeting's transcriptions
func (r *GormTranscriptionRepository) GetTranscriptionStats(ctx context.Context, meetingID string) (*repositories.TranscriptionStats, error) {
	stats := &repositories.TranscriptionStats{}

	// Get total segments count and other analytics using raw query for simplicity
	query := `
		SELECT 
			COALESCE(MAX(ts.end_time), 0) as total_duration,
			COUNT(ts.id) as segment_count,
			COUNT(DISTINCT ts.speaker) as speaker_count,
			COALESCE(AVG(ts.confidence), 0) as average_confidence,
			COALESCE(SUM(array_length(string_to_array(ts.text, ' '), 1)), 0) as word_count,
			EXTRACT(EPOCH FROM (MAX(t.updated_at) - MIN(t.created_at))) as processing_time
		FROM transcriptions t
		LEFT JOIN transcript_segments ts ON t.id = ts.transcription_id
		WHERE t.meeting_id = ?
		GROUP BY t.meeting_id
	`

	err := r.db.WithContext(ctx).Raw(query, meetingID).Scan(stats).Error
	if err != nil {
		// Return empty stats if no data found
		return &repositories.TranscriptionStats{}, nil
	}

	return stats, nil
}

// GetSegmentCount returns the number of segments for a transcription
func (r *GormTranscriptionRepository) GetSegmentCount(ctx context.Context, transcriptionID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.TranscriptSegment{}).Where("transcription_id = ?", transcriptionID).Count(&count).Error
	return int(count), err
}
