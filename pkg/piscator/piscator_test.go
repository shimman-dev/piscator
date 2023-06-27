package piscator

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type MockHttpClient struct {
	httpStatus int
	httpBody   string
	httpError  error
	Headers    http.Header
}

func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.httpError != nil {
		return nil, m.httpError
	}

	// Simulate a successful HTTP response
	return &http.Response{
		StatusCode: m.httpStatus,
		Body:       io.NopCloser(bytes.NewReader([]byte(m.httpBody))),
		Header:     m.Headers,
	}, nil
}

type MockSleeper struct {
	Durations []time.Duration
}

func (ms *MockSleeper) Sleep(d time.Duration) {
	ms.Durations = append(ms.Durations, d)
}

func TestGetRepos(t *testing.T) {
	repos := []RepoModel{
		{Repo: Repo{Name: "repo1", URL: "http://example.com/repo1"}, Lang: "Go", Fork: false, Private: false, Size: 100},
		{Repo: Repo{Name: "repo2", URL: "http://example.com/repo2"}, Lang: "Python", Fork: true, Private: false, Size: 200},
		{Repo: Repo{Name: "repo3", URL: "http://example.com/repo3"}, Lang: "Java", Fork: false, Private: true, Size: 300},
	}

	tests := []struct {
		name           string
		token          string
		username       string
		password       string
		enterpriseHost string
		team           string
		isSelf         bool
		isOrg          bool
		isForked       bool
		makeFile       bool
		httpError      error
		httpStatus     int
		httpBody       string
		rateLimit      string
		rateReset      string
		wantError      bool
	}{
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          false,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github",
			team:           "",
			isSelf:         true,
			isOrg:          true,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "",
			password:       "",
			enterpriseHost: "",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "",
			password:       "",
			enterpriseHost: "",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			token:          "token",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "",
			password:       "",
			enterpriseHost: "",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       true,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "user1-makefile",
			token:          "token",
			username:       "",
			password:       "",
			enterpriseHost: "",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      false,
		},
		{
			name:           "forked repos-makefile",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       true,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody: `[
				{
					"name": "repo1",
					"html_url": "https://github.com/user1/repo1",
					"language": "Go",
					"fork": true,
					"private": false,
					"size": 100
				},
				{
					"name": "repo2",
					"html_url": "https://github.com/user1/repo2",
					"language": "JavaScript",
					"fork": false,
					"private": false,
					"size": 200
				},
				{
					"name": "repo3",
					"html_url": "https://github.com/user1/repo3",
					"language": "",
					"fork": false,
					"private": true,
					"size": 300
				}
			]`,
			rateLimit: "5000",
			rateReset: "0",
			wantError: false,
		},
		{
			name:           "forked repos-makefile",
			token:          "token",
			username:       "",
			password:       "",
			enterpriseHost: "",
			team:           "",
			isSelf:         false,
			isOrg:          true,
			isForked:       false,
			makeFile:       true,
			httpError:      nil,
			httpStatus:     200,
			httpBody: `[
				{
					"name": "repo1",
					"html_url": "https://github.com/user1/repo1",
					"language": "Go",
					"fork": true,
					"private": false,
					"size": 100
				},
				{
					"name": "repo2",
					"html_url": "https://github.com/user1/repo2",
					"language": "JavaScript",
					"fork": false,
					"private": false,
					"size": 200
				},
				{
					"name": "repo3",
					"html_url": "https://github.com/user1/repo3",
					"language": "",
					"fork": false,
					"private": true,
					"size": 300
				}
			]`,
			rateLimit: "5000",
			rateReset: "0",
			wantError: false,
		},
		{
			name:           "network error",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      errors.New("network error"),
			httpStatus:     200,
			httpBody:       "",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      true,
		},
		{
			name:           "github 404",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     404,
			httpBody:       "",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      true,
		},
		{
			name:           "invalid json",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "{invalid json}",
			rateLimit:      "5000",
			rateReset:      "0",
			wantError:      true,
		},
		{
			name:           "rate limited",
			token:          "token",
			username:       "tester_mctesterson",
			password:       "hunter2",
			enterpriseHost: "acme.github.com",
			team:           "",
			isSelf:         false,
			isOrg:          false,
			isForked:       false,
			makeFile:       false,
			httpError:      nil,
			httpStatus:     200,
			httpBody:       "[]",
			rateLimit:      "0",
			rateReset:      strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10),
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("X-RateLimit-Remaining", tt.rateLimit)
			headers.Set("X-RateLimit-Reset", tt.rateReset)

			client := MockHttpClient{
				httpError:  tt.httpError,
				httpStatus: tt.httpStatus,
				httpBody:   tt.httpBody,
				Headers:    headers,
			}

			sleeper := &MockSleeper{}

			_, err := GetRepos(client, sleeper, tt.name, tt.token, tt.username, tt.password, tt.enterpriseHost, tt.team, tt.isSelf, tt.isOrg, tt.isForked, tt.makeFile)
			if (err != nil) != tt.wantError {
				t.Errorf("GetRepos() error = %v, wantError %v", err, tt.wantError)
			}

			if tt.name == "forked repos" {
				filteredJSON, err := GetRepos(client, sleeper, tt.name, tt.token, tt.username, tt.password, tt.enterpriseHost, tt.team, tt.isSelf, tt.isOrg, tt.isForked, tt.makeFile)
				if (err != nil) != tt.wantError {
					t.Errorf("GetRepos() error = %v, wantError %v", err, tt.wantError)
				}

				// Unmarshal the filtered JSON for assertion
				var filteredRepos []RepoModel
				if err := json.Unmarshal([]byte(filteredJSON), &filteredRepos); err != nil {
					t.Errorf("Failed to unmarshal filtered JSON: %v", err)
				}

				// Perform assertions on filteredRepos
				if tt.isForked {
					// Check if all repos are present in the filteredRepos slice
					if len(filteredRepos) != len(repos) {
						t.Errorf("Expected all repos to be present in filteredRepos. Expected %d, got %d", len(repos), len(filteredRepos))
					}
				} else {
					// Check if only non-forked repos are present in the filteredRepos slice
					for _, repo := range filteredRepos {
						if repo.Fork {
							t.Errorf("Found a forked repo in filteredRepos: %v", repo)
						}
					}
				}
			}

			// If this is the rate limited test case, check that the correct sleep duration was used
			if tt.name == "rate limited" {
				if len(sleeper.Durations) != 3 {
					t.Errorf("Expected 3 sleeper calls, got %v", len(sleeper.Durations))
				} else {
					// Parse reset time from the test case
					resetTimeUnix, _ := strconv.ParseInt(tt.rateReset, 10, 64)
					resetTime := time.Unix(resetTimeUnix, 0)
					expectedSleep := time.Until(resetTime)

					if sleeper.Durations[0] < expectedSleep-time.Millisecond || sleeper.Durations[0] > expectedSleep+time.Millisecond {
						t.Errorf("Expected sleep for %v, got %v", expectedSleep, sleeper.Durations[0])
					}
				}
			}

			// test cleanup, after it has run
			defer func() {
				if strings.Contains(tt.name, "-makefile") {
					err := os.Remove("repos.json")
					if err != nil {
						t.Errorf("Failed to remove repos.json file: %v", err)
					}
				}
			}()
		})
	}
}

