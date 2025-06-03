package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"

	"github.com/google/uuid"
)

// PostgresTranscriptionRepository implements TranscriptionRepository using PostgreSQL
type PostgresTranscriptionRepository struct {
	db *sql.DB
}

// NewPostgresTranscriptionRepository creates a new PostgreSQL transcription repository
func NewPostgresTranscriptionRepository(db *sql.DB) *PostgresTranscriptionRepository {
	return &PostgresTranscriptionRepository{db: db}
}

// Save stores a transcription in the database
func (r *PostgresTranscriptionRepository) Save(ctx context.Context, transcription *entities.Transcription) error {
	query := `
		INSERT INTO transcriptions (id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			content = EXCLUDED.content,
			confidence = EXCLUDED.confidence,
			updated_at = EXCLUDED.updated_at
	`

	if transcription.ID == "" {
		transcription.ID = uuid.New().String()
	}

	now := time.Now()
	if transcription.CreatedAt.IsZero() {
		transcription.CreatedAt = now
	}
	transcription.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		transcription.ID,
		transcription.MeetingID,
		transcription.AudioFilePath,
		string(transcription.Status),
		transcription.Content,
		transcription.Confidence,
		transcription.Provider,
		transcription.CreatedAt,
		transcription.UpdatedAt,
	)

	return err
}

// FindByID retrieves a transcription by its ID
func (r *PostgresTranscriptionRepository) FindByID(ctx context.Context, id string) (*entities.Transcription, error) {
	query := `
		SELECT id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at
		FROM transcriptions
		WHERE id = $1
	`

	var transcription entities.Transcription
	var status string
	var content sql.NullString
	var confidence sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transcription.ID,
		&transcription.MeetingID,
		&transcription.AudioFilePath,
		&status,
		&content,
		&confidence,
		&transcription.Provider,
		&transcription.CreatedAt,
		&transcription.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transcription not found: %s", id)
		}
		return nil, err
	}

	transcription.Status = entities.TranscriptionStatus(status)
	if content.Valid {
		transcription.Content = content.String
	}
	if confidence.Valid {
		transcription.Confidence = confidence.Float64
	}

	return &transcription, nil
}

