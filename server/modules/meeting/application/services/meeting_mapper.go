package services

import (
	"teammate/server/modules/meeting/domain/entities"
	"teammate/server/modules/meeting/domain/repositories"
	"teammate/server/seedwork/domain"
)

// MeetingMapper implements DomainMapper for Meeting entities
type MeetingMapper struct {
	domain.BaseDomainMapper
}

// NewMeetingMapper creates a new meeting mapper
func NewMeetingMapper() *MeetingMapper {
	return &MeetingMapper{}
}

// ToRepository converts domain Meeting to repository Meeting
func (m *MeetingMapper) ToRepository(meeting *entities.Meeting) repositories.Meeting {
	repo := repositories.Meeting{
		UserID:        meeting.UserID,
		Title:         meeting.Title,
		Type:          string(meeting.Type),
		Status:        string(meeting.Status),
		StartTime:     meeting.StartTime,
		EndTime:       meeting.EndTime,
		MeetingURL:    meeting.MeetingURL,
		BotJoinURL:    m.StringToPointer(meeting.BotJoinURL),
		RecordingPath: m.StringToPointer(meeting.RecordingPath),
	}

	// Set repository model fields
	repo.SetID(meeting.GetID())
	repo.CreatedAt = meeting.GetCreatedAt()
	repo.UpdatedAt = meeting.GetUpdatedAt()

	return repo
}

// ToDomain converts repository Meeting to domain Meeting
func (m *MeetingMapper) ToDomain(repo repositories.Meeting) *entities.Meeting {
	meeting := &entities.Meeting{
		UserID:        repo.UserID,
		Title:         repo.Title,
		Type:          entities.MeetingType(repo.Type),
		Status:        entities.MeetingStatus(repo.Status),
		StartTime:     repo.StartTime,
		MeetingURL:    repo.MeetingURL,
		BotJoinURL:    m.PointerToString(repo.BotJoinURL),
		RecordingPath: m.PointerToString(repo.RecordingPath),
	}

	// Set entity metadata
	meeting.SetID(repo.GetID())
	meeting.CreatedAt = repo.CreatedAt
	meeting.UpdatedAt = repo.UpdatedAt

	// Handle EndTime pointer
	if repo.EndTime != nil {
		meeting.EndTime = repo.EndTime
	}

	return meeting
}

// ToRepositoryList converts slice of domain Meetings to repository Meetings
func (m *MeetingMapper) ToRepositoryList(meetings []*entities.Meeting) []repositories.Meeting {
	result := make([]repositories.Meeting, len(meetings))
	for i, meeting := range meetings {
		result[i] = m.ToRepository(meeting)
	}
	return result
}

// ToDomainList converts slice of repository Meetings to domain Meetings
func (m *MeetingMapper) ToDomainList(repos []repositories.Meeting) []*entities.Meeting {
	result := make([]*entities.Meeting, len(repos))
	for i := range repos {
		result[i] = m.ToDomain(repos[i])
	}
	return result
}

// ParticipantMapper implements DomainMapper for Participant entities
type ParticipantMapper struct {
	domain.BaseDomainMapper
}

// NewParticipantMapper creates a new participant mapper
func NewParticipantMapper() *ParticipantMapper {
	return &ParticipantMapper{}
}

// ToRepository converts domain Participant to repository Participant
func (m *ParticipantMapper) ToRepository(participant *entities.Participant) repositories.Participant {
	repo := repositories.Participant{
		MeetingID: participant.MeetingID,
		Name:      participant.Name,
		Email:     m.StringToPointer(participant.Email),
		Role:      m.StringToPointer(participant.Role),
	}

	// Set repository model fields
	repo.SetID(participant.GetID())
	repo.CreatedAt = participant.GetCreatedAt()
	repo.UpdatedAt = participant.GetUpdatedAt()

	return repo
}

// ToDomain converts repository Participant to domain Participant
func (m *ParticipantMapper) ToDomain(repo repositories.Participant) *entities.Participant {
	participant := &entities.Participant{
		MeetingID: repo.MeetingID,
		Name:      repo.Name,
		Email:     m.PointerToString(repo.Email),
		Role:      m.PointerToString(repo.Role),
	}

	participant.SetID(repo.GetID())
	participant.CreatedAt = repo.CreatedAt
	participant.UpdatedAt = repo.UpdatedAt

	return participant
}

// ToRepositoryList converts slice of domain Participants to repository Participants
func (m *ParticipantMapper) ToRepositoryList(participants []*entities.Participant) []repositories.Participant {
	result := make([]repositories.Participant, len(participants))
	for i, participant := range participants {
		result[i] = m.ToRepository(participant)
	}
	return result
}

// ToDomainList converts slice of repository Participants to domain Participants
func (m *ParticipantMapper) ToDomainList(repos []repositories.Participant) []*entities.Participant {
	result := make([]*entities.Participant, len(repos))
	for i := range repos {
		result[i] = m.ToDomain(repos[i])
	}
	return result
}

// BotSessionMapper implements DomainMapper for BotSession entities
type BotSessionMapper struct {
	domain.BaseDomainMapper
}

// NewBotSessionMapper creates a new bot session mapper
func NewBotSessionMapper() *BotSessionMapper {
	return &BotSessionMapper{}
}

// ToRepository converts domain BotSession to repository BotSession
func (m *BotSessionMapper) ToRepository(session *entities.BotSession) repositories.BotSession {
	repo := repositories.BotSession{
		MeetingID: session.MeetingID,
		SessionID: session.SessionID,
		Status:    string(session.Status),
		JoinedAt:  session.JoinedAt,
		LeftAt:    session.LeftAt,
		BotUserID: m.StringToPointer(session.BotUserID),
		Metadata:  session.Metadata,
	}

	// Set repository model fields
	repo.SetID(session.GetID())
	repo.CreatedAt = session.GetCreatedAt()
	repo.UpdatedAt = session.GetUpdatedAt()

	return repo
}

// ToDomain converts repository BotSession to domain BotSession
func (m *BotSessionMapper) ToDomain(repo repositories.BotSession) *entities.BotSession {
	session := &entities.BotSession{
		MeetingID: repo.MeetingID,
		SessionID: repo.SessionID,
		Status:    entities.BotSessionStatus(repo.Status),
		JoinedAt:  repo.JoinedAt,
		LeftAt:    repo.LeftAt,
		BotUserID: m.PointerToString(repo.BotUserID),
		Metadata:  repo.Metadata,
	}

	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}

	session.SetID(repo.GetID())
	session.CreatedAt = repo.CreatedAt
	session.UpdatedAt = repo.UpdatedAt

	return session
}

// ToRepositoryList converts slice of domain BotSessions to repository BotSessions
func (m *BotSessionMapper) ToRepositoryList(sessions []*entities.BotSession) []repositories.BotSession {
	result := make([]repositories.BotSession, len(sessions))
	for i, session := range sessions {
		result[i] = m.ToRepository(session)
	}
	return result
}

// ToDomainList converts slice of repository BotSessions to domain BotSessions
func (m *BotSessionMapper) ToDomainList(repos []repositories.BotSession) []*entities.BotSession {
	result := make([]*entities.BotSession, len(repos))
	for i := range repos {
		result[i] = m.ToDomain(repos[i])
	}
	return result
}
