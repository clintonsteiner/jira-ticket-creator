package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock server for Jira API
func setupMockServer(responseCode int, responseBody any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseCode)
		if responseBody != nil {
			json.NewEncoder(w).Encode(responseBody)
		}
	}))
}

func TestCreateIssue_Success(t *testing.T) {
	mock := setupMockServer(201, map[string]string{"key": "TEST-123"})
	defer mock.Close()

	key, err := CreateIssue(mock.URL, "email", "token", "TEST", "Summary", "Desc", "Task")
	if err != nil {
		t.Fatal(err)
	}
	if key != "TEST-123" {
		t.Fatalf("Expected TEST-123, got %s", key)
	}
}

func TestCreateIssue_Failure(t *testing.T) {
	mock := setupMockServer(400, map[string]string{"error": "bad request"})
	defer mock.Close()

	_, err := CreateIssue(mock.URL, "email", "token", "TEST", "Summary", "Desc", "Task")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestLinkIssue_Success(t *testing.T) {
	mock := setupMockServer(201, nil)
	defer mock.Close()

	err := LinkIssue(mock.URL, "email", "token", "BLOCKER-1", "BLOCKED-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinkIssue_Failure(t *testing.T) {
	mock := setupMockServer(400, nil)
	defer mock.Close()

	err := LinkIssue(mock.URL, "email", "token", "BLOCKER-1", "BLOCKED-1")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}
