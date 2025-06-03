package routes

import (
	"teammate/server/modules/meeting/interfaces/http/handlers"
	"teammate/server/modules/user/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// MeetingRoutes sets up all meeting-related routes
type MeetingRoutes struct {
	meetingHandlers *handlers.MeetingHandlers
	authMiddleware  *middleware.AuthMiddleware
}

// NewMeetingRoutes creates a new meeting routes instance
func NewMeetingRoutes(meetingHandlers *handlers.MeetingHandlers, authMiddleware *middleware.AuthMiddleware) *MeetingRoutes {
	return &MeetingRoutes{
		meetingHandlers: meetingHandlers,
		authMiddleware:  authMiddleware,
	}
}

// SetupProtectedRoutes sets up protected meeting routes (authentication required)
func (mr *MeetingRoutes) SetupProtectedRoutes(protected *gin.RouterGroup) {
	// Apply auth middleware to protected routes
	protected.Use(mr.authMiddleware.FirebaseAuth())

	// Meeting endpoints
	meetings := protected.Group("/meetings")
	{
		meetings.POST("", mr.meetingHandlers.CreateMeeting)     // Create new meeting
		meetings.GET("", mr.meetingHandlers.GetMeetings)        // Get user's meetings
		meetings.GET("/:id", mr.meetingHandlers.GetMeetingByID) // Get specific meeting
	}
}
