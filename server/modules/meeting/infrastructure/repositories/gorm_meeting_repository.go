package repositories

import (
	"context"
	"fmt"
	"time"

	"teammate/server/modules/meeting/domain/repositories"
	"teammate/server/seedwork/infrastructure/database"

	"gorm.io/gorm"
)

// GormMeetingRepository implements MeetingRepository using GORM
type GormMeetingRepository struct {
	db *gorm.DB
}

// NewGormMeetingRepository creates a new GORM meeting repository
func NewGormMeetingRepository() *GormMeetingRepository {
	return &GormMeetingRepository{db: database.GetDB()}
}

// SaveMeeting stores a meeting in the database
func (r *GormMeetingRepository) SaveMeeting(ctx context.Context, meeting *repositories.Meeting) error {
	return r.db.WithContext(ctx).Save(meeting).Error
}

// FindMeetingByID retrieves a meeting by its ID
func (r *GormMeetingRepository) FindMeetingByID(ctx context.Context, id string) (*repositories.Meeting, error) {
	var meeting repositories.Meeting
	err := r.db.WithContext(ctx).First(&meeting, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

// UpdateMeeting updates an existing meeting
func (r *GormMeetingRepository) UpdateMeeting(ctx context.Context, meeting *repositories.Meeting) error {
	result := r.db.WithContext(ctx).Save(meeting)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("meeting not found: %s", meeting.ID)
	}
	return nil
}

// DeleteMeeting removes a meeting from the database
func (r *GormMeetingRepository) DeleteMeeting(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&repositories.Meeting{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("meeting not found: %s", id)
	}
	return nil
}

// SaveBotSession stores a bot session
func (r *GormMeetingRepository) SaveBotSession(ctx context.Context, session *repositories.BotSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// FindBotSessionByID retrieves a bot session by ID
func (r *GormMeetingRepository) FindBotSessionByID(ctx context.Context, id string) (*repositories.BotSession, error) {
	var session repositories.BotSession
	err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateBotSession updates a bot session
func (r *GormMeetingRepository) UpdateBotSession(ctx context.Context, session *repositories.BotSession) error {
	result := r.db.WithContext(ctx).Save(session)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("bot session not found: %s", session.ID)
	}
	return nil
}

// FindBotSessionsByMeetingID retrieves bot sessions for a meeting
func (r *GormMeetingRepository) FindBotSessionsByMeetingID(ctx context.Context, meetingID string) ([]*repositories.BotSession, error) {
	var sessions []*repositories.BotSession
	err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Find(&sessions).Error
	return sessions, err
}

// FindMeetingsByUserID retrieves meetings for a user
func (r *GormMeetingRepository) FindMeetingsByUserID(ctx context.Context, userID string) ([]*repositories.Meeting, error) {
	var meetings []*repositories.Meeting
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&meetings).Error
	return meetings, err
}

// SaveParticipant stores a participant
func (r *GormMeetingRepository) SaveParticipant(ctx context.Context, participant *repositories.Participant) error {
	return r.db.WithContext(ctx).Save(participant).Error
}

// FindParticipantsByMeetingID retrieves participants for a meeting
func (r *GormMeetingRepository) FindParticipantsByMeetingID(ctx context.Context, meetingID string) ([]*repositories.Participant, error) {
	var participants []*repositories.Participant
	err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Find(&participants).Error
	return participants, err
}

// UpdateParticipant updates a participant
func (r *GormMeetingRepository) UpdateParticipant(ctx context.Context, participant *repositories.Participant) error {
	result := r.db.WithContext(ctx).Save(participant)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("participant not found: %s", participant.ID)
	}
	return nil
}

// DeleteParticipant removes a participant
func (r *GormMeetingRepository) DeleteParticipant(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&repositories.Participant{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("participant not found: %s", id)
	}
	return nil
}

// FindBotSessionBySessionID retrieves a bot session by session ID
func (r *GormMeetingRepository) FindBotSessionBySessionID(ctx context.Context, sessionID string) (*repositories.BotSession, error) {
	var session repositories.BotSession
	err := r.db.WithContext(ctx).First(&session, "session_id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindMeetingsByStatus retrieves meetings by status
func (r *GormMeetingRepository) FindMeetingsByStatus(ctx context.Context, status string) ([]*repositories.Meeting, error) {
	var meetings []*repositories.Meeting
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&meetings).Error
	return meetings, err
}

// FindActiveMeetings retrieves active meetings
func (r *GormMeetingRepository) FindActiveMeetings(ctx context.Context) ([]*repositories.Meeting, error) {
	return r.FindMeetingsByStatus(ctx, repositories.MeetingStatusInProgress)
}

// FindMeetingsByTimeRange retrieves meetings in a time range
func (r *GormMeetingRepository) FindMeetingsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*repositories.Meeting, error) {
	var meetings []*repositories.Meeting
	err := r.db.WithContext(ctx).Where("start_time BETWEEN ? AND ?", startTime, endTime).Find(&meetings).Error
	return meetings, err
}
