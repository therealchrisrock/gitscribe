package routes

import (
	"teammate/server/modules/transcription/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

type TranscriptionRoutes struct {
	audioHandlers *handlers.AudioHandlers
}

func NewTranscriptionRoutes(audioHandlers *handlers.AudioHandlers) *TranscriptionRoutes {
	return &TranscriptionRoutes{
		audioHandlers: audioHandlers,
	}
}

// SetupRoutes sets up all transcription-related routes
func (r *TranscriptionRoutes) SetupRoutes(router *gin.RouterGroup) {
	// Audio processing routes
	audioGroup := router.Group("/audio")
	{
		// REST endpoint for provider capabilities
		audioGroup.GET("/providers", r.audioHandlers.GetProviderCapabilities)
	}

	// WebSocket route for audio streaming (both real-time and batch)
	router.GET("/ws/audio", r.audioHandlers.HandleAudioWebSocket)
}
