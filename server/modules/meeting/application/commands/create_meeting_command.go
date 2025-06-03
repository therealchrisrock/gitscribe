package commands

import (
	"teammate/server/modules/meeting/domain/entities"
)

// CreateMeetingCommand represents the command to create a new meeting
type CreateMeetingCommand struct {
	UserID     string               `json:"user_id" validate:"required"`
	Title      string               `json:"title" validate:"required"`
	Type       entities.MeetingType `json:"type" validate:"required"`
	MeetingURL string               `json:"meeting_url" validate:"required,url"`
}
