package assemblyai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// CreateTranscript creates a new transcript
func (c *Client) CreateTranscript(ctx context.Context, request *TranscriptRequest) (*Transcript, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/transcript", request)
	if err != nil {
		return nil, err
	}

	var transcript Transcript
	if err := c.handleResponse(resp, &transcript); err != nil {
		return nil, err
	}

	return &transcript, nil
}

// GetTranscript retrieves a transcript by ID
func (c *Client) GetTranscript(ctx context.Context, transcriptID string) (*Transcript, error) {
	endpoint := fmt.Sprintf("/transcript/%s", transcriptID)
	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var transcript Transcript
	if err := c.handleResponse(resp, &transcript); err != nil {
		return nil, err
	}

	return &transcript, nil
}

// ListTranscripts lists transcripts with optional pagination
func (c *Client) ListTranscripts(ctx context.Context, limit *int, beforeID, afterID *string) (*ListTranscriptsResponse, error) {
	endpoint := "/transcript"

	// Build query parameters
	params := url.Values{}
	if limit != nil {
		params.Set("limit", strconv.Itoa(*limit))
	}
	if beforeID != nil {
		params.Set("before_id", *beforeID)
	}
	if afterID != nil {
		params.Set("after_id", *afterID)
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var listResponse ListTranscriptsResponse
	if err := c.handleResponse(resp, &listResponse); err != nil {
		return nil, err
	}

	return &listResponse, nil
}

// DeleteTranscript deletes a transcript by ID
func (c *Client) DeleteTranscript(ctx context.Context, transcriptID string) error {
	endpoint := fmt.Sprintf("/transcript/%s", transcriptID)
	resp, err := c.makeRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// WaitForTranscript waits for a transcript to complete processing
func (c *Client) WaitForTranscript(ctx context.Context, transcriptID string, pollInterval time.Duration) (*Transcript, error) {
	if pollInterval == 0 {
		pollInterval = 3 * time.Second
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			transcript, err := c.GetTranscript(ctx, transcriptID)
			if err != nil {
				return nil, err
			}

			switch transcript.Status {
			case StatusCompleted:
				return transcript, nil
			case StatusError:
				if transcript.Error != nil {
					return nil, fmt.Errorf("transcript failed: %s", *transcript.Error)
				}
				return nil, fmt.Errorf("transcript failed with unknown error")
			case StatusQueued, StatusProcessing:
				// Continue polling
				continue
			default:
				return nil, fmt.Errorf("unknown transcript status: %s", transcript.Status)
			}
		}
	}
}

// TranscribeFromURL creates a transcript from an audio URL and waits for completion
func (c *Client) TranscribeFromURL(ctx context.Context, audioURL string, options *TranscriptRequest) (*Transcript, error) {
	request := &TranscriptRequest{
		AudioURL: audioURL,
	}

	// Merge options if provided
	if options != nil {
		// Copy all non-zero values from options to request
		if options.LanguageCode != nil {
			request.LanguageCode = options.LanguageCode
		}
		if options.Punctuate != nil {
			request.Punctuate = options.Punctuate
		}
		if options.FormatText != nil {
			request.FormatText = options.FormatText
		}
		if options.DualChannel != nil {
			request.DualChannel = options.DualChannel
		}
		if options.WebhookURL != nil {
			request.WebhookURL = options.WebhookURL
		}
		if options.WebhookAuthHeaderName != nil {
			request.WebhookAuthHeaderName = options.WebhookAuthHeaderName
		}
		if options.WebhookAuthHeaderValue != nil {
			request.WebhookAuthHeaderValue = options.WebhookAuthHeaderValue
		}
		if options.AutoHighlights != nil {
			request.AutoHighlights = options.AutoHighlights
		}
		if options.AudioStartFrom != nil {
			request.AudioStartFrom = options.AudioStartFrom
		}
		if options.AudioEndAt != nil {
			request.AudioEndAt = options.AudioEndAt
		}
		if options.WordBoost != nil {
			request.WordBoost = options.WordBoost
		}
		if options.BoostParam != nil {
			request.BoostParam = options.BoostParam
		}
		if options.FilterProfanity != nil {
			request.FilterProfanity = options.FilterProfanity
		}
		if options.RedactPII != nil {
			request.RedactPII = options.RedactPII
		}
		if options.RedactPIIAudio != nil {
			request.RedactPIIAudio = options.RedactPIIAudio
		}
		if options.RedactPIIPolicies != nil {
			request.RedactPIIPolicies = options.RedactPIIPolicies
		}
		if options.RedactPIISub != nil {
			request.RedactPIISub = options.RedactPIISub
		}
		if options.SpeakerLabels != nil {
			request.SpeakerLabels = options.SpeakerLabels
		}
		if options.SpeakersExpected != nil {
			request.SpeakersExpected = options.SpeakersExpected
		}
		if options.ContentSafety != nil {
			request.ContentSafety = options.ContentSafety
		}
		if options.IabCategories != nil {
			request.IabCategories = options.IabCategories
		}
		if options.LanguageDetection != nil {
			request.LanguageDetection = options.LanguageDetection
		}
		if options.CustomSpelling != nil {
			request.CustomSpelling = options.CustomSpelling
		}
		if options.Disfluencies != nil {
			request.Disfluencies = options.Disfluencies
		}
		if options.SentimentAnalysis != nil {
			request.SentimentAnalysis = options.SentimentAnalysis
		}
		if options.AutoChapters != nil {
			request.AutoChapters = options.AutoChapters
		}
		if options.EntityDetection != nil {
			request.EntityDetection = options.EntityDetection
		}
		if options.SpeechThreshold != nil {
			request.SpeechThreshold = options.SpeechThreshold
		}
		if options.Summarization != nil {
			request.Summarization = options.Summarization
		}
		if options.SummaryModel != nil {
			request.SummaryModel = options.SummaryModel
		}
		if options.SummaryType != nil {
			request.SummaryType = options.SummaryType
		}
		if options.CustomTopics != nil {
			request.CustomTopics = options.CustomTopics
		}
		if options.Topics != nil {
			request.Topics = options.Topics
		}
		if options.AdditionalProperties != nil {
			request.AdditionalProperties = options.AdditionalProperties
		}
	}

	// Create the transcript
	transcript, err := c.CreateTranscript(ctx, request)
	if err != nil {
		return nil, err
	}

	// Wait for completion
	return c.WaitForTranscript(ctx, transcript.ID, 0)
}
