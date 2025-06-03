package services

import (
	"testing"

	"teammate/server/modules/transcription/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestEnhancedTranscriptionService_segmentsToText(t *testing.T) {
	service := &EnhancedTranscriptionService{}

	tests := []struct {
		name     string
		segments []entities.TranscriptSegment
		expected string
	}{
		{
			name:     "empty segments",
			segments: []entities.TranscriptSegment{},
			expected: "",
		},
		{
			name: "normal segments with speakers",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "Speaker A",
					Text:    "Hello everyone, welcome to the meeting.",
				},
				{
					Speaker: "Speaker B",
					Text:    "Thank you for having me here today.",
				},
			},
			expected: "Speaker A: Hello everyone, welcome to the meeting.\nSpeaker B: Thank you for having me here today.\n",
		},
		{
			name: "segments with Speaker Unknown",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "Speaker Unknown",
					Text:    "This is from an unknown speaker.",
				},
				{
					Speaker: "Speaker A",
					Text:    "This is from a known speaker.",
				},
			},
			expected: "Speaker Unknown: This is from an unknown speaker.\nSpeaker A: This is from a known speaker.\n",
		},
		{
			name: "segments with speaker_unknown (should be handled)",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "speaker_unknown",
					Text:    "This should not show speaker label.",
				},
				{
					Speaker: "Speaker A",
					Text:    "This should show speaker label.",
				},
			},
			expected: "This should not show speaker label.\nSpeaker A: This should show speaker label.\n",
		},
		{
			name: "segments with empty text (should be filtered out)",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "Speaker A",
					Text:    "This should appear.",
				},
				{
					Speaker: "speaker_unknown",
					Text:    "",
				},
				{
					Speaker: "Speaker B",
					Text:    "   ",
				},
				{
					Speaker: "Speaker C",
					Text:    "This should also appear.",
				},
			},
			expected: "Speaker A: This should appear.\nSpeaker C: This should also appear.\n",
		},
		{
			name: "mixed scenario with empty text and speaker_unknown",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "Speaker A",
					Text:    "Hello world.",
				},
				{
					Speaker: "speaker_unknown",
					Text:    "",
				},
				{
					Speaker: "speaker_unknown",
					Text:    "Some text from unknown speaker.",
				},
				{
					Speaker: "Speaker B",
					Text:    "   \n  ",
				},
				{
					Speaker: "Speaker C",
					Text:    "Final message.",
				},
			},
			expected: "Speaker A: Hello world.\nSome text from unknown speaker.\nSpeaker C: Final message.\n",
		},
		{
			name: "empty speaker field",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "",
					Text:    "Text with no speaker.",
				},
				{
					Speaker: "Speaker A",
					Text:    "Text with speaker.",
				},
			},
			expected: "Text with no speaker.\nSpeaker A: Text with speaker.\n",
		},
		{
			name: "segments with speaker_unknown labels",
			segments: []entities.TranscriptSegment{
				{
					Speaker: "speaker_unknown", // This should be filtered out (no speaker label)
					Text:    "Hello there",
				},
				{
					Speaker: "Speaker A",
					Text:    "How are you?",
				},
				{
					Speaker: "speaker_0", // This should be preserved as-is
					Text:    "I'm doing well",
				},
			},
			expected: "Hello there\nSpeaker A: How are you?\nspeaker_0: I'm doing well\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.segmentsToText(tt.segments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnhancedTranscriptionService_calculateAverageConfidence(t *testing.T) {
	service := &EnhancedTranscriptionService{}

	tests := []struct {
		name     string
		segments []entities.TranscriptSegment
		expected float64
	}{
		{
			name:     "empty segments",
			segments: []entities.TranscriptSegment{},
			expected: 0.0,
		},
		{
			name: "single segment",
			segments: []entities.TranscriptSegment{
				{Confidence: 0.85},
			},
			expected: 0.85,
		},
		{
			name: "multiple segments",
			segments: []entities.TranscriptSegment{
				{Confidence: 0.90},
				{Confidence: 0.80},
				{Confidence: 0.70},
			},
			expected: 0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateAverageConfidence(tt.segments)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}
