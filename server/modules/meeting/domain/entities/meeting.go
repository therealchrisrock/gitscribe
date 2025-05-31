package entities

import (
	"teammate/server/seedwork/domain"
	"time"
)

type MeetingType string

const (
	ZoomMeeting    MeetingType = "zoom"
	GoogleMeet     MeetingType = "google_meet"
	MicrosoftTeams MeetingType = "microsoft_teams"
	GenericMeeting MeetingType = "generic"
)

type MeetingStatus string

const (
	Scheduled  MeetingStatus = "scheduled"
	InProgress MeetingStatus = "in_progress"
	Completed  MeetingStatus = "completed"
	Failed     MeetingStatus = "failed"
)

// Meeting represents a meeting entity in the domain
type Meeting struct {
	domain.BaseEntity
	UserID        string        `json:"user_id" gorm:"column:user_id;not null"`
	Title         string        `json:"title" gorm:"column:title;not null"`
	Type          MeetingType   `json:"type" gorm:"column:type;not null"`
	Status        MeetingStatus `json:"status" gorm:"column:status;not null"`
	StartTime     time.Time     `json:"start_time" gorm:"column:start_time;not null"`
	EndTime       *time.Time    `json:"end_time,omitempty" gorm:"column:end_time"`
	MeetingURL    string        `json:"meeting_url" gorm:"column:meeting_url;not null"`
	BotJoinURL    string        `json:"bot_join_url,omitempty" gorm:"column:bot_join_url"`
	RecordingPath string        `json:"recording_path,omitempty" gorm:"column:recording_path"`
	Participants  []Participant `json:"participants" gorm:"foreignKey:MeetingID"`
}

// NewMeeting creates a new Meeting entity
func NewMeeting(userID, title string, meetingType MeetingType, meetingURL string) Meeting {
	meeting := Meeting{
		UserID:     userID,
		Title:      title,
		Type:       meetingType,
		Status:     Scheduled,
		StartTime:  time.Now(),
		MeetingURL: meetingURL,
	}
	meeting.SetID(domain.GenerateID())
	return meeting
}

// StartMeeting transitions the meeting to in-progress status
func (m *Meeting) StartMeeting(botJoinURL string) {
	m.Status = InProgress
	m.BotJoinURL = botJoinURL
}

// CompleteMeeting transitions the meeting to completed status
func (m *Meeting) CompleteMeeting(recordingPath string) {
	m.Status = Completed
	m.RecordingPath = recordingPath
	now := time.Now()
	m.EndTime = &now
}

// FailMeeting transitions the meeting to failed status
func (m *Meeting) FailMeeting() {
	m.Status = Failed
	now := time.Now()
	m.EndTime = &now
}

// IsActive returns true if the meeting is currently in progress
func (m *Meeting) IsActive() bool {
	return m.Status == InProgress
}

// IsCompleted returns true if the meeting has completed successfully
func (m *Meeting) IsCompleted() bool {
	return m.Status == Completed
}

// HasRecording returns true if the meeting has a recording
func (m *Meeting) HasRecording() bool {
	return m.RecordingPath != ""
}

// GetDuration returns the meeting duration if it has ended
func (m *Meeting) GetDuration() *time.Duration {
	if m.EndTime == nil {
		return nil
	}
	duration := m.EndTime.Sub(m.StartTime)
	return &duration
}

// AddParticipant adds a participant to the meeting
func (m *Meeting) AddParticipant(name, email, role string) {
	participant := NewParticipant(m.GetID(), name, email, role)
	m.Participants = append(m.Participants, participant)
}

// TableName sets the table name for GORM
func (Meeting) TableName() string {
	return "meetings"
}

// Participant represents a meeting participant
type Participant struct {
	domain.BaseEntity
	MeetingID string `json:"meeting_id" gorm:"column:meeting_id;not null"`
	Name      string `json:"name" gorm:"column:name;not null"`
	Email     string `json:"email,omitempty" gorm:"column:email"`
	Role      string `json:"role,omitempty" gorm:"column:role"`
}

// NewParticipant creates a new Participant entity
func NewParticipant(meetingID, name, email, role string) Participant {
	participant := Participant{
		MeetingID: meetingID,
		Name:      name,
		Email:     email,
		Role:      role,
	}
	participant.SetID(domain.GenerateID())
	return participant
}

// TableName sets the table name for GORM
func (Participant) TableName() string {
	return "participants"
}

// BotSession represents a bot's participation in a meeting
type BotSession struct {
	domain.BaseEntity
	MeetingID string                 `json:"meeting_id" gorm:"column:meeting_id;not null"`
	SessionID string                 `json:"session_id" gorm:"column:session_id;not null"`
	Status    BotSessionStatus       `json:"status" gorm:"column:status;not null"`
	JoinedAt  time.Time              `json:"joined_at" gorm:"column:joined_at;not null"`
	LeftAt    *time.Time             `json:"left_at,omitempty" gorm:"column:left_at"`
	BotUserID string                 `json:"bot_user_id,omitempty" gorm:"column:bot_user_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

type BotSessionStatus string

const (
	BotJoining   BotSessionStatus = "joining"
	BotActive    BotSessionStatus = "active"
	BotRecording BotSessionStatus = "recording"
	BotCompleted BotSessionStatus = "completed"
	BotFailed    BotSessionStatus = "failed"
)

// NewBotSession creates a new BotSession entity
func NewBotSession(meetingID, sessionID string) BotSession {
	botSession := BotSession{
		MeetingID: meetingID,
		SessionID: sessionID,
		Status:    BotJoining,
		JoinedAt:  time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	botSession.SetID(domain.GenerateID())
	return botSession
}

// StartRecording transitions the bot session to recording status
func (bs *BotSession) StartRecording() {
	bs.Status = BotRecording
}

// Complete transitions the bot session to completed status
func (bs *BotSession) Complete() {
	bs.Status = BotCompleted
	now := time.Now()
	bs.LeftAt = &now
}

// Fail transitions the bot session to failed status
func (bs *BotSession) Fail() {
	bs.Status = BotFailed
	now := time.Now()
	bs.LeftAt = &now
}

// SetBotUserID sets the bot's user ID in the meeting platform
func (bs *BotSession) SetBotUserID(botUserID string) {
	bs.BotUserID = botUserID
}

// AddMetadata adds metadata to the bot session
func (bs *BotSession) AddMetadata(key string, value interface{}) {
	if bs.Metadata == nil {
		bs.Metadata = make(map[string]interface{})
	}
	bs.Metadata[key] = value
}

// IsActive returns true if the bot session is currently active
func (bs *BotSession) IsActive() bool {
	return bs.Status == BotActive || bs.Status == BotRecording
}

// TableName sets the table name for GORM
func (BotSession) TableName() string {
	return "bot_sessions"
}
