package persistent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"teammate/server/modules/transcription/application/services"
	domainServices "teammate/server/modules/transcription/domain/services"
	"teammate/server/seedwork/infrastructure/events"

	"github.com/gorilla/websocket"
)

// AudioMessage represents the structure of messages sent over websocket
type AudioMessage struct {
	Type      string                                 `json:"type"`
	SessionID string                                 `json:"session_id,omitempty"`
	Metadata  *domainServices.AudioStreamMetadata    `json:"metadata,omitempty"`
	Options   *domainServices.AudioProcessingOptions `json:"options,omitempty"`
	Chunk     *domainServices.AudioChunk             `json:"chunk,omitempty"`
	Error     string                                 `json:"error,omitempty"`
	Result    *domainServices.AudioProcessingResult  `json:"result,omitempty"`
	Message   string                                 `json:"message,omitempty"`
}

// PersistentAudioHandler provides database-integrated WebSocket audio processing
type PersistentAudioHandler struct {
	transcriptionService *services.EnhancedTranscriptionService
	upgrader             websocket.Upgrader
	activeSessions       map[string]*PersistentAudioSession
	eventBus             events.EventBus
}

// PersistentAudioSession represents an active WebSocket session with database persistence
type PersistentAudioSession struct {
	ID                   string
	Conn                 *websocket.Conn
	TranscriptionSession *services.TranscriptionSession
	MeetingID            string
	UserID               string
	CreatedAt            time.Time
	LastActivity         time.Time
	IsActive             bool
	// Buffer for chunks received before session starts
	ChunkBuffer []AudioMessage
	BufferMutex sync.Mutex
}

// NewPersistentAudioHandler creates a new persistent audio handler
func NewPersistentAudioHandler(transcriptionService *services.EnhancedTranscriptionService, eventBus events.EventBus) *PersistentAudioHandler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	handler := &PersistentAudioHandler{
		transcriptionService: transcriptionService,
		upgrader:             upgrader,
		activeSessions:       make(map[string]*PersistentAudioSession),
		eventBus:             eventBus,
	}

	// Subscribe to transcription events for real-time updates
	handler.subscribeToEvents()

	return handler
}

// HandleWebSocketConnection handles enhanced WebSocket connections for audio processing
func (h *PersistentAudioHandler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Parse query parameters
	queryParams := h.parseQueryParams(r.URL.Query())
	log.Printf("Enhanced WebSocket connection established with params: %+v", queryParams)

	// Create session
	session := &PersistentAudioSession{
		ID:           fmt.Sprintf("enhanced_session_%d", time.Now().UnixNano()),
		Conn:         conn,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		IsActive:     true,
		ChunkBuffer:  make([]AudioMessage, 0),
	}

	// Add to active sessions
	h.activeSessions[session.ID] = session

	// Setup connection close cleanup
	defer func() {
		session.IsActive = false
		if session.TranscriptionSession != nil {
			// End transcription session on disconnect
			h.transcriptionService.EndTranscriptionSession(context.Background(), session.TranscriptionSession)
		}
		// Clear any buffered chunks
		session.BufferMutex.Lock()
		session.ChunkBuffer = nil
		session.BufferMutex.Unlock()

		delete(h.activeSessions, session.ID)
		log.Printf("Enhanced WebSocket session %s cleaned up", session.ID)
	}()

	// Send welcome message
	h.sendMessage(conn, AudioMessage{
		Type:    "connection_established",
		Message: fmt.Sprintf("Enhanced WebSocket connection established. Session ID: %s", session.ID),
	})

	// Handle messages
	h.handleMessages(session, queryParams)
}

// handleMessages processes incoming WebSocket messages
func (h *PersistentAudioHandler) handleMessages(session *PersistentAudioSession, params PersistentQueryParams) {
	for {
		var msg AudioMessage // We read as AudioMessage but send as PersistentAudioMessage
		err := session.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Enhanced WebSocket error: %v", err)
			}
			break
		}

		session.LastActivity = time.Now()

		log.Printf("Enhanced Handler: Received message type: %s for session: %s", msg.Type, session.ID)

		switch msg.Type {
		case "start_session":
			log.Printf("Enhanced Handler: Processing start_session message for session: %s", session.ID)
			h.handleStartSession(session, msg, params)
		case "audio_chunk":
			h.handleAudioChunk(session, msg)
		case "end_session":
			log.Printf("Enhanced Handler: Processing end_session message for session: %s", session.ID)
			h.handleEndSession(session, msg)
		case "get_session_status":
			h.handleGetSessionStatus(session, msg)
		case "get_transcription_history":
			h.handleGetTranscriptionHistory(session, msg)
		case "get_transcription_stats":
			h.handleGetTranscriptionStats(session, msg)
		default:
			log.Printf("Enhanced Handler: Unknown message type: %s for session: %s", msg.Type, session.ID)
			h.sendMessage(session.Conn, AudioMessage{
				Type:  "error",
				Error: fmt.Sprintf("Unknown message type: %s", msg.Type),
			})
		}
	}
}

