package repositories

import (
	"teammate/server/modules/meeting/domain/entities"
)

// MeetingRepository defines the interface for meeting data access
type MeetingRepository interface {
	// Base repository methods
	FindAll() ([]*entities.Meeting, error)
	FindByID(id string) (*entities.Meeting, error)
	Create(entity *entities.Meeting) error
	Update(entity *entities.Meeting) error
	Delete(id string) error
	Count() (int64, error)

	// Domain-specific query methods
	FindByUserID(userID string) ([]*entities.Meeting, error)
	FindByStatus(status entities.MeetingStatus) ([]*entities.Meeting, error)
	FindByType(meetingType entities.MeetingType) ([]*entities.Meeting, error)
	FindActiveByUserID(userID string) ([]*entities.Meeting, error)
	FindCompletedByUserID(userID string) ([]*entities.Meeting, error)
	FindByDateRange(userID string, startDate, endDate string) ([]*entities.Meeting, error)
}

// ParticipantRepository defines the interface for participant data access
type ParticipantRepository interface {
	// Base repository methods
	FindAll() ([]*entities.Participant, error)
	FindByID(id string) (*entities.Participant, error)
	Create(entity *entities.Participant) error
	Update(entity *entities.Participant) error
	Delete(id string) error
	Count() (int64, error)

	// Domain-specific query methods
	FindByMeetingID(meetingID string) ([]*entities.Participant, error)
	FindByEmail(email string) ([]*entities.Participant, error)
}

// BotSessionRepository defines the interface for bot session data access
type BotSessionRepository interface {
	// Base repository methods
	FindAll() ([]*entities.BotSession, error)
	FindByID(id string) (*entities.BotSession, error)
	Create(entity *entities.BotSession) error
	Update(entity *entities.BotSession) error
	Delete(id string) error
	Count() (int64, error)

	// Domain-specific query methods
	FindByMeetingID(meetingID string) ([]*entities.BotSession, error)
	FindBySessionID(sessionID string) (*entities.BotSession, error)
	FindByStatus(status entities.BotSessionStatus) ([]*entities.BotSession, error)
	FindActiveSessions() ([]*entities.BotSession, error)
}
