package handlers

import (
	"log"
	"net/http"
	"time"

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
	// TODO: Add transcription service when implemented
}

func NewAudioHandlers() *AudioHandlers {
	return &AudioHandlers{}
}

// HandleAudioWebSocket handles incoming websocket connections for audio streaming
// @Summary Handle audio websocket connection
// @Description Handle websocket connection for real-time audio streaming and processing
// @Tags transcription
// @Accept json
// @Produce json
// @Router /ws/audio [get]
func (h *AudioHandlers) HandleAudioWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	log.Printf("New audio WebSocket connection established from %s", conn.RemoteAddr())

	// Set read deadline for the connection
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Set ping/pong handlers for connection health
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	chunkCount := 0
	totalBytes := 0

	// Read audio chunks from the WebSocket connection
	for {
		// Read message from client
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			} else {
				log.Printf("WebSocket connection closed: %v", err)
			}
			break
		}

		// Only process binary messages (audio data)
		if messageType == websocket.BinaryMessage {
			chunkCount++
			totalBytes += len(data)

			log.Printf("Received audio chunk #%d: %d bytes (total: %d bytes)",
				chunkCount, len(data), totalBytes)

			// TODO: Process the audio chunk here
			// This is where you would:
			// 1. Send to transcription service (AssemblyAI, Whisper, etc.)
			// 2. Store chunks for later processing
			// 3. Real-time processing if needed

			// For now, just log the receipt
			// Example of what might happen:
			// err := h.transcriptionService.ProcessAudioChunk(data, timestamp)
			// if err != nil {
			//     log.Printf("Error processing audio chunk: %v", err)
			//     conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to process audio chunk"}`))
			//     continue
			// }

			// Reset read deadline on successful message
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}
	}

	log.Printf("Audio WebSocket session ended. Processed %d chunks (%d total bytes)", chunkCount, totalBytes)
}
