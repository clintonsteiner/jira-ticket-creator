//go:build !test

package main

/*
#include <stdlib.h>
*/
import "C"

import (
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

	request := CreateIssueRequest{
		Fields: fields,
	}

	var response CreateIssueResponse
	if err := client.Do("POST", "/rest/api/2/issue", request, &response); err != nil {
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

//export GetTicket
func GetTicket(urlC *C.char, emailC *C.char, tokenC *C.char, keyC *C.char) *C.char {
	url := C.GoString(urlC)
	email := C.GoString(emailC)
	token := C.GoString(tokenC)
	key := C.GoString(keyC)

	client := &jira.Client{
		BaseURL: url,
		Email:   email,
		Token:   token,
	}

	issue, err := client.GetIssue(key)
	if err != nil {
		resp := map[string]interface{}{"error": err.Error()}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	resp := map[string]interface{}{
		"key":         issue.Key,
		"id":          issue.ID,
		"summary":     issue.Fields.Summary,
		"description": issue.Fields.Description,
		"status":      issue.Fields.Status.Name,
		"issue_type":  issue.Fields.IssueType.Name,
		"priority":    getStringPtr(issue.Fields.Priority, "Name"),
		"assignee":    getStringPtr(issue.Fields.Assignee, "DisplayName"),
		"labels":      issue.Fields.Labels,
		"url":         issue.Self,
	}

	data, _ := json.Marshal(resp)
	return C.CString(string(data))
}

// Helper to safely extract string pointers
func getStringPtr(v interface{}, field string) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(*jira.Priority); ok && field == "Name" && s != nil {
		return s.Name
	}
	if u, ok := v.(*jira.User); ok && field == "DisplayName" && u != nil {
		return u.DisplayName
	}
	return ""
}

//export SearchTickets
func SearchTickets(urlC *C.char, emailC *C.char, tokenC *C.char, jqlC *C.char) *C.char {
	url := C.GoString(urlC)
	email := C.GoString(emailC)
	token := C.GoString(tokenC)
	jql := C.GoString(jqlC)

	client := &jira.Client{
		BaseURL: url,
		Email:   email,
		Token:   token,
	}

	result, err := client.GetIssueByJQL(jql, 0, 50)
	if err != nil {
		resp := map[string]interface{}{"error": err.Error()}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	tickets := make([]map[string]interface{}, 0, len(result.Issues))
	for _, issue := range result.Issues {
		ticket := map[string]interface{}{
			"key":        issue.Key,
			"summary":    issue.Fields.Summary,
			"status":     issue.Fields.Status.Name,
			"issue_type": issue.Fields.IssueType.Name,
			"priority":   getStringPtr(issue.Fields.Priority, "Name"),
			"assignee":   getStringPtr(issue.Fields.Assignee, "DisplayName"),
		}
		tickets = append(tickets, ticket)
	}

	resp := map[string]interface{}{
		"total":   result.Total,
		"count":   len(tickets),
		"tickets": tickets,
	}

	data, _ := json.Marshal(resp)
	return C.CString(string(data))
}

//export UpdateTicket
func UpdateTicket(urlC *C.char, emailC *C.char, tokenC *C.char, keyC *C.char, jsonC *C.char) *C.char {
	url := C.GoString(urlC)
	email := C.GoString(emailC)
	token := C.GoString(tokenC)
	key := C.GoString(keyC)
	jsonStr := C.GoString(jsonC)

	var updateReq map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &updateReq); err != nil {
		resp := map[string]interface{}{"error": fmt.Sprintf("Invalid JSON: %v", err)}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	client := &jira.Client{
		BaseURL: url,
		Email:   email,
		Token:   token,
	}

	// Build IssueFields from the update request
	fields := jira.IssueFields{}

	if summary, ok := updateReq["summary"].(string); ok && summary != "" {
		fields.Summary = summary
	}
	if description, ok := updateReq["description"].(string); ok && description != "" {
		fields.Description = description
	}
	if priority, ok := updateReq["priority"].(string); ok && priority != "" {
		fields.Priority = &jira.Priority{Name: priority}
	}
	if assignee, ok := updateReq["assignee"].(string); ok && assignee != "" {
		fields.Assignee = &jira.User{EmailAddress: assignee}
	}

	// Perform update using client's Do method directly
	updateRequest := struct {
		Fields jira.IssueFields `json:"fields"`
	}{
		Fields: fields,
	}

	path := fmt.Sprintf("/rest/api/2/issue/%s", key)
	if err := client.Do("PUT", path, updateRequest, nil); err != nil {
		resp := map[string]interface{}{"error": err.Error()}
		data, _ := json.Marshal(resp)
		return C.CString(string(data))
	}

	resp := map[string]interface{}{
		"key":     key,
		"message": "Ticket updated successfully",
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

func main() {
	// This function is not called when building as a C library
	// It's here to make this a valid main package for cgo
}