func TestRepoByLanguage(t *testing.T) {
	tests := []struct {
		name         string
		jsonStr      string
		languages    string
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
			languages: "Go",
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
			name: "valid JSON with multiple language repos",
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
				},
				{
					"name": "Repo3",
					"html_url": "https://github.com/user/repo3",
					"language": "Python",
					"fork": false,
					"private": false,
					"size": 300
				}
			]`,
			languages: "Go,Python",
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
				{
					Repo: Repo{
						Name: "Repo3",
						URL:  "https://github.com/user/repo3",
					}, Lang: "Python",
					Fork:    false,
					Private: false,
					Size:    300,
				},
			},
		},
		{
			name: "valid JSON with case-insensitive language filtering",
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
					"language": "javascript",
					"fork": true,
					"private": false,
					"size": 200
				}
			]`,
			languages: "go",
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
			name:         "invalid JSON",
			jsonStr:      `{[}`,
			languages:    "Go",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, err := RepoByLanguage(tt.jsonStr, tt.languages)
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

func TestSplitLanguages(t *testing.T) {
	tests := []struct {
		language string
		expected []string
	}{
		{"go,rust,python", []string{"go", "rust", "python"}},
		{"  typescript,   erlang,   c++  ", []string{"typescript", "erlang", "c++"}},
		{"swift", []string{"swift"}},
	}

	for _, test := range tests {
		result := splitLanguages(test.language)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected splitLanguages(%q) to return %v, but got %v", test.language, test.expected, result)
		}
	}
}

type MockCommandExecutor struct {
	errors map[string]error
}

func (m MockCommandExecutor) ExecuteCommand(name string, arg ...string) ([]byte, error) {
	if err, ok := m.errors["ExecuteCommand_"+name]; ok {
		return []byte("error"), err
	}
	// Always return success
	return []byte("ok"), nil
}

func (m MockCommandExecutor) ExecuteCommandInDir(dir, name string, arg ...string) ([]byte, error) {
	if err, ok := m.errors["ExecuteCommandInDir_"+name]; ok {
		return []byte("error"), err
	}
	// Always return success
	return []byte("ok"), nil
}

func TestCloneReposFromJson(t *testing.T) {
	successExecutor := MockCommandExecutor{
		errors: map[string]error{},
	}

	cloneErrorExecutor := MockCommandExecutor{
		errors: map[string]error{
			"ExecuteCommand_git": errors.New("Simulated clone error"),
		},
	}

	tests := []struct {
		name            string
		jsonStr         string
		executor        CommandExecutor
		concurrentLimit int8
		verboseLog      bool
		wantError       bool
	}{
		// happy path
		{"happy_path", `[{ "name": "repo1", "html_url": "url1" }]`, successExecutor, 2, false, false},
		// trigger an error while cloning
		{"err_cloning", `[{ "name": "repo1", "html_url": "url1" }]`, cloneErrorExecutor, 2, false, true},
		// trigger an error in json.Unmarshal
		{"error_unmarshal", `[{ "name": "repo1", "html_url": "url1", "invalid": : "json" }]`, successExecutor, 2, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloneReposFromJson(tt.executor, tt.jsonStr, tt.name, tt.concurrentLimit, tt.verboseLog)
			if (err != nil) != tt.wantError {
				t.Errorf("CloneReposFromJson() error = %v, wantError %v", err, tt.wantError)
			}

			// test cleanup, after it has run
			defer func() {
				err := os.RemoveAll(tt.name)
				if err != nil {
					t.Errorf("Failed to remove %s directory: %v", tt.name, err)
				}
			}()
		})

	}
}

func TestIsSSHURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "SSH URL",
			url:      "ssh://github.com/user/repo",
			expected: true,
		},
		{
			name:     "HTTP URL",
			url:      "http://github.com/user/repo",
			expected: false,
		},
		{
			name:     "Invalid URL",
			url:      "invalid-url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSSHURL(tt.url)
			if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}
