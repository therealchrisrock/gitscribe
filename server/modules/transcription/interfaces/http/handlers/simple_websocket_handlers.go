package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	appServices "teammate/server/modules/transcription/application/services"
	"teammate/server/modules/transcription/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// In production, you should restrict this to your domain
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type AudioHandlers struct {
	processorFactory services.AudioProcessorFactory
}

func NewAudioHandlers(factory services.AudioProcessorFactory) *AudioHandlers {
	return &AudioHandlers{
		processorFactory: factory,
	}
}

// CreateDefaultAudioHandlers creates handlers with a default factory
func CreateDefaultAudioHandlers() *AudioHandlers {
	factory := appServices.NewAudioProcessorFactory()
	return NewAudioHandlers(factory)
}

// AudioMessage represents the structure of messages sent over websocket
type AudioMessage struct {
	Type      string                           `json:"type"`
	SessionID string                           `json:"session_id,omitempty"`
	Metadata  *services.AudioStreamMetadata    `json:"metadata,omitempty"`
	Options   *services.AudioProcessingOptions `json:"options,omitempty"`
	Chunk     *services.AudioChunk             `json:"chunk,omitempty"`
	Error     string                           `json:"error,omitempty"`
	Result    *services.AudioProcessingResult  `json:"result,omitempty"`
	Message   string                           `json:"message,omitempty"`
}

// HandleAudioWebSocket handles incoming websocket connections for audio streaming
// @Summary Handle audio websocket connection
// @Description Handle websocket connections for real-time or batch audio processing
// @Tags audio
// @Accept json
// @Produce json
// @Param mode query string false "Processing mode (realtime or batch)" default(realtime)
// @Param provider query string false "Transcription provider" default(assemblyai)
// @Param language query string false "Language code" default(en)
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/audio [get]
func (h *AudioHandlers) HandleAudioWebSocket(c *gin.Context) {
	// Parse query parameters for processing configuration
	modeParam := c.DefaultQuery("mode", "batch")
	provider := c.DefaultQuery("provider", "assemblyai")
	language := c.DefaultQuery("language", "en")
	costOptimized := c.DefaultQuery("cost_optimized", "false") == "true"
	realTimeTranscription := c.DefaultQuery("real_time", "true") == "true"
	speakerDiarization := c.DefaultQuery("speaker_diarization", "false") == "true"

	// Parse processing mode
	var mode services.ProcessingMode
	switch modeParam {
	case "batch":
		mode = services.BatchMode
	case "realtime":
		mode = services.RealTimeMode
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid processing mode. Use 'realtime' or 'batch'"})
		return
	}

	// Create processing options
	options := services.AudioProcessingOptions{
		Mode:                  mode,
		Provider:              provider,
		Language:              language,
		CostOptimized:         costOptimized,
		RealTimeTranscription: realTimeTranscription && mode == services.RealTimeMode, // Only enable real-time for real-time mode
		SpeakerDiarization:    speakerDiarization,
		ConfidenceThreshold:   0.7,
	}

	// Get provider capabilities for validation
	capabilities, err := h.processorFactory.GetProviderCapabilities(provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown provider: " + provider})
		return
	}

	// Validate that the provider supports the requested mode
	supportedModes := capabilities.SupportedModes
	modeSupported := false
	for _, supportedMode := range supportedModes {
		if supportedMode == mode {
			modeSupported = true
			break
		}
	}
	if !modeSupported {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Provider " + provider + " does not support " + string(mode) + " mode",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	// Create audio processor
	processor, err := h.processorFactory.CreateProcessor(mode, options)
	if err != nil {
		log.Printf("Failed to create audio processor: %v", err)
		h.sendError(conn, "", "Failed to create audio processor: "+err.Error())
		return
	}

	log.Printf("WebSocket connection established for %s mode with provider %s", mode, provider)

	// Send welcome message with capabilities
	welcomeMsg := AudioMessage{
		Type:    "connected",
		Message: "Connected to audio processing service",
		Result: &services.AudioProcessingResult{
			ProcessingMode: mode,
			Message:        "Ready to process audio",
		},
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
		return
	}

	var currentSessionID string
	ctx := context.Background()

	// Handle incoming messages
	for {
		var msg AudioMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "start_session":
			currentSessionID, err = h.handleStartSession(ctx, processor, msg, options)
			if err != nil {
				h.sendError(conn, currentSessionID, "Failed to start session: "+err.Error())
				continue
			}

			response := AudioMessage{
				Type:      "session_started",
				SessionID: currentSessionID,
				Message:   "Audio processing session started",
			}
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Failed to send session_started response: %v", err)
			}

		case "audio_chunk":
			if currentSessionID == "" {
				h.sendError(conn, "", "No active session. Start a session first.")
				continue
			}

			if msg.Chunk == nil {
				h.sendError(conn, currentSessionID, "No audio chunk provided")
				continue
			}

			err = processor.ProcessChunk(ctx, currentSessionID, *msg.Chunk)
			if err != nil {
				h.sendError(conn, currentSessionID, "Failed to process chunk: "+err.Error())
				continue
			}

			// For real-time mode, we might send partial results here
			if mode == services.RealTimeMode {
				// TODO: Implement real-time partial results
				response := AudioMessage{
					Type:      "chunk_processed",
					SessionID: currentSessionID,
					Message:   "Audio chunk processed",
				}
				if err := conn.WriteJSON(response); err != nil {
					log.Printf("Failed to send chunk_processed response: %v", err)
				}
			}

		case "end_session":
			if currentSessionID == "" {
				h.sendError(conn, "", "No active session to end")
				continue
			}

			result, err := processor.EndSession(ctx, currentSessionID)
			if err != nil {
				h.sendError(conn, currentSessionID, "Failed to end session: "+err.Error())
				continue
			}

			response := AudioMessage{
				Type:      "session_ended",
				SessionID: currentSessionID,
				Result:    result,
				Message:   "Session completed successfully",
			}
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Failed to send session_ended response: %v", err)
			}

			currentSessionID = "" // Reset session

		case "get_status":
			if currentSessionID == "" {
				h.sendError(conn, "", "No active session")
				continue
			}

			status, err := processor.GetSessionStatus(ctx, currentSessionID)
			if err != nil {
				h.sendError(conn, currentSessionID, "Failed to get status: "+err.Error())
				continue
			}

			response := AudioMessage{
				Type:      "status",
				SessionID: currentSessionID,
				Result:    status,
			}
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Failed to send status response: %v", err)
			}

		default:
			h.sendError(conn, currentSessionID, "Unknown message type: "+msg.Type)
		}
	}

	// Clean up session if still active
	if currentSessionID != "" {
		if err := processor.AbortSession(ctx, currentSessionID); err != nil {
			log.Printf("Failed to abort session %s: %v", currentSessionID, err)
		}
	}

	log.Printf("WebSocket connection closed")
}