// handleStartSession starts a new transcription session with database persistence
func (h *PersistentAudioHandler) handleStartSession(session *PersistentAudioSession, msg AudioMessage, params PersistentQueryParams) {
	log.Printf("Enhanced Handler: handleStartSession called for session: %s", session.ID)

	if session.TranscriptionSession != nil {
		log.Printf("Enhanced Handler: Session already started for session: %s", session.ID)
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "Session already started",
		})
		return
	}

	log.Printf("Enhanced Handler: Parsing metadata for session: %s, metadata: %+v", session.ID, msg.Metadata)

	// Parse metadata
	metadata, err := h.parseMetadata(msg.Metadata)
	if err != nil {
		log.Printf("Enhanced Handler: Failed to parse metadata: %v", err)
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Invalid metadata: %v", err),
		})
		return
	}

	log.Printf("Enhanced Handler: Parsed metadata successfully: %+v", metadata)

	// Parse options
	options, err := h.parseOptions(msg.Options, params)
	if err != nil {
		log.Printf("Enhanced Handler: Failed to parse options: %v", err)
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Invalid options: %v", err),
		})
		return
	}

	log.Printf("Enhanced Handler: Parsed options successfully: %+v", options)

	// Create meeting if needed (simplified - in production this would be more sophisticated)
	// TODO: DDD Violation - Direct meeting access should be through commands/queries
	// For now, we'll simplify by assuming meetings exist or using simplified meeting creation
	log.Printf("Enhanced Handler: Creating transcription for meeting: %s", metadata.MeetingID)

	// Create transcription request
	req := &services.StartTranscriptionRequest{
		MeetingID:        metadata.MeetingID,
		Metadata:         metadata,
		Options:          options,
		CreateBotSession: false, // This is direct audio streaming, no bot needed
	}

	// Start transcription session
	transcriptionSession, err := h.transcriptionService.StartTranscriptionSession(context.Background(), req)
	if err != nil {
		log.Printf("Enhanced Handler: Failed to start transcription session: %v", err)
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to start transcription session: %v", err),
		})
		return
	}

	log.Printf("Enhanced Handler: Transcription session started successfully: %s", transcriptionSession.SessionID)
	session.TranscriptionSession = transcriptionSession
	session.MeetingID = metadata.MeetingID
	session.UserID = metadata.UserID

	// Process any buffered chunks
	h.processBufferedChunks(session)

	h.sendMessage(session.Conn, AudioMessage{
		Type:      "session_started",
		SessionID: transcriptionSession.SessionID,
		Message:   "Enhanced transcription session started successfully",
	})
}

// handleAudioChunk processes audio chunks with database updates
func (h *PersistentAudioHandler) handleAudioChunk(session *PersistentAudioSession, msg AudioMessage) {
	if session.TranscriptionSession == nil {
		// Buffer the chunk instead of rejecting it
		session.BufferMutex.Lock()
		session.ChunkBuffer = append(session.ChunkBuffer, msg)
		session.BufferMutex.Unlock()

		log.Printf("Enhanced Handler: Buffering audio chunk %d (session not started yet)",
			func() int {
				if msg.Chunk != nil {
					return msg.Chunk.SequenceNum
				}
				return -1
			}())

		h.sendMessage(session.Conn, AudioMessage{
			Type:    "chunk_buffered",
			Message: fmt.Sprintf("Audio chunk buffered (session not started yet). Send start_session message to begin processing."),
		})
		return
	}

	if msg.Chunk == nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "Missing audio chunk data",
		})
		return
	}

	// Create audio chunk
	chunk := domainServices.AudioChunk{
		Data:        msg.Chunk.Data,
		Timestamp:   msg.Chunk.Timestamp,
		SequenceNum: msg.Chunk.SequenceNum,
	}

	// Process chunk
	err := h.transcriptionService.ProcessAudioChunk(context.Background(), session.TranscriptionSession, chunk)
	if err != nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to process audio chunk: %v", err),
		})
		return
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "chunk_processed",
		Message: fmt.Sprintf("Audio chunk %d processed successfully", msg.Chunk.SequenceNum),
	})
}

// handleEndSession ends the transcription session and saves results
func (h *PersistentAudioHandler) handleEndSession(session *PersistentAudioSession, msg AudioMessage) {
	if session.TranscriptionSession == nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "No active session to end",
		})
		return
	}

	// End transcription session
	_, err := h.transcriptionService.EndTranscriptionSession(context.Background(), session.TranscriptionSession)
	if err != nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to end transcription session: %v", err),
		})
		return
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "session_ended",
		Message: "Enhanced transcription session ended successfully",
	})

	session.TranscriptionSession = nil
}