// FindByMeetingID retrieves all transcriptions for a meeting
func (r *PostgresTranscriptionRepository) FindByMeetingID(ctx context.Context, meetingID string) ([]*entities.Transcription, error) {
	query := `
		SELECT id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at
		FROM transcriptions
		WHERE meeting_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transcriptions []*entities.Transcription

	for rows.Next() {
		var transcription entities.Transcription
		var status string
		var content sql.NullString
		var confidence sql.NullFloat64

		err := rows.Scan(
			&transcription.ID,
			&transcription.MeetingID,
			&transcription.AudioFilePath,
			&status,
			&content,
			&confidence,
			&transcription.Provider,
			&transcription.CreatedAt,
			&transcription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		transcription.Status = entities.TranscriptionStatus(status)
		if content.Valid {
			transcription.Content = content.String
		}
		if confidence.Valid {
			transcription.Confidence = confidence.Float64
		}

		transcriptions = append(transcriptions, &transcription)
	}

	return transcriptions, rows.Err()
}

// Update updates an existing transcription
func (r *PostgresTranscriptionRepository) Update(ctx context.Context, transcription *entities.Transcription) error {
	query := `
		UPDATE transcriptions 
		SET status = $2, content = $3, confidence = $4, updated_at = $5
		WHERE id = $1
	`

	transcription.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		transcription.ID,
		string(transcription.Status),
		transcription.Content,
		transcription.Confidence,
		transcription.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transcription not found: %s", transcription.ID)
	}

	return nil
}

// Delete removes a transcription from the database
func (r *PostgresTranscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transcriptions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transcription not found: %s", id)
	}

	return nil
}

// SaveSegments stores transcript segments for a transcription
func (r *PostgresTranscriptionRepository) SaveSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error {
	if len(segments) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing segments
	_, err = tx.ExecContext(ctx, "DELETE FROM transcript_segments WHERE transcription_id = $1", transcriptionID)
	if err != nil {
		return err
	}

	// Insert new segments
	query := `
		INSERT INTO transcript_segments (id, transcription_id, speaker, text, start_time, end_time, confidence, sequence_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for _, segment := range segments {
		if segment.ID == "" {
			segment.ID = uuid.New().String()
		}
		now := time.Now()
		if segment.CreatedAt.IsZero() {
			segment.CreatedAt = now
		}
		if segment.UpdatedAt.IsZero() {
			segment.UpdatedAt = now
		}

		_, err = tx.ExecContext(ctx, query,
			segment.ID,
			transcriptionID,
			segment.Speaker,
			segment.Text,
			segment.StartTime,
			segment.EndTime,
			segment.Confidence,
			segment.SequenceNumber,
			segment.CreatedAt,
			segment.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindSegmentsByTranscriptionID retrieves all segments for a transcription
func (r *PostgresTranscriptionRepository) FindSegmentsByTranscriptionID(ctx context.Context, transcriptionID string) ([]entities.TranscriptSegment, error) {
	query := `
		SELECT id, transcription_id, speaker, text, start_time, end_time, confidence, sequence_number, created_at, updated_at
		FROM transcript_segments
		WHERE transcription_id = $1
		ORDER BY sequence_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, transcriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []entities.TranscriptSegment

	for rows.Next() {
		var segment entities.TranscriptSegment
		var speaker sql.NullString
		var confidence sql.NullFloat64

		err := rows.Scan(
			&segment.ID,
			&segment.TranscriptionID,
			&speaker,
			&segment.Text,
			&segment.StartTime,
			&segment.EndTime,
			&confidence,
			&segment.SequenceNumber,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if speaker.Valid {
			segment.Speaker = speaker.String
		}
		if confidence.Valid {
			segment.Confidence = confidence.Float64
		}

		segments = append(segments, segment)
	}

	return segments, rows.Err()
}

// UpdateSegments updates transcript segments for a transcription
func (r *PostgresTranscriptionRepository) UpdateSegments(ctx context.Context, transcriptionID string, segments []entities.TranscriptSegment) error {
	return r.SaveSegments(ctx, transcriptionID, segments) // Use upsert behavior
}

// UpdateSegment updates a single transcript segment
func (r *PostgresTranscriptionRepository) UpdateSegment(ctx context.Context, segment *entities.TranscriptSegment) error {
	if segment.GetID() == "" {
		return fmt.Errorf("segment ID is required for update")
	}

	query := `
		UPDATE transcript_segments 
		SET speaker = $2, text = $3, start_time = $4, end_time = $5, confidence = $6, 
		    sequence_number = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		segment.GetID(),
		segment.Speaker,
		segment.Text,
		segment.StartTime,
		segment.EndTime,
		segment.Confidence,
		segment.SequenceNumber,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("segment not found: %s", segment.GetID())
	}

	// Update the segment's UpdatedAt field to reflect the database change
	segment.UpdatedAt = time.Now()

	return nil
}

// FindByStatus retrieves transcriptions by status
func (r *PostgresTranscriptionRepository) FindByStatus(ctx context.Context, status entities.TranscriptionStatus) ([]*entities.Transcription, error) {
	query := `
		SELECT id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at
		FROM transcriptions
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTranscriptions(rows)
}

// FindByProvider retrieves transcriptions by provider
func (r *PostgresTranscriptionRepository) FindByProvider(ctx context.Context, provider string) ([]*entities.Transcription, error) {
	query := `
		SELECT id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at
		FROM transcriptions
		WHERE provider = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTranscriptions(rows)
}

// FindPendingTranscriptions retrieves transcriptions that need processing
func (r *PostgresTranscriptionRepository) FindPendingTranscriptions(ctx context.Context) ([]*entities.Transcription, error) {
	query := `
		SELECT id, meeting_id, audio_file_path, status, content, confidence, provider, created_at, updated_at
		FROM transcriptions
		WHERE status IN ('pending', 'processing')
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTranscriptions(rows)
}

// GetTranscriptionStats retrieves analytics data for a meeting's transcriptions
func (r *PostgresTranscriptionRepository) GetTranscriptionStats(ctx context.Context, meetingID string) (*repositories.TranscriptionStats, error) {
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
		WHERE t.meeting_id = $1
		GROUP BY t.meeting_id
	`

	var stats repositories.TranscriptionStats
	err := r.db.QueryRowContext(ctx, query, meetingID).Scan(
		&stats.TotalDuration,
		&stats.SegmentCount,
		&stats.SpeakerCount,
		&stats.AverageConfidence,
		&stats.WordCount,
		&stats.ProcessingTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &repositories.TranscriptionStats{}, nil
		}
		return nil, err
	}

	return &stats, nil
}

// GetSegmentCount returns the number of segments for a transcription
func (r *PostgresTranscriptionRepository) GetSegmentCount(ctx context.Context, transcriptionID string) (int, error) {
	query := `SELECT COUNT(*) FROM transcript_segments WHERE transcription_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, transcriptionID).Scan(&count)
	return count, err
}

// Helper method to scan transcriptions from rows
func (r *PostgresTranscriptionRepository) scanTranscriptions(rows *sql.Rows) ([]*entities.Transcription, error) {
	var transcriptions []*entities.Transcription

	for rows.Next() {
		var transcription entities.Transcription
		var status string
		var content sql.NullString
		var confidence sql.NullFloat64

		err := rows.Scan(
			&transcription.ID,
			&transcription.MeetingID,
			&transcription.AudioFilePath,
			&status,
			&content,
			&confidence,
			&transcription.Provider,
			&transcription.CreatedAt,
			&transcription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		transcription.Status = entities.TranscriptionStatus(status)
		if content.Valid {
			transcription.Content = content.String
		}
		if confidence.Valid {
			transcription.Confidence = confidence.Float64
		}

		transcriptions = append(transcriptions, &transcription)
	}

	return transcriptions, rows.Err()
}
