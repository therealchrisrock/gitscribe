package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"teammate/server/modules/meeting/domain/repositories"

	"github.com/google/uuid"
)

// PostgresMeetingRepository implements MeetingRepository using PostgreSQL
type PostgresMeetingRepository struct {
	db *sql.DB
}

// NewPostgresMeetingRepository creates a new PostgreSQL meeting repository
func NewPostgresMeetingRepository(db *sql.DB) *PostgresMeetingRepository {
	return &PostgresMeetingRepository{db: db}
}

// SaveMeeting stores a meeting in the database
func (r *PostgresMeetingRepository) SaveMeeting(ctx context.Context, meeting *repositories.Meeting) error {
	query := `
		INSERT INTO meetings (id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			status = EXCLUDED.status,
			end_time = EXCLUDED.end_time,
			bot_join_url = EXCLUDED.bot_join_url,
			recording_path = EXCLUDED.recording_path,
			updated_at = EXCLUDED.updated_at
	`

	if meeting.ID == "" {
		meeting.ID = uuid.New().String()
	}

	now := time.Now()
	if meeting.CreatedAt.IsZero() {
		meeting.CreatedAt = now
	}
	meeting.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		meeting.ID,
		meeting.UserID,
		meeting.Title,
		meeting.Type,
		meeting.Status,
		meeting.StartTime,
		meeting.EndTime,
		meeting.MeetingURL,
		meeting.BotJoinURL,
		meeting.RecordingPath,
		meeting.CreatedAt,
		meeting.UpdatedAt,
	)

	return err
}

// FindMeetingByID retrieves a meeting by its ID
func (r *PostgresMeetingRepository) FindMeetingByID(ctx context.Context, id string) (*repositories.Meeting, error) {
	query := `
		SELECT id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at
		FROM meetings
		WHERE id = $1 AND deleted_at IS NULL
	`

	var meeting repositories.Meeting
	var endTime sql.NullTime
	var botJoinURL, recordingPath sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&meeting.ID,
		&meeting.UserID,
		&meeting.Title,
		&meeting.Type,
		&meeting.Status,
		&meeting.StartTime,
		&endTime,
		&meeting.MeetingURL,
		&botJoinURL,
		&recordingPath,
		&meeting.CreatedAt,
		&meeting.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("meeting not found: %s", id)
		}
		return nil, err
	}

	if endTime.Valid {
		meeting.EndTime = &endTime.Time
	}
	if botJoinURL.Valid {
		meeting.BotJoinURL = &botJoinURL.String
	}
	if recordingPath.Valid {
		meeting.RecordingPath = &recordingPath.String
	}

	return &meeting, nil
}

// FindMeetingsByUserID retrieves all meetings for a user
func (r *PostgresMeetingRepository) FindMeetingsByUserID(ctx context.Context, userID string) ([]*repositories.Meeting, error) {
	query := `
		SELECT id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at
		FROM meetings
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY start_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMeetings(rows)
}

// UpdateMeeting updates an existing meeting
func (r *PostgresMeetingRepository) UpdateMeeting(ctx context.Context, meeting *repositories.Meeting) error {
	query := `
		UPDATE meetings 
		SET title = $2, status = $3, end_time = $4, bot_join_url = $5, recording_path = $6, updated_at = $7
		WHERE id = $1
	`

	meeting.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		meeting.ID,
		meeting.Title,
		meeting.Status,
		meeting.EndTime,
		meeting.BotJoinURL,
		meeting.RecordingPath,
		meeting.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("meeting not found: %s", meeting.ID)
	}

	return nil
}

// DeleteMeeting soft deletes a meeting
func (r *PostgresMeetingRepository) DeleteMeeting(ctx context.Context, id string) error {
	query := `UPDATE meetings SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("meeting not found: %s", id)
	}

	return nil
}

// SaveParticipant stores a participant in the database
func (r *PostgresMeetingRepository) SaveParticipant(ctx context.Context, participant *repositories.Participant) error {
	query := `
		INSERT INTO participants (id, meeting_id, name, email, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			role = EXCLUDED.role
	`

	if participant.ID == "" {
		participant.ID = uuid.New().String()
	}

	if participant.CreatedAt.IsZero() {
		participant.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		participant.ID,
		participant.MeetingID,
		participant.Name,
		participant.Email,
		participant.Role,
		participant.CreatedAt,
	)

	return err
}