// handleStartSession handles session initialization
func (h *AudioHandlers) handleStartSession(ctx context.Context, processor services.AudioProcessor, msg AudioMessage, options services.AudioProcessingOptions) (string, error) {
	if msg.Metadata == nil {
		return "", fmt.Errorf("no metadata provided")
	}

	// Set processing mode in metadata
	msg.Metadata.Mode = options.Mode
	msg.Metadata.StartTime = time.Now().Unix()

	// Generate session ID if not provided
	if msg.Metadata.SessionID == "" {
		msg.Metadata.SessionID = fmt.Sprintf("ws_session_%d", time.Now().UnixNano())
	}

	// Override options if provided in message
	if msg.Options != nil {
		if msg.Options.Language != "" {
			options.Language = msg.Options.Language
		}
		if msg.Options.Provider != "" {
			options.Provider = msg.Options.Provider
		}
		options.SpeakerDiarization = msg.Options.SpeakerDiarization
		options.PunctuationFiltering = msg.Options.PunctuationFiltering
		options.ProfanityFiltering = msg.Options.ProfanityFiltering
		if msg.Options.ConfidenceThreshold > 0 {
			options.ConfidenceThreshold = msg.Options.ConfidenceThreshold
		}
	}

	return processor.StartSession(ctx, *msg.Metadata, options)
}

// sendError sends an error message via websocket
func (h *AudioHandlers) sendError(conn *websocket.Conn, sessionID, errorMsg string) {
	response := AudioMessage{
		Type:      "error",
		SessionID: sessionID,
		Error:     errorMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Failed to send error message: %v", err)
	}
}

// GetProviderCapabilities returns the capabilities of transcription providers
// @Summary Get provider capabilities
// @Description Get the capabilities and pricing of available transcription providers
// @Tags audio
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /audio/providers [get]
func (h *AudioHandlers) GetProviderCapabilities(c *gin.Context) {
	providers := h.processorFactory.GetAvailableProviders()
	capabilities := make(map[string]*services.ProviderCapabilities)

	for _, provider := range providers {
		caps, err := h.processorFactory.GetProviderCapabilities(provider)
		if err != nil {
			log.Printf("Failed to get capabilities for provider %s: %v", provider, err)
			continue
		}
		capabilities[provider] = caps
	}

	c.JSON(http.StatusOK, gin.H{
		"providers":    providers,
		"capabilities": capabilities,
	})
}
