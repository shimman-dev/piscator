package piscator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockHttpClient struct{}

func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	// Simulate a successful HTTP response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`[{ "name": "repo1", "html_url": "url1" }]`))),
	}, nil
}

func TestGetRepos(t *testing.T) {
	client := MockHttpClient{}
	tests := []struct {
		name      string
		isOrg     bool
		isPrivate bool
		isForked  bool
		makeFile  bool
		wantError bool
	}{
		{"user1", false, false, false, false, false},
		// Add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetRepos(client, tt.name, tt.isOrg, tt.isPrivate, tt.isForked, tt.makeFile)
			if (err != nil) != tt.wantError {
				t.Errorf("GetRepos() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestRepoByLanguage(t *testing.T) {
	tests := []struct {
		name         string
		jsonStr      string
		language     string
		expected     []RepoModel
		expectingErr bool
	}{
		{
			name: "valid JSON with Go repos",
			jsonStr: `[
				{
					"name": "Repo1",
					"html_url": "https://github.com/user/repo1",
					"language": "Go",
					"fork": false,
					"private": false,
					"size": 100
				},
				{
					"name": "Repo2",
					"html_url": "https://github.com/user/repo2",
					"language": "JavaScript",
					"fork": true,
					"private": false,
					"size": 200
				}
			]`,
			language: "Go",
			expected: []RepoModel{
				{
					Repo: Repo{
						Name: "Repo1",
						URL:  "https://github.com/user/repo1",
					}, Lang: "Go",
					Fork:    false,
					Private: false,
					Size:    100,
				},
			},
		},
		{
			name: "valid JSON with no Go repos",
			jsonStr: `[
				{
					"name": "Repo1",
					"html_url": "https://github.com/user/repo1",
					"language": "Python",
					"fork": false,
					"private": false,
					"size": 100
				},
				{
					"name": "Repo2",
					"html_url": "https://github.com/user/repo2",
					"language": "JavaScript",
					"fork": true,
					"private": false,
					"size": 200
				}
			]`,
			language: "Go",
			expected: []RepoModel{}, // empty slice as no Go repos are expected
		},
		{
			name:         "invalid JSON",
			jsonStr:      `{[}`,
			language:     "Go",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, err := RepoByLanguage(tt.jsonStr, tt.language)
			if tt.expectingErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one")
				}
				return
			}
			if err != nil {
				t.Errorf("Did not expect an error but got: %v", err)
				return
			}
			var got []RepoModel
			err = json.Unmarshal([]byte(gotStr), &got)
			if err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

type MockCommandExecutor struct{}

func (m MockCommandExecutor) ExecuteCommand(name string, arg ...string) ([]byte, error) {
	// Always return success
	return []byte("ok"), nil
}

func (m MockCommandExecutor) ExecuteCommandInDir(dir, name string, arg ...string) ([]byte, error) {
	// Always return success
	return []byte("ok"), nil
}

func TestCloneReposFromJson(t *testing.T) {
	executor := MockCommandExecutor{}
	tests := []struct {
		name            string
		jsonStr         string
		concurrentLimit int8
		verboseLog      bool
		wantError       bool
	}{
		{"user1", `[{ "name": "repo1", "html_url": "url1" }]`, 2, false, false},
		// Add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloneReposFromJson(executor, tt.jsonStr, tt.name, tt.concurrentLimit, tt.verboseLog)
			if (err != nil) != tt.wantError {
				t.Errorf("CloneReposFromJson() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