// FindParticipantsByMeetingID retrieves all participants for a meeting
func (r *PostgresMeetingRepository) FindParticipantsByMeetingID(ctx context.Context, meetingID string) ([]*repositories.Participant, error) {
	query := `
		SELECT id, meeting_id, name, email, role, created_at
		FROM participants
		WHERE meeting_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*repositories.Participant

	for rows.Next() {
		var participant repositories.Participant
		var email, role sql.NullString

		err := rows.Scan(
			&participant.ID,
			&participant.MeetingID,
			&participant.Name,
			&email,
			&role,
			&participant.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if email.Valid {
			participant.Email = &email.String
		}
		if role.Valid {
			participant.Role = &role.String
		}

		participants = append(participants, &participant)
	}

	return participants, rows.Err()
}

// UpdateParticipant updates an existing participant
func (r *PostgresMeetingRepository) UpdateParticipant(ctx context.Context, participant *repositories.Participant) error {
	query := `
		UPDATE participants 
		SET name = $2, email = $3, role = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		participant.ID,
		participant.Name,
		participant.Email,
		participant.Role,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("participant not found: %s", participant.ID)
	}

	return nil
}

// DeleteParticipant removes a participant from the database
func (r *PostgresMeetingRepository) DeleteParticipant(ctx context.Context, id string) error {
	query := `DELETE FROM participants WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("participant not found: %s", id)
	}

	return nil
}

// SaveBotSession stores a bot session in the database
func (r *PostgresMeetingRepository) SaveBotSession(ctx context.Context, session *repositories.BotSession) error {
	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO bot_sessions (id, meeting_id, session_id, status, joined_at, left_at, bot_user_id, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			left_at = EXCLUDED.left_at,
			bot_user_id = EXCLUDED.bot_user_id,
			metadata = EXCLUDED.metadata
	`

	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}

	_, err = r.db.ExecContext(ctx, query,
		session.ID,
		session.MeetingID,
		session.SessionID,
		session.Status,
		session.JoinedAt,
		session.LeftAt,
		session.BotUserID,
		metadataJSON,
		session.CreatedAt,
	)

	return err
}

// FindBotSessionByID retrieves a bot session by its ID
func (r *PostgresMeetingRepository) FindBotSessionByID(ctx context.Context, id string) (*repositories.BotSession, error) {
	query := `
		SELECT id, meeting_id, session_id, status, joined_at, left_at, bot_user_id, metadata, created_at
		FROM bot_sessions
		WHERE id = $1
	`

	var session repositories.BotSession
	var leftAt sql.NullTime
	var botUserID sql.NullString
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.MeetingID,
		&session.SessionID,
		&session.Status,
		&session.JoinedAt,
		&leftAt,
		&botUserID,
		&metadataJSON,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bot session not found: %s", id)
		}
		return nil, err
	}

	if leftAt.Valid {
		session.LeftAt = &leftAt.Time
	}
	if botUserID.Valid {
		session.BotUserID = &botUserID.String
	}

	if len(metadataJSON) > 0 {
		err = json.Unmarshal(metadataJSON, &session.Metadata)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}

