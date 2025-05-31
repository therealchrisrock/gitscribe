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
	// WebSocket route for audio streaming
	router.GET("/ws/audio", r.audioHandlers.HandleAudioWebSocket)
}
