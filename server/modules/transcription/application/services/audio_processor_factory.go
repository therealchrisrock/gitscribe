package services

import (
	"fmt"
	"os"
	"teammate/server/modules/transcription/domain/services"
	"teammate/server/modules/transcription/infrastructure/providers"
)

// AudioProcessorFactory creates audio processors based on requirements
// This concrete implementation satisfies the domain services.AudioProcessorFactory interface
type AudioProcessorFactory struct {
	firebaseUploader *providers.FirebaseStorageUploader
}

// Ensure AudioProcessorFactory implements the domain interface
var _ services.AudioProcessorFactory = (*AudioProcessorFactory)(nil)

// NewAudioProcessorFactory creates a new factory instance
func NewAudioProcessorFactory() *AudioProcessorFactory {
	// Initialize Firebase uploader
	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")

	if bucketName == "" {
		bucketName = "gitscribe-default"
	}
	if credentialsPath == "" {
		credentialsPath = "./firebase-credentials/serviceAccountKey.json"
	}

	firebaseUploader, err := providers.NewFirebaseStorageUploader(bucketName, credentialsPath)
	if err != nil {
		// Log warning but continue with nil uploader for development
		fmt.Printf("Warning: Failed to initialize Firebase storage uploader: %v\n", err)
		firebaseUploader = nil
	}

	return &AudioProcessorFactory{
		firebaseUploader: firebaseUploader,
	}
}

// CreateProcessor creates an audio processor based on the specified mode and options
func (f *AudioProcessorFactory) CreateProcessor(mode services.ProcessingMode, options services.AudioProcessingOptions) (services.AudioProcessor, error) {
	switch options.Provider {
	case "assemblyai":
		return f.createAssemblyAIProcessor(mode, options)
	case "mock":
		return f.createMockProcessor(mode, options)
	default:
		// Default to mock for development
		return f.createMockProcessor(mode, options)
	}
}

// createAssemblyAIProcessor creates a real AssemblyAI processor
func (f *AudioProcessorFactory) createAssemblyAIProcessor(mode services.ProcessingMode, options services.AudioProcessingOptions) (services.AudioProcessor, error) {
	// Check if AssemblyAI API key is available
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		// Fall back to mock provider if no API key is configured
		fmt.Printf("Warning: ASSEMBLYAI_API_KEY not set, using mock provider\n")
		return providers.NewMockAssemblyAIProvider(f.firebaseUploader), nil
	}

	// For production AssemblyAI, we require Firebase uploader
	if f.firebaseUploader == nil {
		// Fall back to mock if Firebase is not available
		fmt.Printf("Warning: Firebase not available, falling back to mock provider\n")
		return providers.NewMockAssemblyAIProvider(f.firebaseUploader), nil
	}

	// Use real AssemblyAI provider
	return providers.NewAssemblyAIProvider(apiKey, f.firebaseUploader), nil
}

// createMockProcessor creates a mock processor for development/testing
func (f *AudioProcessorFactory) createMockProcessor(mode services.ProcessingMode, options services.AudioProcessingOptions) (services.AudioProcessor, error) {
	// Mock provider doesn't require Firebase uploader - it can work with nil
	return providers.NewMockAssemblyAIProvider(f.firebaseUploader), nil
}

// GetAvailableProviders returns a list of available transcription providers
func (f *AudioProcessorFactory) GetAvailableProviders() []string {
	return []string{
		"assemblyai",
		"mock", // For development/testing
	}
}

// GetProviderCapabilities returns the capabilities of a specific provider
func (f *AudioProcessorFactory) GetProviderCapabilities(provider string) (*services.ProviderCapabilities, error) {
	switch provider {
	case "assemblyai":
		return &services.ProviderCapabilities{
			Provider:            "assemblyai",
			SupportedModes:      []services.ProcessingMode{services.RealTimeMode, services.BatchMode},
			SupportedLanguages:  []string{"en", "es", "fr", "de", "it", "pt", "hi", "ja", "ko", "zh"},
			SupportsDiarization: true,
			SupportsRealTime:    true,
			SupportsBatch:       true,
			MaxAudioDuration:    7200, // 2 hours
			SupportedFormats:    []string{"wav", "mp3", "m4a", "flac", "opus"},
			PricingPerMinute: map[string]float64{
				"realtime": 0.0025,  // $0.0025/minute for real-time
				"batch":    0.00065, // $0.00065/minute for batch
			},
			EstimatedLatency: map[string]int{
				"realtime": 1,   // ~1 second
				"batch":    300, // ~5 minutes
			},
		}, nil
	case "mock":
		return &services.ProviderCapabilities{
			Provider:            "mock",
			SupportedModes:      []services.ProcessingMode{services.RealTimeMode, services.BatchMode},
			SupportedLanguages:  []string{"en"},
			SupportsDiarization: true,
			SupportsRealTime:    true,
			SupportsBatch:       true,
			MaxAudioDuration:    3600,
			SupportedFormats:    []string{"wav", "mp3"},
			PricingPerMinute: map[string]float64{
				"realtime": 0.001,
				"batch":    0.0005,
			},
			EstimatedLatency: map[string]int{
				"realtime": 1,
				"batch":    30,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

// RecommendProvider recommends the best provider based on requirements
func (f *AudioProcessorFactory) RecommendProvider(options services.AudioProcessingOptions) (string, error) {
	providers := f.GetAvailableProviders()

	// Simple recommendation logic based on requirements
	if options.CostOptimized {
		return "mock", nil // Cheapest for development
	}

	if options.RealTimeTranscription {
		return "assemblyai", nil // Best real-time support
	}

	if options.SpeakerDiarization {
		return "assemblyai", nil // Excellent diarization
	}

	// Default recommendation
	if len(providers) > 0 {
		return "assemblyai", nil // Default to AssemblyAI
	}

	return "", fmt.Errorf("no suitable provider found")
}
