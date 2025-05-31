package assemblyai

import (
	"context"
	"fmt"
	"net/http"
)

// LemurTask performs a custom LeMUR task
func (c *Client) LemurTask(ctx context.Context, request *LemurRequest) (*LemurResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/lemur/v3/generate/task", request)
	if err != nil {
		return nil, err
	}

	var lemurResponse LemurResponse
	if err := c.handleResponse(resp, &lemurResponse); err != nil {
		return nil, err
	}

	return &lemurResponse, nil
}

// LemurSummary generates a summary using LeMUR
func (c *Client) LemurSummary(ctx context.Context, request *LemurSummaryRequest) (*LemurResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/lemur/v3/generate/summary", request)
	if err != nil {
		return nil, err
	}

	var lemurResponse LemurResponse
	if err := c.handleResponse(resp, &lemurResponse); err != nil {
		return nil, err
	}

	return &lemurResponse, nil
}

// LemurQuestionAnswer performs question and answer using LeMUR
func (c *Client) LemurQuestionAnswer(ctx context.Context, request *LemurQuestionAnswerRequest) (*LemurQuestionAnswerResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/lemur/v3/generate/question-answer", request)
	if err != nil {
		return nil, err
	}

	var lemurResponse LemurQuestionAnswerResponse
	if err := c.handleResponse(resp, &lemurResponse); err != nil {
		return nil, err
	}

	return &lemurResponse, nil
}

// LemurActionItems extracts action items using LeMUR
func (c *Client) LemurActionItems(ctx context.Context, request *LemurActionItemsRequest) (*LemurActionItemsResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/lemur/v3/generate/action-items", request)
	if err != nil {
		return nil, err
	}

	var lemurResponse LemurActionItemsResponse
	if err := c.handleResponse(resp, &lemurResponse); err != nil {
		return nil, err
	}

	return &lemurResponse, nil
}

// GetLemurResponse retrieves a LeMUR response by request ID
func (c *Client) GetLemurResponse(ctx context.Context, requestID string) (*LemurResponse, error) {
	endpoint := fmt.Sprintf("/lemur/v3/%s", requestID)
	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var lemurResponse LemurResponse
	if err := c.handleResponse(resp, &lemurResponse); err != nil {
		return nil, err
	}

	return &lemurResponse, nil
}

// DeleteLemurRequest deletes a LeMUR request by request ID
func (c *Client) DeleteLemurRequest(ctx context.Context, requestID string) error {
	endpoint := fmt.Sprintf("/lemur/v3/%s", requestID)
	resp, err := c.makeRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}
