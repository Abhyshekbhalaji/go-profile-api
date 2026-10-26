package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGetUserE2E exercises the real ServeMux and GET handler through HTTP.
func TestGetUserE2E(t *testing.T) {
	server := httptest.NewServer(newMux())
	defer server.Close()

	tests := []struct {
		name            string
		path            string
		wantStatus      int
		wantBodyContains string
	}{
		{
			name:             "existing user",
			path:             "/users/1",
			wantStatus:       http.StatusOK,
			wantBodyContains: "Abhyshek",
		},
		{
			name:             "missing user",
			path:             "/users/3",
			wantStatus:       http.StatusNotFound,
			wantBodyContains: "No such profile exists",
		},
	}

	client := server.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(server.URL + tt.path)
			if err != nil {
				t.Fatalf("GET %s failed: %v", tt.path, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("reading response body failed: %v", err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("GET %s: got status %d, want %d",
					tt.path, resp.StatusCode, tt.wantStatus)
			}

			if !strings.Contains(string(body), tt.wantBodyContains) {
				t.Errorf("GET %s: body %q does not contain %q",
					tt.path, string(body), tt.wantBodyContains)
			}
		})
	}
}
