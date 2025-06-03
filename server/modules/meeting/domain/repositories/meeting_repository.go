package repositories

import (
	"context"
	"teammate/server/seedwork/domain"
	"time"
)

// Meeting represents a meeting entity for database operations
type Meeting struct {
	domain.BaseRepositoryModel
	UserID        string     `json:"user_id"`
	Title         string     `json:"title"`
	Type          string     `json:"type"`
	Status        string     `json:"status"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	MeetingURL    string     `json:"meeting_url"`
	BotJoinURL    *string    `json:"bot_join_url,omitempty"`
	RecordingPath *string    `json:"recording_path,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the database table name for meetings
func (Meeting) TableName() string {
	return "meetings"
}

// Participant represents a meeting participant
type Participant struct {
	domain.BaseRepositoryModel
	MeetingID string  `json:"meeting_id"`
	Name      string  `json:"name"`
	Email     *string `json:"email,omitempty"`
	Role      *string `json:"role,omitempty"`
}

// TableName returns the database table name for participants
func (Participant) TableName() string {
	return "participants"
}

// BotSession represents a bot session in a meeting
type BotSession struct {
	domain.BaseRepositoryModel
	MeetingID string                 `json:"meeting_id"`
	SessionID string                 `json:"session_id"`
	Status    string                 `json:"status"`
	JoinedAt  time.Time              `json:"joined_at"`
	LeftAt    *time.Time             `json:"left_at,omitempty"`
	BotUserID *string                `json:"bot_user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// TableName returns the database table name for bot sessions
func (BotSession) TableName() string {
	return "bot_sessions"
}

// MeetingRepository defines the interface for meeting persistence
type MeetingRepository interface {
	// Meeting CRUD operations
	SaveMeeting(ctx context.Context, meeting *Meeting) error
	FindMeetingByID(ctx context.Context, id string) (*Meeting, error)
	FindMeetingsByUserID(ctx context.Context, userID string) ([]*Meeting, error)
	UpdateMeeting(ctx context.Context, meeting *Meeting) error
	DeleteMeeting(ctx context.Context, id string) error

	// Participant operations
	SaveParticipant(ctx context.Context, participant *Participant) error
	FindParticipantsByMeetingID(ctx context.Context, meetingID string) ([]*Participant, error)
	UpdateParticipant(ctx context.Context, participant *Participant) error
	DeleteParticipant(ctx context.Context, id string) error

	// Bot session operations
	SaveBotSession(ctx context.Context, session *BotSession) error
	FindBotSessionByID(ctx context.Context, id string) (*BotSession, error)
	FindBotSessionsByMeetingID(ctx context.Context, meetingID string) ([]*BotSession, error)
	FindBotSessionBySessionID(ctx context.Context, sessionID string) (*BotSession, error)
	UpdateBotSession(ctx context.Context, session *BotSession) error

	// Query operations
	FindMeetingsByStatus(ctx context.Context, status string) ([]*Meeting, error)
	FindActiveMeetings(ctx context.Context) ([]*Meeting, error)
	FindMeetingsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*Meeting, error)
}

// Meeting status constants
const (
	MeetingStatusScheduled  = "scheduled"
	MeetingStatusInProgress = "in_progress"
	MeetingStatusCompleted  = "completed"
	MeetingStatusFailed     = "failed"
)

// Bot session status constants
const (
	BotSessionStatusJoining   = "joining"
	BotSessionStatusActive    = "active"
	BotSessionStatusRecording = "recording"
	BotSessionStatusCompleted = "completed"
	BotSessionStatusFailed    = "failed"
)

// Meeting type constants
const (
	MeetingTypeZoom           = "zoom"
	MeetingTypeGoogleMeet     = "google_meet"
	MeetingTypeMicrosoftTeams = "microsoft_teams"
	MeetingTypeGeneric        = "generic"
)
