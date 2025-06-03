package dtos

import (
	"time"

	"teammate/server/modules/meeting/domain/entities"
)

// CreateMeetingRequest represents the request to create a meeting
type CreateMeetingRequest struct {
	UserID     string               `json:"-"` // Set from authentication context
	Title      string               `json:"title" binding:"required"`
	Type       entities.MeetingType `json:"type" binding:"required"`
	MeetingURL string               `json:"meeting_url" binding:"required,url"`
}

// MeetingResponse represents the response containing meeting data
type MeetingResponse struct {
	ID            string                 `json:"id"`
	UserID        string                 `json:"user_id"`
	Title         string                 `json:"title"`
	Type          entities.MeetingType   `json:"type"`
	Status        entities.MeetingStatus `json:"status"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       *time.Time             `json:"end_time,omitempty"`
	MeetingURL    string                 `json:"meeting_url"`
	BotJoinURL    string                 `json:"bot_join_url,omitempty"`
	RecordingPath string                 `json:"recording_path,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// MeetingsListResponse represents the response containing a list of meetings
type MeetingsListResponse struct {
	Meetings []MeetingResponse `json:"meetings"`
	Total    int64             `json:"total"`
}

// ToMeetingResponse converts a Meeting entity to MeetingResponse DTO
func ToMeetingResponse(meeting *entities.Meeting) MeetingResponse {
	return MeetingResponse{
		ID:            meeting.GetID(),
		UserID:        meeting.UserID,
		Title:         meeting.Title,
		Type:          meeting.Type,
		Status:        meeting.Status,
		StartTime:     meeting.StartTime,
		EndTime:       meeting.EndTime,
		MeetingURL:    meeting.MeetingURL,
		BotJoinURL:    meeting.BotJoinURL,
		RecordingPath: meeting.RecordingPath,
		CreatedAt:     meeting.GetCreatedAt(),
		UpdatedAt:     meeting.GetUpdatedAt(),
	}
}

// ToMeetingsListResponse converts a slice of Meeting entities to MeetingsListResponse DTO
func ToMeetingsListResponse(meetings []*entities.Meeting, total int64) MeetingsListResponse {
	meetingResponses := make([]MeetingResponse, len(meetings))
	for i, meeting := range meetings {
		meetingResponses[i] = ToMeetingResponse(meeting)
	}

	return MeetingsListResponse{
		Meetings: meetingResponses,
		Total:    total,
	}
}
