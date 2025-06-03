package queries

// GetMeetingsQuery represents the query to get meetings for a user
type GetMeetingsQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// GetMeetingByIDQuery represents the query to get a specific meeting by ID
type GetMeetingByIDQuery struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}
