package jira

import "time"

// Issue represents a JIRA issue
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
	Self   string      `json:"self,omitempty"`
	ID     string      `json:"id,omitempty"`
}

// IssueFields contains the fields for creating/updating a JIRA issue
type IssueFields struct {
	Project      Project                `json:"project"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	IssueType    IssueType              `json:"issuetype"`
	Priority     *Priority              `json:"priority,omitempty"`
	Assignee     *User                  `json:"assignee,omitempty"`
	Status       *Status                `json:"status,omitempty"`
	Labels       []string               `json:"labels,omitempty"`
	Components   []Component            `json:"components,omitempty"`
	CustomFields map[string]interface{} `json:"customfields,omitempty"`
}

// Project represents a JIRA project reference
type Project struct {
	Key string `json:"key"`
}

// IssueType represents a JIRA issue type
type IssueType struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// Priority represents a JIRA priority level
type Priority struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// Status represents a JIRA issue status
type Status struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// User represents a JIRA user
type User struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	AccountID    string `json:"accountId,omitempty"`
}

// Component represents a JIRA component
type Component struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// IssueLink represents a link between two issues
type IssueLink struct {
	Type         LinkType `json:"type"`
	InwardIssue  *Issue   `json:"inwardIssue,omitempty"`
	OutwardIssue *Issue   `json:"outwardIssue,omitempty"`
}

// LinkType represents the type of link (e.g., "Blocks", "Relates")
type LinkType struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// CreateIssueRequest is the request body for creating an issue
type CreateIssueRequest struct {
	Fields IssueFields `json:"fields"`
}

// CreateIssueResponse is the response from creating an issue
type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// LinkIssueRequest is the request body for linking issues
type LinkIssueRequest struct {
	Type         LinkType `json:"type"`
	InwardIssue  *Issue   `json:"inwardIssue,omitempty"`
	OutwardIssue *Issue   `json:"outwardIssue,omitempty"`
}

// TicketRecord represents a ticket that was created, used for persistence
type TicketRecord struct {
	Key              string     `json:"key"`
	Summary          string     `json:"summary"`
	Status           string     `json:"status"`
	BlockedBy        []string   `json:"blocked_by"`
	CreatedAt        time.Time  `json:"created_at"`
	Creator          string     `json:"creator"`                 // Who created this ticket
	Assignee         string     `json:"assignee"`                // Current assignee
	EstimatedEndDate *time.Time `json:"estimated_end,omitempty"` // When is it estimated to be done
	Priority         string     `json:"priority"`
	IssueType        string     `json:"issue_type"`
	Project          string     `json:"project,omitempty"` // Logical project name for grouping
}

// SearchResponse is the response from a search query
type SearchResponse struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// Transition represents a workflow transition
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   struct {
		Self           string `json:"self"`
		Description    string `json:"description"`
		IconURL        string `json:"iconUrl"`
		Name           string `json:"name"`
		ID             string `json:"id"`
		StatusCategory struct {
			Self      string `json:"self"`
			ID        int    `json:"id"`
			Key       string `json:"key"`
			ColorName string `json:"colorName"`
			Name      string `json:"name"`
		} `json:"statusCategory"`
	} `json:"to"`
}

// TransitionsResponse is the response from getting transitions
type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

// TransitionRequest is the request body for transitioning an issue
type TransitionRequest struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Fields IssueFields `json:"fields,omitempty"`
}

// UpdateIssueRequest is the request body for updating an issue
type UpdateIssueRequest struct {
	Fields IssueFields `json:"fields"`
}
