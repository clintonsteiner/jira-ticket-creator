package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Issue structures
type IssueFields struct {
	Project     map[string]string `json:"project"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	IssueType   map[string]string `json:"issuetype"`
}

type Issue struct {
	Fields IssueFields `json:"fields"`
}

// Create a Jira issue
func CreateIssue(jiraURL, email, apiToken, projectKey, summary, description, issueType string) (string, error) {
	issue := Issue{
		Fields: IssueFields{
			Project:     map[string]string{"key": projectKey},
			Summary:     summary,
			Description: description,
			IssueType:   map[string]string{"name": issueType},
		},
	}

	data, _ := json.Marshal(issue)
	req, _ := http.NewRequest("POST", jiraURL+"/rest/api/2/issue", bytes.NewBuffer(data))
	req.SetBasicAuth(email, apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if key, ok := res["key"].(string); ok {
		return key, nil
	}
	return "", fmt.Errorf("failed to create issue: %v", res)
}

// Link two issues
func LinkIssue(jiraURL, email, apiToken, blocker, blocked string) error {
	link := map[string]interface{}{
		"type":         map[string]string{"name": "Blocks"},
		"inwardIssue":  map[string]string{"key": blocked},
		"outwardIssue": map[string]string{"key": blocker},
	}
	data, _ := json.Marshal(link)
	req, _ := http.NewRequest("POST", jiraURL+"/rest/api/2/issueLink", bytes.NewBuffer(data))
	req.SetBasicAuth(email, apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("failed to link issues, status: %v", resp.Status)
}
