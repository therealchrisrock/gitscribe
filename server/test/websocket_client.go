package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// AudioMessage matches the structure in audio_handlers.go
type AudioMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Options   interface{} `json:"options,omitempty"`
	Chunk     *AudioChunk `json:"chunk,omitempty"`
	Error     string      `json:"error,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Message   string      `json:"message,omitempty"`
}

type AudioChunk struct {
	Data        []byte  `json:"data"`
	Timestamp   int64   `json:"timestamp"`
	SequenceNum int     `json:"sequence_num"`
	Size        int     `json:"size"`
	Duration    float64 `json:"duration,omitempty"`
}

type AudioStreamMetadata struct {
	SessionID     string `json:"session_id"`
	MeetingID     string `json:"meeting_id"`
	UserID        string `json:"user_id"`
	SampleRate    int    `json:"sample_rate"`
	Channels      int    `json:"channels"`
	BitsPerSample int    `json:"bits_per_sample"`
	MimeType      string `json:"mime_type"`
	StartTime     int64  `json:"start_time"`
	Mode          string `json:"mode"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// WebSocket URL with speaker diarization enabled
	u := url.URL{
		Scheme:   "ws",
		Host:     "localhost:8080",
		Path:     "/ws/audio",
		RawQuery: "provider=mock&speaker_diarization=true&mode=batch&language=en",
	}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Read messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var msg AudioMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Received: %s - %s", msg.Type, msg.Message)
			if msg.Error != "" {
				log.Printf("Error: %s", msg.Error)
			}
			if msg.Result != nil {
				log.Printf("Result: %+v", msg.Result)
			}
		}
	}()

	// Test sequence
	go func() {
		// 1. Start session
		log.Println("Starting session...")
		startMsg := AudioMessage{
			Type:      "start_session",
			SessionID: "test-session-123",
			Metadata: AudioStreamMetadata{
				SessionID:     "test-session-123",
				MeetingID:     "meeting-456",
				UserID:        "user-789",
				SampleRate:    44100,
				Channels:      1,
				BitsPerSample: 16,
				MimeType:      "audio/wav",
				StartTime:     time.Now().Unix(),
				Mode:          "batch",
			},
		}

		if err := c.WriteJSON(startMsg); err != nil {
			log.Println("write start_session:", err)
			return
		}

		time.Sleep(1 * time.Second)

		// 2. Send audio chunks
		log.Println("Sending audio chunks...")
		for i := 1; i <= 5; i++ {
			chunk := AudioChunk{
				Data:        []byte(fmt.Sprintf("fake-audio-data-chunk-%d", i)),
				Timestamp:   time.Now().Unix(),
				SequenceNum: i,
				Size:        len(fmt.Sprintf("fake-audio-data-chunk-%d", i)),
				Duration:    1.0, // 1 second per chunk
			}

			chunkMsg := AudioMessage{
				Type:  "audio_chunk",
				Chunk: &chunk,
			}

			if err := c.WriteJSON(chunkMsg); err != nil {
				log.Println("write audio_chunk:", err)
				return
			}

			log.Printf("Sent chunk %d", i)
			time.Sleep(500 * time.Millisecond) // Small delay between chunks
		}

		time.Sleep(1 * time.Second)

		// 3. End session
		log.Println("Ending session...")
		endMsg := AudioMessage{
			Type: "end_session",
		}

		if err := c.WriteJSON(endMsg); err != nil {
			log.Println("write end_session:", err)
			return
		}

		// Wait a bit for the final response
		time.Sleep(5 * time.Second)

		// Close connection
		log.Println("Test completed, closing connection...")
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
