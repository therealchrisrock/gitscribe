package assemblyai

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// UploadFile uploads a local file to AssemblyAI and returns the upload URL
func (c *Client) UploadFile(ctx context.Context, reader io.Reader) (*UploadResponse, error) {
	// Create the upload request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/upload", reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	var uploadResponse UploadResponse
	if err := c.handleResponse(resp, &uploadResponse); err != nil {
		return nil, err
	}

	return &uploadResponse, nil
}

// TranscribeFromReader uploads a file from a reader and transcribes it
func (c *Client) TranscribeFromReader(ctx context.Context, reader io.Reader, options *TranscriptRequest) (*Transcript, error) {
	// Upload the file first
	uploadResp, err := c.UploadFile(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Now transcribe using the upload URL
	return c.TranscribeFromURL(ctx, uploadResp.UploadURL, options)
}