// processBufferedChunks processes any audio chunks that were buffered before session started
func (h *PersistentAudioHandler) processBufferedChunks(session *PersistentAudioSession) {
	session.BufferMutex.Lock()
	bufferedChunks := session.ChunkBuffer
	session.ChunkBuffer = nil // Clear the buffer
	session.BufferMutex.Unlock()

	if len(bufferedChunks) == 0 {
		return
	}

	log.Printf("Enhanced Handler: Processing %d buffered audio chunks", len(bufferedChunks))

	for _, chunkMsg := range bufferedChunks {
		if chunkMsg.Chunk == nil {
			continue
		}

		// Create audio chunk
		chunk := domainServices.AudioChunk{
			Data:        chunkMsg.Chunk.Data,
			Timestamp:   chunkMsg.Chunk.Timestamp,
			SequenceNum: chunkMsg.Chunk.SequenceNum,
		}

		// Process chunk
		err := h.transcriptionService.ProcessAudioChunk(context.Background(), session.TranscriptionSession, chunk)
		if err != nil {
			log.Printf("Enhanced Handler: Failed to process buffered chunk %d: %v", chunk.SequenceNum, err)
			h.sendMessage(session.Conn, AudioMessage{
				Type:  "error",
				Error: fmt.Sprintf("Failed to process buffered chunk %d: %v", chunk.SequenceNum, err),
			})
		} else {
			log.Printf("Enhanced Handler: Successfully processed buffered chunk %d", chunk.SequenceNum)
		}
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "buffered_chunks_processed",
		Message: fmt.Sprintf("Processed %d buffered audio chunks", len(bufferedChunks)),
	})
}

// handleGetSessionStatus retrieves the current session status
func (h *PersistentAudioHandler) handleGetSessionStatus(session *PersistentAudioSession, msg AudioMessage) {
	if session.TranscriptionSession == nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "No active session",
		})
		return
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "session_status",
		Message: "Session status retrieved",
	})
}

// handleGetTranscriptionHistory retrieves transcription history for a meeting
func (h *PersistentAudioHandler) handleGetTranscriptionHistory(session *PersistentAudioSession, msg AudioMessage) {
	meetingID := session.MeetingID
	if meetingID == "" {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "No meeting ID available",
		})
		return
	}

	_, err := h.transcriptionService.GetTranscriptionHistory(context.Background(), meetingID)
	if err != nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to get transcription history: %v", err),
		})
		return
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "transcription_history",
		Message: "Transcription history retrieved",
	})
}

// handleGetTranscriptionStats retrieves analytics for the meeting's transcriptions
func (h *PersistentAudioHandler) handleGetTranscriptionStats(session *PersistentAudioSession, msg AudioMessage) {
	meetingID := session.MeetingID
	if meetingID == "" {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: "No meeting ID available",
		})
		return
	}

	_, err := h.transcriptionService.GetTranscriptionStats(context.Background(), meetingID)
	if err != nil {
		h.sendMessage(session.Conn, AudioMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to get transcription stats: %v", err),
		})
		return
	}

	h.sendMessage(session.Conn, AudioMessage{
		Type:    "transcription_stats",
		Message: "Transcription statistics retrieved",
	})
}

// subscribeToEvents sets up event subscriptions for real-time updates
func (h *PersistentAudioHandler) subscribeToEvents() {
	// Subscribe to transcription events
	h.eventBus.Subscribe("transcription.session.started", func(event interface{}) {
		if startedEvent, ok := event.(*services.TranscriptionSessionStartedEvent); ok {
			h.broadcastToMeeting(startedEvent.MeetingID, AudioMessage{
				Type:    "transcription_started",
				Message: "Transcription session started",
			})
		}
	})

	h.eventBus.Subscribe("transcription.processing", func(event interface{}) {
		if processingEvent, ok := event.(*services.TranscriptionProcessingEvent); ok {
			h.broadcastToMeeting(processingEvent.MeetingID, AudioMessage{
				Type:    "transcription_processing",
				Message: "Transcription processing update",
			})
		}
	})

	h.eventBus.Subscribe("transcription.completed", func(event interface{}) {
		if completedEvent, ok := event.(*services.TranscriptionCompletedEvent); ok {
			h.broadcastToMeeting(completedEvent.MeetingID, AudioMessage{
				Type:    "transcription_completed",
				Message: "Transcription completed",
			})
		}
	})
}

