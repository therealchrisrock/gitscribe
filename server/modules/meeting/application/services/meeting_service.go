package services

import (
	"context"
	"fmt"

	"teammate/server/modules/meeting/application/commands"
	"teammate/server/modules/meeting/application/queries"
	"teammate/server/modules/meeting/domain/entities"
	"teammate/server/modules/meeting/domain/repositories"
	"teammate/server/seedwork/domain"
)

// MeetingService handles meeting business logic orchestration
type MeetingService struct {
	meetingRepo   repositories.MeetingRepository
	meetingMapper *MeetingMapper
}

// NewMeetingService creates a new meeting service
func NewMeetingService(meetingRepo repositories.MeetingRepository) *MeetingService {
	return &MeetingService{
		meetingRepo:   meetingRepo,
		meetingMapper: NewMeetingMapper(),
	}
}

// CreateMeeting creates a new meeting using domain factory method
func (s *MeetingService) CreateMeeting(ctx context.Context, cmd commands.CreateMeetingCommand) (*entities.Meeting, error) {
	// Use domain factory method with direct parameters
	meeting, err := entities.CreateMeeting(cmd.UserID, cmd.Title, cmd.Type, cmd.MeetingURL)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	// Convert to repository model using mapper
	repoMeeting := s.meetingMapper.ToRepository(meeting)

	if err := s.meetingRepo.SaveMeeting(ctx, &repoMeeting); err != nil {
		return nil, fmt.Errorf("failed to persist meeting: %w", err)
	}

	return meeting, nil
}

// StartMeeting starts a meeting using domain method
func (s *MeetingService) StartMeeting(ctx context.Context, meetingID, userID, botJoinURL string) error {
	// Load aggregate
	meeting, err := s.getMeetingAggregate(ctx, meetingID, userID)
	if err != nil {
		return err
	}

	// Execute domain method with direct parameters
	if err := meeting.Start(botJoinURL); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	// Persist changes using mapper
	return s.saveMeetingAggregate(ctx, meeting)
}

// CompleteMeeting completes a meeting using domain method
func (s *MeetingService) CompleteMeeting(ctx context.Context, meetingID, userID, recordingPath string) error {
	// Load aggregate
	meeting, err := s.getMeetingAggregate(ctx, meetingID, userID)
	if err != nil {
		return err
	}

	// Execute domain method with direct parameters
	if err := meeting.Complete(recordingPath); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	// Persist changes using mapper
	return s.saveMeetingAggregate(ctx, meeting)
}

// GetMeetings retrieves meetings for a user (query - no domain logic needed)
func (s *MeetingService) GetMeetings(ctx context.Context, query queries.GetMeetingsQuery) ([]*entities.Meeting, int64, error) {
	repoMeetings, err := s.meetingRepo.FindMeetingsByUserID(ctx, query.UserID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get meetings: %w", err)
	}

	// Convert repository models to domain entities using mapper
	// Convert pointers to values for the mapper
	repoValues := make([]repositories.Meeting, len(repoMeetings))
	for i, repoMeeting := range repoMeetings {
		repoValues[i] = *repoMeeting
	}
	meetings := s.meetingMapper.ToDomainList(repoValues)

	return meetings, int64(len(meetings)), nil
}

// GetMeetingByID retrieves a specific meeting by ID (query - no domain logic needed)
func (s *MeetingService) GetMeetingByID(ctx context.Context, query queries.GetMeetingByIDQuery) (*entities.Meeting, error) {
	repoMeeting, err := s.meetingRepo.FindMeetingByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meeting: %w", err)
	}

	// Verify ownership
	if repoMeeting.UserID != query.UserID {
		return nil, domain.NewDomainError("UNAUTHORIZED", "Meeting not found or access denied", nil)
	}

	// Convert repository model to domain entity using mapper
	return s.meetingMapper.ToDomain(*repoMeeting), nil
}

// Helper methods for aggregate loading/saving
func (s *MeetingService) getMeetingAggregate(ctx context.Context, meetingID, userID string) (*entities.Meeting, error) {
	repoMeeting, err := s.meetingRepo.FindMeetingByID(ctx, meetingID)
	if err != nil {
		return nil, fmt.Errorf("failed to load meeting: %w", err)
	}

	// Verify ownership
	if repoMeeting.UserID != userID {
		return nil, domain.NewDomainError("UNAUTHORIZED", "Meeting not found or access denied", nil)
	}

	// Convert repository model to domain entity using mapper
	return s.meetingMapper.ToDomain(*repoMeeting), nil
}

func (s *MeetingService) saveMeetingAggregate(ctx context.Context, meeting *entities.Meeting) error {
	// Convert domain entity to repository model using mapper
	repoMeeting := s.meetingMapper.ToRepository(meeting)
	return s.meetingRepo.UpdateMeeting(ctx, &repoMeeting)
}