// FindBotSessionsByMeetingID retrieves all bot sessions for a meeting
func (r *PostgresMeetingRepository) FindBotSessionsByMeetingID(ctx context.Context, meetingID string) ([]*repositories.BotSession, error) {
	query := `
		SELECT id, meeting_id, session_id, status, joined_at, left_at, bot_user_id, metadata, created_at
		FROM bot_sessions
		WHERE meeting_id = $1
		ORDER BY joined_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBotSessions(rows)
}

// FindBotSessionBySessionID retrieves a bot session by session ID
func (r *PostgresMeetingRepository) FindBotSessionBySessionID(ctx context.Context, sessionID string) (*repositories.BotSession, error) {
	query := `
		SELECT id, meeting_id, session_id, status, joined_at, left_at, bot_user_id, metadata, created_at
		FROM bot_sessions
		WHERE session_id = $1
	`

	var session repositories.BotSession
	var leftAt sql.NullTime
	var botUserID sql.NullString
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.MeetingID,
		&session.SessionID,
		&session.Status,
		&session.JoinedAt,
		&leftAt,
		&botUserID,
		&metadataJSON,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bot session not found for session ID: %s", sessionID)
		}
		return nil, err
	}

	if leftAt.Valid {
		session.LeftAt = &leftAt.Time
	}
	if botUserID.Valid {
		session.BotUserID = &botUserID.String
	}

	if len(metadataJSON) > 0 {
		err = json.Unmarshal(metadataJSON, &session.Metadata)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}

// UpdateBotSession updates an existing bot session
func (r *PostgresMeetingRepository) UpdateBotSession(ctx context.Context, session *repositories.BotSession) error {
	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE bot_sessions 
		SET status = $2, left_at = $3, bot_user_id = $4, metadata = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.Status,
		session.LeftAt,
		session.BotUserID,
		metadataJSON,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bot session not found: %s", session.ID)
	}

	return nil
}

// FindMeetingsByStatus retrieves meetings by status
func (r *PostgresMeetingRepository) FindMeetingsByStatus(ctx context.Context, status string) ([]*repositories.Meeting, error) {
	query := `
		SELECT id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at
		FROM meetings
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY start_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMeetings(rows)
}

// FindActiveMeetings retrieves currently active meetings
func (r *PostgresMeetingRepository) FindActiveMeetings(ctx context.Context) ([]*repositories.Meeting, error) {
	query := `
		SELECT id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at
		FROM meetings
		WHERE status IN ('in_progress', 'scheduled') AND deleted_at IS NULL
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMeetings(rows)
}

// FindMeetingsByTimeRange retrieves meetings within a time range
func (r *PostgresMeetingRepository) FindMeetingsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*repositories.Meeting, error) {
	query := `
		SELECT id, user_id, title, type, status, start_time, end_time, meeting_url, bot_join_url, recording_path, created_at, updated_at
		FROM meetings
		WHERE start_time >= $1 AND start_time <= $2 AND deleted_at IS NULL
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMeetings(rows)
}

// Helper method to scan meetings from rows
func (r *PostgresMeetingRepository) scanMeetings(rows *sql.Rows) ([]*repositories.Meeting, error) {
	var meetings []*repositories.Meeting

	for rows.Next() {
		var meeting repositories.Meeting
		var endTime sql.NullTime
		var botJoinURL, recordingPath sql.NullString

		err := rows.Scan(
			&meeting.ID,
			&meeting.UserID,
			&meeting.Title,
			&meeting.Type,
			&meeting.Status,
			&meeting.StartTime,
			&endTime,
			&meeting.MeetingURL,
			&botJoinURL,
			&recordingPath,
			&meeting.CreatedAt,
			&meeting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endTime.Valid {
			meeting.EndTime = &endTime.Time
		}
		if botJoinURL.Valid {
			meeting.BotJoinURL = &botJoinURL.String
		}
		if recordingPath.Valid {
			meeting.RecordingPath = &recordingPath.String
		}

		meetings = append(meetings, &meeting)
	}

	return meetings, rows.Err()
}

// Helper method to scan bot sessions from rows
func (r *PostgresMeetingRepository) scanBotSessions(rows *sql.Rows) ([]*repositories.BotSession, error) {
	var sessions []*repositories.BotSession

	for rows.Next() {
		var session repositories.BotSession
		var leftAt sql.NullTime
		var botUserID sql.NullString
		var metadataJSON []byte

		err := rows.Scan(
			&session.ID,
			&session.MeetingID,
			&session.SessionID,
			&session.Status,
			&session.JoinedAt,
			&leftAt,
			&botUserID,
			&metadataJSON,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if leftAt.Valid {
			session.LeftAt = &leftAt.Time
		}
		if botUserID.Valid {
			session.BotUserID = &botUserID.String
		}

		if len(metadataJSON) > 0 {
			err = json.Unmarshal(metadataJSON, &session.Metadata)
			if err != nil {
				return nil, err
			}
		}

		sessions = append(sessions, &session)
	}

	return sessions, rows.Err()
}