// broadcastToMeeting sends a message to all sessions in a meeting
func (h *PersistentAudioHandler) broadcastToMeeting(meetingID string, message AudioMessage) {
	for _, session := range h.activeSessions {
		if session.MeetingID == meetingID && session.IsActive {
			h.sendMessage(session.Conn, message)
		}
	}
}

// Helper methods (reusing from original handler)

func (h *PersistentAudioHandler) parseMetadata(metadata interface{}) (domainServices.AudioStreamMetadata, error) {
	var result domainServices.AudioStreamMetadata

	if metadata == nil {
		return result, fmt.Errorf("metadata is required")
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(jsonData, &result)
	return result, err
}

func (h *PersistentAudioHandler) parseOptions(options interface{}, params PersistentQueryParams) (domainServices.AudioProcessingOptions, error) {
	var result domainServices.AudioProcessingOptions

	// Set defaults from query params
	result.Provider = params.Provider
	result.Mode = params.Mode
	result.SpeakerDiarization = params.SpeakerDiarization
	result.Language = params.Language

	// Override with options from message if provided
	if options != nil {
		jsonData, err := json.Marshal(options)
		if err != nil {
			return result, err
		}

		var msgOptions domainServices.AudioProcessingOptions
		err = json.Unmarshal(jsonData, &msgOptions)
		if err != nil {
			return result, err
		}

		// Merge options
		if msgOptions.Provider != "" {
			result.Provider = msgOptions.Provider
		}
		if msgOptions.Mode != "" {
			result.Mode = msgOptions.Mode
		}
		if msgOptions.Language != "" {
			result.Language = msgOptions.Language
		}
		// Override boolean fields
		result.SpeakerDiarization = msgOptions.SpeakerDiarization
		result.RealTimeTranscription = msgOptions.RealTimeTranscription
		result.CostOptimized = msgOptions.CostOptimized
	}

	return result, nil
}

func (h *PersistentAudioHandler) parseQueryParams(values url.Values) PersistentQueryParams {
	params := PersistentQueryParams{
		Provider:           values.Get("provider"),
		Mode:               domainServices.ProcessingMode(values.Get("mode")),
		Language:           values.Get("language"),
		SpeakerDiarization: false,
	}

	if params.Provider == "" {
		params.Provider = "mock"
	}
	if params.Mode == "" {
		params.Mode = domainServices.BatchMode
	}
	if params.Language == "" {
		params.Language = "en"
	}

	if values.Get("speaker_diarization") == "true" {
		params.SpeakerDiarization = true
	}

	return params
}

func (h *PersistentAudioHandler) sendMessage(conn *websocket.Conn, message AudioMessage) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Printf("Failed to send WebSocket message: %v", err)
	}
}

// GetActiveSessionsCount returns the number of active sessions
func (h *PersistentAudioHandler) GetActiveSessionsCount() int {
	return len(h.activeSessions)
}

// GetSessionsByMeeting returns sessions for a specific meeting
func (h *PersistentAudioHandler) GetSessionsByMeeting(meetingID string) []*PersistentAudioSession {
	var sessions []*PersistentAudioSession
	for _, session := range h.activeSessions {
		if session.MeetingID == meetingID {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// GetProviderCapabilities returns the capabilities of transcription providers
// @Summary Get provider capabilities
// @Description Get the capabilities and pricing of available transcription providers
// @Tags audio
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /audio/providers [get]
func (h *PersistentAudioHandler) GetProviderCapabilities(w http.ResponseWriter, r *http.Request) {
	// Get the audio factory from the transcription service
	// We'll need to access it through the service or inject it separately
	w.Header().Set("Content-Type", "application/json")

	// For now, return a basic response - this would need proper implementation
	// depending on how you want to expose the factory capabilities
	response := map[string]interface{}{
		"providers": []string{"mock", "assemblyai"},
		"message":   "Provider capabilities endpoint - implement based on your factory access pattern",
	}

	json.NewEncoder(w).Encode(response)
}

// PersistentQueryParams represents query parameters for WebSocket connections (for persistent handlers)
type PersistentQueryParams struct {
	Provider           string                        `json:"provider"`
	Mode               domainServices.ProcessingMode `json:"mode"`
	Language           string                        `json:"language"`
	SpeakerDiarization bool                          `json:"speaker_diarization"`
}

// PersistentAudioMessage represents the structure of messages sent over websocket for persistent handlers
// This differs from AudioMessage in that it uses interface{} for Result to allow more flexible data types
type PersistentAudioMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Options   interface{} `json:"options,omitempty"`
	Chunk     *AudioChunk `json:"chunk,omitempty"`
}

// AudioChunk represents an audio data chunk
type AudioChunk struct {
	Data        []byte `json:"data"`
	Timestamp   int64  `json:"timestamp"`
	SequenceNum int    `json:"sequence_num"`
}
