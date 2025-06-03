package services

import (
	"os"
	"testing"

	"teammate/server/modules/transcription/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestNewAudioProcessorFactory(t *testing.T) {
	factory := NewAudioProcessorFactory()
	assert.NotNil(t, factory)
}

func TestAudioProcessorFactory_GetAvailableProviders(t *testing.T) {
	factory := NewAudioProcessorFactory()
	providers := factory.GetAvailableProviders()

	assert.Contains(t, providers, "assemblyai")
	assert.Contains(t, providers, "mock")
	assert.Len(t, providers, 2)
}

func TestAudioProcessorFactory_GetProviderCapabilities(t *testing.T) {
	factory := NewAudioProcessorFactory()

	// Test AssemblyAI capabilities
	assemblyAICapabilities, err := factory.GetProviderCapabilities("assemblyai")
	assert.NoError(t, err)
	assert.Equal(t, "assemblyai", assemblyAICapabilities.Provider)
	assert.True(t, assemblyAICapabilities.SupportsDiarization)
	assert.True(t, assemblyAICapabilities.SupportsRealTime)
	assert.True(t, assemblyAICapabilities.SupportsBatch)
	assert.Contains(t, assemblyAICapabilities.SupportedModes, services.RealTimeMode)
	assert.Contains(t, assemblyAICapabilities.SupportedModes, services.BatchMode)

	// Test Mock capabilities
	mockCapabilities, err := factory.GetProviderCapabilities("mock")
	assert.NoError(t, err)
	assert.Equal(t, "mock", mockCapabilities.Provider)
	assert.True(t, mockCapabilities.SupportsDiarization)

	// Test unknown provider
	_, err = factory.GetProviderCapabilities("unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown provider")
}

func TestAudioProcessorFactory_CreateProcessor_MockProvider(t *testing.T) {
	factory := NewAudioProcessorFactory()

	options := services.AudioProcessingOptions{
		Provider:           "mock",
		Mode:               services.BatchMode,
		SpeakerDiarization: true,
	}

	processor, err := factory.CreateProcessor(services.BatchMode, options)
	assert.NoError(t, err)
	assert.NotNil(t, processor)

	// Verify it supports the expected modes
	supportedModes := processor.GetSupportedModes()
	assert.Contains(t, supportedModes, services.BatchMode)
	assert.Contains(t, supportedModes, services.RealTimeMode)
}

func TestAudioProcessorFactory_CreateProcessor_AssemblyAIWithoutAPIKey(t *testing.T) {
	// Ensure no API key is set
	originalKey := os.Getenv("ASSEMBLYAI_API_KEY")
	os.Unsetenv("ASSEMBLYAI_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("ASSEMBLYAI_API_KEY", originalKey)
		}
	}()

	factory := NewAudioProcessorFactory()

	options := services.AudioProcessingOptions{
		Provider:           "assemblyai",
		Mode:               services.BatchMode,
		SpeakerDiarization: true,
	}

	processor, err := factory.CreateProcessor(services.BatchMode, options)
	assert.NoError(t, err)
	assert.NotNil(t, processor)
	// Should fall back to mock provider when no API key is available
}

func TestAudioProcessorFactory_CreateProcessor_AssemblyAIWithAPIKey(t *testing.T) {
	// Set a test API key
	originalKey := os.Getenv("ASSEMBLYAI_API_KEY")
	os.Setenv("ASSEMBLYAI_API_KEY", "test-api-key-12345")
	defer func() {
		if originalKey != "" {
			os.Setenv("ASSEMBLYAI_API_KEY", originalKey)
		} else {
			os.Unsetenv("ASSEMBLYAI_API_KEY")
		}
	}()

	factory := NewAudioProcessorFactory()

	options := services.AudioProcessingOptions{
		Provider:           "assemblyai",
		Mode:               services.BatchMode,
		SpeakerDiarization: true,
	}

	processor, err := factory.CreateProcessor(services.BatchMode, options)
	assert.NoError(t, err)
	assert.NotNil(t, processor)
}

func TestAudioProcessorFactory_CreateProcessor_DefaultProvider(t *testing.T) {
	factory := NewAudioProcessorFactory()

	options := services.AudioProcessingOptions{
		Provider: "", // Empty provider should default to mock
		Mode:     services.BatchMode,
	}

	processor, err := factory.CreateProcessor(services.BatchMode, options)
	assert.NoError(t, err)
	assert.NotNil(t, processor)
}

func TestAudioProcessorFactory_RecommendProvider(t *testing.T) {
	factory := NewAudioProcessorFactory()

	// Test cost-optimized recommendation
	costOptimizedOptions := services.AudioProcessingOptions{
		CostOptimized: true,
		Mode:          services.BatchMode,
	}
	provider, err := factory.RecommendProvider(costOptimizedOptions)
	assert.NoError(t, err)
	assert.Equal(t, "mock", provider)

	// Test real-time recommendation
	realTimeOptions := services.AudioProcessingOptions{
		RealTimeTranscription: true,
		Mode:                  services.RealTimeMode,
	}
	provider, err = factory.RecommendProvider(realTimeOptions)
	assert.NoError(t, err)
	assert.Equal(t, "assemblyai", provider)

	// Test speaker diarization recommendation
	diarizationOptions := services.AudioProcessingOptions{
		SpeakerDiarization: true,
		Mode:               services.BatchMode,
	}
	provider, err = factory.RecommendProvider(diarizationOptions)
	assert.NoError(t, err)
	assert.Equal(t, "assemblyai", provider)

	// Test default recommendation
	defaultOptions := services.AudioProcessingOptions{
		Mode: services.BatchMode,
	}
	provider, err = factory.RecommendProvider(defaultOptions)
	assert.NoError(t, err)
	assert.Equal(t, "assemblyai", provider)
}

func TestAudioProcessorFactory_CreateProcessor_ErrorCases(t *testing.T) {
	factory := NewAudioProcessorFactory()

	// Test with unknown provider
	options := services.AudioProcessingOptions{
		Provider: "unknown-provider",
		Mode:     services.BatchMode,
	}

	processor, err := factory.CreateProcessor(services.BatchMode, options)
	// Should fall back to mock for unknown provider
	assert.NoError(t, err)
	assert.NotNil(t, processor)
}

// Integration test that verifies the full flow
func TestAudioProcessorFactory_IntegrationTest(t *testing.T) {
	factory := NewAudioProcessorFactory()

	// Test that we can create both types of processors
	providers := []string{"mock", "assemblyai"}
	modes := []services.ProcessingMode{services.BatchMode, services.RealTimeMode}

	for _, provider := range providers {
		for _, mode := range modes {
			options := services.AudioProcessingOptions{
				Provider:           provider,
				Mode:               mode,
				SpeakerDiarization: true,
				Language:           "en",
			}

			processor, err := factory.CreateProcessor(mode, options)
			assert.NoError(t, err)
			assert.NotNil(t, processor)

			// Verify the processor supports the expected modes
			supportedModes := processor.GetSupportedModes()
			assert.Contains(t, supportedModes, mode)
		}
	}
}
