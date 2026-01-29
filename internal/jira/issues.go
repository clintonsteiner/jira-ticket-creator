package jira

import (
	"fmt"
)

// IssueService handles JIRA issue operations
type IssueService struct {
	client *Client
}

// NewIssueService creates a new issue service
func NewIssueService(client *Client) *IssueService {
	return &IssueService{client: client}
}

// CreateIssue creates a new JIRA issue
func (s *IssueService) CreateIssue(projectKey, summary, description, issueType string) (*CreateIssueResponse, error) {
	fields := IssueFields{
		Project: Project{
			Key: projectKey,
		},
		Summary:     summary,
		Description: description,
		IssueType: IssueType{
			Name: issueType,
		},
	}

	req := CreateIssueRequest{
		Fields: fields,
	}

	var resp CreateIssueResponse
	if err := s.client.Do("POST", "/rest/api/2/issue", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return &resp, nil
}

// CreateIssueWithFields creates a new JIRA issue with full field control
func (s *IssueService) CreateIssueWithFields(fields IssueFields) (*CreateIssueResponse, error) {
	req := CreateIssueRequest{
		Fields: fields,
	}

	var resp CreateIssueResponse
	if err := s.client.Do("POST", "/rest/api/2/issue", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return &resp, nil
}

// GetIssue retrieves an issue by key
func (s *IssueService) GetIssue(key string) (*Issue, error) {
	issue, err := s.client.GetIssue(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %s: %w", key, err)
	}
	return issue, nil
}

// UpdateIssue updates an existing issue
func (s *IssueService) UpdateIssue(key string, fields IssueFields) error {
	req := UpdateIssueRequest{
		Fields: fields,
	}

	path := fmt.Sprintf("/rest/api/2/issue/%s", key)
	if err := s.client.Do("PUT", path, req, nil); err != nil {
		return fmt.Errorf("failed to update issue %s: %w", key, err)
	}

	return nil
}

// GetTransitions retrieves available transitions for an issue
func (s *IssueService) GetTransitions(key string) ([]Transition, error) {
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", key)
	var resp TransitionsResponse
	if err := s.client.Do("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get transitions for %s: %w", key, err)
	}
	return resp.Transitions, nil
}

// TransitionIssue transitions an issue to a new state
func (s *IssueService) TransitionIssue(key, transitionID string) error {
	req := TransitionRequest{}
	req.Transition.ID = transitionID

	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", key)
	if err := s.client.Do("POST", path, req, nil); err != nil {
		return fmt.Errorf("failed to transition issue %s: %w", key, err)
	}

	return nil
}

// SearchIssues searches for issues using JQL
func (s *IssueService) SearchIssues(jql string, startAt, maxResults int) ([]Issue, error) {
	resp, err := s.client.GetIssueByJQL(jql, startAt, maxResults)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}
	return resp.Issues, nil
}
