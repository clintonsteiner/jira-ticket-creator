package api

import (
	"C"
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// CreateTicketRequest represents ticket creation parameters
type CreateTicketRequest struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	IssueType   string   `json:"issue_type"`
	Priority    string   `json:"priority"`
	Assignee    string   `json:"assignee"`
	Labels      []string `json:"labels"`
	BlockedBy   []string `json:"blocked_by"`
}

// CreateTicketResponse represents the response from ticket creation
type CreateTicketResponse struct {
	Key   string `json:"key"`
	ID    string `json:"id"`
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

//export CreateTicketJSON
func CreateTicketJSON(urlC *C.char, emailC *C.char, tokenC *C.char, projectC *C.char, jsonC *C.char) *C.char {
	url := C.GoString(urlC)
	email := C.GoString(emailC)
	token := C.GoString(tokenC)
	project := C.GoString(projectC)
	jsonStr := C.GoString(jsonC)

	var req CreateTicketRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		resp := CreateTicketResponse{Error: fmt.Sprintf("Invalid JSON: %v", err)}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	if req.Summary == "" {
		resp := CreateTicketResponse{Error: "summary is required"}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	client := &jira.Client{
		BaseURL: url,
		Email:   email,
		Token:   token,
	}

	if req.IssueType == "" {
		req.IssueType = "Task"
	}
	if req.Priority == "" {
		req.Priority = "Medium"
	}

	fields := jira.IssueFields{
		Project:     jira.Project{Key: project},
		Summary:     req.Summary,
		Description: req.Description,
		IssueType:   jira.IssueType{Name: req.IssueType},
		Priority:    &jira.Priority{Name: req.Priority},
		Labels:      req.Labels,
	}

	if req.Assignee != "" {
		fields.Assignee = &jira.User{EmailAddress: req.Assignee}
	}

	response, err := client.CreateIssue(fields)
	if err != nil {
		resp := CreateTicketResponse{Error: err.Error()}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	resp := CreateTicketResponse{
		Key: response.Key,
		ID:  response.ID,
		URL: response.Self,
	}
	data, _ := json.Marshal(resp)
	return C.CString(string(data))
}

//export ExtractProjectKey
func ExtractProjectKey(ticketKeyC *C.char) *C.char {
	ticketKey := C.GoString(ticketKeyC)
	parts := strings.Split(strings.TrimSpace(ticketKey), "-")
	if len(parts) < 2 {
		return C.CString("")
	}
	return C.CString(parts[0])
}

//export FreeMemory
func FreeMemory(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

//export Version
func Version() *C.char {
	return C.CString("1.0.0")
}
