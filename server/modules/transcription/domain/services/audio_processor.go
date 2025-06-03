package services

import (
	"context"
	"teammate/server/modules/transcription/domain/entities"
)

// ProcessingMode defines whether audio should be processed in real-time or batch mode
type ProcessingMode string

const (
	RealTimeMode ProcessingMode = "realtime"
	BatchMode    ProcessingMode = "batch"
)

// AudioStreamMetadata contains metadata about an audio stream session
type AudioStreamMetadata struct {
	SessionID     string         `json:"session_id"`
	MeetingID     string         `json:"meeting_id,omitempty"`
	UserID        string         `json:"user_id,omitempty"`
	SampleRate    int            `json:"sample_rate,omitempty"`
	Channels      int            `json:"channels,omitempty"`
	BitsPerSample int            `json:"bits_per_sample,omitempty"`
	MimeType      string         `json:"mime_type,omitempty"`
	StartTime     int64          `json:"start_time"`
	Mode          ProcessingMode `json:"mode"` // New field to specify processing mode
}

// AudioChunk represents a chunk of audio data with timing information
type AudioChunk struct {
	Data        []byte  `json:"data"`
	Timestamp   int64   `json:"timestamp"`
	SequenceNum int     `json:"sequence_num"`
	Size        int     `json:"size"`
	Duration    float64 `json:"duration,omitempty"` // Duration in seconds
}

// AudioProcessingResult contains the result of processing an audio stream
type AudioProcessingResult struct {
	TranscriptionID string                       `json:"transcription_id"`
	Status          entities.TranscriptionStatus `json:"status"`
	Segments        []entities.TranscriptSegment `json:"segments,omitempty"`
	Error           error                        `json:"error,omitempty"`
	Message         string                       `json:"message,omitempty"`
	ProcessingMode  ProcessingMode               `json:"processing_mode"`        // Indicates which mode was used
	FirebaseURL     string                       `json:"firebase_url,omitempty"` // URL to the uploaded audio file
}

// AudioProcessingOptions contains configuration options for audio processing
type AudioProcessingOptions struct {
	Language              string         `json:"language,omitempty"`                // Language code (e.g., "en-US")
	SpeakerDiarization    bool           `json:"speaker_diarization,omitempty"`     // Enable speaker identification
	PunctuationFiltering  bool           `json:"punctuation_filtering,omitempty"`   // Enable automatic punctuation
	ProfanityFiltering    bool           `json:"profanity_filtering,omitempty"`     // Enable profanity filtering
	ConfidenceThreshold   float64        `json:"confidence_threshold,omitempty"`    // Minimum confidence threshold
	RealTimeTranscription bool           `json:"real_time_transcription,omitempty"` // Enable real-time processing
	Provider              string         `json:"provider,omitempty"`                // Transcription provider (assemblyai, whisper, etc.)
	Mode                  ProcessingMode `json:"mode,omitempty"`                    // Processing mode (realtime or batch)

	// Batch-specific options
	BatchPriority string `json:"batch_priority,omitempty"` // Priority level for batch processing (low, normal, high)
	MaxLatency    int    `json:"max_latency,omitempty"`    // Maximum acceptable latency in seconds for batch mode
	CostOptimized bool   `json:"cost_optimized,omitempty"` // Whether to prioritize cost over speed
	QualityLevel  string `json:"quality_level,omitempty"`  // Quality level (basic, standard, premium)
}

// AudioProcessor defines the contract for processing audio streams (both real-time and batch)
type AudioProcessor interface {
	// StartSession initializes a new audio processing session
	// Returns a session ID that should be used for all subsequent calls
	StartSession(ctx context.Context, metadata AudioStreamMetadata, options AudioProcessingOptions) (string, error)

	// ProcessChunk processes an individual audio chunk
	// Should be called for each chunk received during the session
	ProcessChunk(ctx context.Context, sessionID string, chunk AudioChunk) error

	// EndSession finalizes the audio processing session
	// Returns the final transcription result
	EndSession(ctx context.Context, sessionID string) (*AudioProcessingResult, error)

	// GetSessionStatus returns the current status of a session
	GetSessionStatus(ctx context.Context, sessionID string) (*AudioProcessingResult, error)

	// AbortSession cancels an ongoing session and cleans up resources
	AbortSession(ctx context.Context, sessionID string) error

	// IsSessionActive checks if a session is currently active
	IsSessionActive(ctx context.Context, sessionID string) bool

	// GetSupportedModes returns the processing modes supported by this processor
	GetSupportedModes() []ProcessingMode
}

// BatchAudioProcessor extends AudioProcessor for batch/post-processing capabilities
type BatchAudioProcessor interface {
	AudioProcessor

	// SubmitForBatchProcessing submits complete audio for batch processing
	// This is an alternative to the chunk-by-chunk approach for post-processing
	SubmitForBatchProcessing(ctx context.Context, audioData []byte, metadata AudioStreamMetadata, options AudioProcessingOptions) (string, error)

	// GetBatchJob returns information about a batch job
	GetBatchJob(ctx context.Context, jobID string) (*BatchJobInfo, error)

	// ListBatchJobs returns a list of batch jobs with optional filtering
	ListBatchJobs(ctx context.Context, filter BatchJobFilter) ([]BatchJobInfo, error)

	// CancelBatchJob cancels a pending batch job
	CancelBatchJob(ctx context.Context, jobID string) error
}

