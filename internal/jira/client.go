package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client represents a JIRA API client with retry logic
type Client struct {
	BaseURL    string
	Email      string
	Token      string
	HTTPClient *http.Client
	MaxRetries int
}

// NewClient creates a new JIRA API client
func NewClient(baseURL, email, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Email:   email,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		MaxRetries: 3,
	}
}

// Do performs an HTTP request with retry logic and exponential backoff
func (c *Client) Do(method, path string, body interface{}, result interface{}) error {
	var retryCount int
	var lastErr error

	for retryCount = 0; retryCount <= c.MaxRetries; retryCount++ {
		err := c.doRequest(method, path, body, result)
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}

		lastErr = err
		if retryCount < c.MaxRetries {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(retryCount)) * time.Second
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set authentication
	auth := base64.StdEncoding.EncodeToString([]byte(c.Email + ":" + c.Token))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		jiraErr := &JiraError{
			StatusCode: resp.StatusCode,
		}

		// Try to parse error response
		var errResp struct {
			ErrorMessages []string               `json:"errorMessages"`
			Errors        map[string]interface{} `json:"errors"`
		}
		json.Unmarshal(respBody, &errResp)
		jiraErr.ErrorMessages = errResp.ErrorMessages
		jiraErr.Errors = errResp.Errors

		// Special handling for specific status codes
		switch resp.StatusCode {
		case 401, 403:
			return &AuthenticationError{Message: jiraErr.Error()}
		case 404:
			return &NotFoundError{Resource: "JIRA resource"}
		case 429:
			retryAfter := 0
			if retryAfterStr := resp.Header.Get("Retry-After"); retryAfterStr != "" {
				retryAfter, _ = strconv.Atoi(retryAfterStr)
			}
			return &RateLimitError{RetryAfter: retryAfter, Message: jiraErr.Error()}
		}

		return jiraErr
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// isRetryableError checks if an error should be retried
func isRetryableError(err error) bool {
	switch e := err.(type) {
	case *JiraError:
		return e.IsRetryable()
	case *RateLimitError:
		return e.IsRetryable()
	default:
		return false
	}
}

// GetIssue retrieves an issue by key
func (c *Client) GetIssue(key string) (*Issue, error) {
	var issue Issue
	path := fmt.Sprintf("/rest/api/2/issue/%s", key)
	if err := c.Do("GET", path, nil, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// GetIssueByJQL retrieves issues using JQL
func (c *Client) GetIssueByJQL(jql string, startAt, maxResults int) (*SearchResponse, error) {
	var result SearchResponse
	escapedJQL := url.QueryEscape(jql)
	path := fmt.Sprintf("/rest/api/2/search?jql=%s&startAt=%d&maxResults=%d",
		escapedJQL, startAt, maxResults)
	if err := c.Do("GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
