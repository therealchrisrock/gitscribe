package providers

import (
	"context"
	"io"
	"testing"
	"time"

	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFirebaseUploader for testing
type MockFirebaseUploader struct {
	mock.Mock
}

func (m *MockFirebaseUploader) UploadAudio(ctx context.Context, audioData []byte, meetingID, sessionID string) (string, error) {
	args := m.Called(ctx, audioData, meetingID, sessionID)
	return args.String(0), args.Error(1)
}

func (m *MockFirebaseUploader) UploadAudioStream(ctx context.Context, reader io.Reader, meetingID, sessionID string) (string, error) {
	args := m.Called(ctx, reader, meetingID, sessionID)
	return args.String(0), args.Error(1)
}

func TestMockAssemblyAIProvider_StartSession(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	metadata := services.AudioStreamMetadata{
		SessionID: "test-session-123",
		MeetingID: "meeting-456",
		UserID:    "user-789",
	}

	options := services.AudioProcessingOptions{
		Mode:               services.BatchMode,
		Provider:           "mock",
		SpeakerDiarization: true,
		Language:           "en",
	}

	ctx := context.Background()
	sessionID, err := provider.StartSession(ctx, metadata, options)

	assert.NoError(t, err)
	assert.Equal(t, "test-session-123", sessionID)
	assert.True(t, provider.IsSessionActive(ctx, sessionID))
}

func TestMockAssemblyAIProvider_ProcessChunks(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	// Start session first
	metadata := services.AudioStreamMetadata{
		SessionID: "test-session-456",
		MeetingID: "meeting-789",
	}
	options := services.AudioProcessingOptions{
		Mode:               services.BatchMode,
		SpeakerDiarization: true,
	}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	// Process audio chunks
	chunk1 := services.AudioChunk{
		Data:        []byte("fake-audio-data-1"),
		Timestamp:   time.Now().Unix(),
		SequenceNum: 1,
	}

	chunk2 := services.AudioChunk{
		Data:        []byte("fake-audio-data-2"),
		Timestamp:   time.Now().Unix(),
		SequenceNum: 2,
	}

	err1 := provider.ProcessChunk(ctx, sessionID, chunk1)
	err2 := provider.ProcessChunk(ctx, sessionID, chunk2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestMockAssemblyAIProvider_EndSession_WithSpeakerDiarization(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	// Mock the Firebase upload
	mockUploader.On("UploadAudio", mock.Anything, mock.Anything, "meeting-789", mock.Anything).
		Return("gs://test-bucket/meeting-789/audio/session_123.wav", nil)

	// Start session with speaker diarization enabled
	metadata := services.AudioStreamMetadata{
		SessionID: "test-session-789",
		MeetingID: "meeting-789",
	}
	options := services.AudioProcessingOptions{
		Mode:               services.BatchMode,
		SpeakerDiarization: true,
		Provider:           "mock",
	}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	// Add some audio chunks
	for i := 1; i <= 3; i++ {
		chunk := services.AudioChunk{
			Data:        []byte("fake-audio-data"),
			Timestamp:   time.Now().Unix(),
			SequenceNum: i,
		}
		provider.ProcessChunk(ctx, sessionID, chunk)
	}

	// End session and get results
	result, err := provider.EndSession(ctx, sessionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities.Completed, result.Status)
	assert.Equal(t, services.BatchMode, result.ProcessingMode)
	assert.True(t, len(result.Segments) > 0)

	// Verify speaker diarization worked
	foundDifferentSpeakers := false
	speakers := make(map[string]bool)
	for _, segment := range result.Segments {
		speakers[segment.Speaker] = true
		if len(speakers) > 1 {
			foundDifferentSpeakers = true
			break
		}
	}
	assert.True(t, foundDifferentSpeakers, "Expected multiple speakers in diarized transcript")

	// Verify Firebase upload was called
	mockUploader.AssertExpectations(t)

	// Verify session is no longer active
	assert.False(t, provider.IsSessionActive(ctx, sessionID))
}

func TestMockAssemblyAIProvider_EndSession_WithoutSpeakerDiarization(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	mockUploader.On("UploadAudio", mock.Anything, mock.Anything, "meeting-123", mock.Anything).
		Return("gs://test-bucket/meeting-123/audio/session_456.wav", nil)

	// Start session WITHOUT speaker diarization
	metadata := services.AudioStreamMetadata{
		SessionID: "test-session-456",
		MeetingID: "meeting-123",
	}
	options := services.AudioProcessingOptions{
		Mode:               services.BatchMode,
		SpeakerDiarization: false, // Disabled
		Provider:           "mock",
	}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	// Add audio chunks
	chunk := services.AudioChunk{
		Data:        []byte("fake-audio-data"),
		Timestamp:   time.Now().Unix(),
		SequenceNum: 1,
	}
	provider.ProcessChunk(ctx, sessionID, chunk)

	// End session
	result, err := provider.EndSession(ctx, sessionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify that without diarization, speakers are marked as unknown
	for _, segment := range result.Segments {
		assert.Equal(t, "Speaker Unknown", segment.Speaker)
	}

	mockUploader.AssertExpectations(t)
}

func TestMockAssemblyAIProvider_GetSessionStatus(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	metadata := services.AudioStreamMetadata{
		SessionID: "status-test-session",
		MeetingID: "meeting-status",
	}
	options := services.AudioProcessingOptions{
		Mode: services.RealTimeMode,
	}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	// Get status before processing
	status1, err1 := provider.GetSessionStatus(ctx, sessionID)
	assert.NoError(t, err1)
	assert.Equal(t, entities.Pending, status1.Status)

	// Process a chunk
	chunk := services.AudioChunk{
		Data:        []byte("audio-data"),
		Timestamp:   time.Now().Unix(),
		SequenceNum: 1,
	}
	provider.ProcessChunk(ctx, sessionID, chunk)

	// Get status after processing
	status2, err2 := provider.GetSessionStatus(ctx, sessionID)
	assert.NoError(t, err2)
	assert.Equal(t, entities.Processing, status2.Status)
}

func TestMockAssemblyAIProvider_AbortSession(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	metadata := services.AudioStreamMetadata{
		SessionID: "abort-test-session",
		MeetingID: "meeting-abort",
	}
	options := services.AudioProcessingOptions{
		Mode: services.BatchMode,
	}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	// Verify session is active
	assert.True(t, provider.IsSessionActive(ctx, sessionID))

	// Abort session
	err := provider.AbortSession(ctx, sessionID)
	assert.NoError(t, err)

	// Verify session is no longer active
	assert.False(t, provider.IsSessionActive(ctx, sessionID))

	// Verify we can't get status of aborted session
	_, err = provider.GetSessionStatus(ctx, sessionID)
	assert.Error(t, err)
}

func TestMockAssemblyAIProvider_GetSupportedModes(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	modes := provider.GetSupportedModes()

	assert.Contains(t, modes, services.RealTimeMode)
	assert.Contains(t, modes, services.BatchMode)
	assert.Len(t, modes, 2)
}

func TestMockAssemblyAIProvider_ErrorCases(t *testing.T) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	ctx := context.Background()

	// Test processing chunk for non-existent session
	chunk := services.AudioChunk{Data: []byte("data")}
	err := provider.ProcessChunk(ctx, "non-existent", chunk)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session non-existent not found")

	// Test ending non-existent session
	_, err = provider.EndSession(ctx, "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session non-existent not found")

	// Test getting status of non-existent session
	_, err = provider.GetSessionStatus(ctx, "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session non-existent not found")

	// Test aborting non-existent session
	err = provider.AbortSession(ctx, "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session non-existent not found")
}

// Benchmark test for chunk processing
func BenchmarkMockAssemblyAIProvider_ProcessChunk(b *testing.B) {
	mockUploader := new(MockFirebaseUploader)
	provider := NewMockAssemblyAIProvider(mockUploader)

	metadata := services.AudioStreamMetadata{
		SessionID: "benchmark-session",
		MeetingID: "benchmark-meeting",
	}
	options := services.AudioProcessingOptions{Mode: services.BatchMode}

	ctx := context.Background()
	sessionID, _ := provider.StartSession(ctx, metadata, options)

	chunk := services.AudioChunk{
		Data:        make([]byte, 1024), // 1KB chunk
		Timestamp:   time.Now().Unix(),
		SequenceNum: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chunk.SequenceNum = i + 1
		provider.ProcessChunk(ctx, sessionID, chunk)
	}
}