// BatchJobInfo contains information about a batch processing job
type BatchJobInfo struct {
	JobID         string                       `json:"job_id"`
	SessionID     string                       `json:"session_id,omitempty"`
	Status        entities.TranscriptionStatus `json:"status"`
	SubmittedAt   int64                        `json:"submitted_at"`
	StartedAt     int64                        `json:"started_at,omitempty"`
	CompletedAt   int64                        `json:"completed_at,omitempty"`
	AudioDuration float64                      `json:"audio_duration"` // Duration in seconds
	AudioSize     int                          `json:"audio_size"`     // Size in bytes
	EstimatedCost float64                      `json:"estimated_cost,omitempty"`
	ActualCost    float64                      `json:"actual_cost,omitempty"`
	Provider      string                       `json:"provider"`
	Options       AudioProcessingOptions       `json:"options"`
	Result        *AudioProcessingResult       `json:"result,omitempty"`
	ErrorMessage  string                       `json:"error_message,omitempty"`
}

// BatchJobFilter for filtering batch jobs
type BatchJobFilter struct {
	Status          entities.TranscriptionStatus `json:"status,omitempty"`
	Provider        string                       `json:"provider,omitempty"`
	SubmittedAfter  int64                        `json:"submitted_after,omitempty"`
	SubmittedBefore int64                        `json:"submitted_before,omitempty"`
	UserID          string                       `json:"user_id,omitempty"`
	Limit           int                          `json:"limit,omitempty"`
	Offset          int                          `json:"offset,omitempty"`
}

// RealTimeAudioProcessor extends AudioProcessor for real-time transcription capabilities
type RealTimeAudioProcessor interface {
	AudioProcessor

	// StartRealTimeTranscription enables real-time transcription with callback
	StartRealTimeTranscription(ctx context.Context, sessionID string, callback RealTimeTranscriptionCallback) error

	// StopRealTimeTranscription disables real-time transcription
	StopRealTimeTranscription(ctx context.Context, sessionID string) error
}

// RealTimeTranscriptionCallback is called when partial transcription results are available
type RealTimeTranscriptionCallback func(sessionID string, partialResult entities.TranscriptSegment) error

// AudioSessionEventType represents different types of events during audio processing
type AudioSessionEventType string

const (
	SessionStarted            AudioSessionEventType = "session_started"
	ChunkProcessed            AudioSessionEventType = "chunk_processed"
	PartialTranscriptionReady AudioSessionEventType = "partial_transcription_ready"
	TranscriptionCompleted    AudioSessionEventType = "transcription_completed"
	SessionEnded              AudioSessionEventType = "session_ended"
	SessionError              AudioSessionEventType = "session_error"
	BatchJobSubmitted         AudioSessionEventType = "batch_job_submitted"
	BatchJobStarted           AudioSessionEventType = "batch_job_started"
	BatchJobCompleted         AudioSessionEventType = "batch_job_completed"
)

// AudioSessionEvent represents an event that occurs during audio processing
type AudioSessionEvent struct {
	Type      AudioSessionEventType `json:"type"`
	SessionID string                `json:"session_id"`
	JobID     string                `json:"job_id,omitempty"` // For batch jobs
	Timestamp int64                 `json:"timestamp"`
	Data      interface{}           `json:"data,omitempty"`
	Error     error                 `json:"error,omitempty"`
}

// AudioProcessorEventHandler handles events from the audio processor
type AudioProcessorEventHandler interface {
	HandleEvent(ctx context.Context, event AudioSessionEvent) error
}

// AudioProcessorFactory creates appropriate audio processors based on requirements
type AudioProcessorFactory interface {
	// CreateProcessor creates an audio processor based on the specified mode and options
	CreateProcessor(mode ProcessingMode, options AudioProcessingOptions) (AudioProcessor, error)

	// GetAvailableProviders returns a list of available transcription providers
	GetAvailableProviders() []string

	// GetProviderCapabilities returns the capabilities of a specific provider
	GetProviderCapabilities(provider string) (*ProviderCapabilities, error)
}

// ProviderCapabilities describes what a transcription provider supports
type ProviderCapabilities struct {
	Provider            string             `json:"provider"`
	SupportedModes      []ProcessingMode   `json:"supported_modes"`
	SupportedLanguages  []string           `json:"supported_languages"`
	SupportsDiarization bool               `json:"supports_diarization"`
	SupportsRealTime    bool               `json:"supports_real_time"`
	SupportsBatch       bool               `json:"supports_batch"`
	MaxAudioDuration    int                `json:"max_audio_duration"` // seconds
	SupportedFormats    []string           `json:"supported_formats"`
	PricingPerMinute    map[string]float64 `json:"pricing_per_minute"` // pricing by mode
	EstimatedLatency    map[string]int     `json:"estimated_latency"`  // latency by mode in seconds
}
