package handlers

import (
	"net/http"

	"teammate/server/modules/meeting/application/commands"
	"teammate/server/modules/meeting/application/queries"
	"teammate/server/modules/meeting/application/services"
	"teammate/server/modules/meeting/interfaces/http/dtos"
	"teammate/server/modules/user/domain/entities"
	"teammate/server/seedwork/domain"

	"github.com/gin-gonic/gin"
)

// MeetingHandlers contains all meeting-related HTTP handlers
type MeetingHandlers struct {
	meetingService *services.MeetingService
}

// NewMeetingHandlers creates a new meeting handlers instance
func NewMeetingHandlers(meetingService *services.MeetingService) *MeetingHandlers {
	return &MeetingHandlers{
		meetingService: meetingService,
	}
}

// CreateMeeting creates a new meeting
// @Summary Create a new meeting
// @Description Create a new meeting for the authenticated user
// @Tags meetings
// @Accept json
// @Produce json
// @Param meeting body dtos.CreateMeetingRequest true "Meeting information"
// @Success 201 {object} dtos.MeetingResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meetings [post]
func (h *MeetingHandlers) CreateMeeting(c *gin.Context) {
	var req dtos.CreateMeetingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, ok := userInterface.(*entities.User)
	if !ok {
		// Try to get user ID directly if user object is not available
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not available"})
			return
		}
		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		req.UserID = userID
	} else {
		req.UserID = user.GetID()
	}

	cmd := commands.CreateMeetingCommand{
		UserID:     req.UserID,
		Title:      req.Title,
		Type:       req.Type,
		MeetingURL: req.MeetingURL,
	}

	meeting, err := h.meetingService.CreateMeeting(c.Request.Context(), cmd)
	if err != nil {
		// Handle domain errors
		if domainErr, ok := err.(*domain.DomainError); ok {
			switch domainErr.Code {
			case "INVALID_MEETING_URL":
				c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
		return
	}

	response := dtos.ToMeetingResponse(meeting)
	c.JSON(http.StatusCreated, response)
}

// GetMeetings returns meetings for the authenticated user
// @Summary Get user meetings
// @Description Get all meetings for the authenticated user
// @Tags meetings
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dtos.MeetingsListResponse
// @Failure 500 {object} map[string]string
// @Router /meetings [get]
func (h *MeetingHandlers) GetMeetings(c *gin.Context) {
	// Get user ID from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, ok := userInterface.(*entities.User)
	if !ok {
		// Try to get user ID directly if user object is not available
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not available"})
			return
		}
		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		query := queries.GetMeetingsQuery{
			UserID: userID,
		}
		meetings, total, err := h.meetingService.GetMeetings(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meetings"})
			return
		}

		response := dtos.ToMeetingsListResponse(meetings, total)
		c.JSON(http.StatusOK, response)
		return
	}

	query := queries.GetMeetingsQuery{
		UserID: user.GetID(),
	}

	meetings, total, err := h.meetingService.GetMeetings(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meetings"})
		return
	}

	response := dtos.ToMeetingsListResponse(meetings, total)
	c.JSON(http.StatusOK, response)
}

// GetMeetingByID returns a specific meeting by ID
// @Summary Get a meeting by ID
// @Description Get a meeting by its ID for the authenticated user
// @Tags meetings
// @Accept json
// @Produce json
// @Param id path string true "Meeting ID"
// @Success 200 {object} dtos.MeetingResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /meetings/{id} [get]
func (h *MeetingHandlers) GetMeetingByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting ID is required"})
		return
	}

	// Get user ID from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, ok := userInterface.(*entities.User)
	if !ok {
		// Try to get user ID directly if user object is not available
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not available"})
			return
		}
		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		query := queries.GetMeetingByIDQuery{
			ID:     id,
			UserID: userID,
		}
		meeting, err := h.meetingService.GetMeetingByID(c.Request.Context(), query)
		if err != nil {
			if domainErr, ok := err.(*domain.DomainError); ok && domainErr.Code == "UNAUTHORIZED" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meeting"})
			return
		}

		response := dtos.ToMeetingResponse(meeting)
		c.JSON(http.StatusOK, response)
		return
	}

	query := queries.GetMeetingByIDQuery{
		ID:     id,
		UserID: user.GetID(),
	}

	meeting, err := h.meetingService.GetMeetingByID(c.Request.Context(), query)
	if err != nil {
		if domainErr, ok := err.(*domain.DomainError); ok && domainErr.Code == "UNAUTHORIZED" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meeting"})
		return
	}

	response := dtos.ToMeetingResponse(meeting)
	c.JSON(http.StatusOK, response)
}
