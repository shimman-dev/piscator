package piscator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/briandowns/spinner"
)

// Repo is a struct for a GitHub repository
type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

// RepoModel is the struct for a GitHub repository
type RepoModel struct {
	Repo           // embed Repo struct
	Lang    string `json:"language"`
	Fork    bool   `json:"fork"`
	Private bool   `json:"private"`
	Size    uint   `json:"size"`
}

// RepoCollection is a collection of RepoModel structs
type RepoCollection struct {
	Repos []*RepoModel `json:"repos"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RealHttpClient struct{}

func (c RealHttpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

type Sleeper interface {
	Sleep(d time.Duration)
}

type RealSleeper struct{}

func (rs RealSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Retrieves repositories of a user/organization/self from GitHub.
// Optionally filters based on fork status, and returns them as a JSON string or writes to a file.
func GetRepos(client HttpClient, sleeper Sleeper, name, token string, isSelf, isOrg, isForked, makeFile bool) (string, error) {
	var githubURL string

	gh, err := url.Parse("https://api.github.com/")
	if err != nil {
		return "", err
	}

	switch {
	case isSelf:
		gh.Path = path.Join("user", "repos")
	case isOrg:
		gh.Path = path.Join("orgs", name, "repos")
	default:
		gh.Path = path.Join("users", name, "repos")
	}

	params := url.Values{}
	params.Add("per_page", "1000")
	githubURL = gh.String() + "?" + params.Encode()

	var res *http.Response
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", githubURL, nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("User-Agent", "shimman-dev/piscator")
		if token != "" {
			req.Header.Set("Accept", "application/vnd.github+json")
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, err = client.Do(req)

		if err != nil {
			log.Printf("Attempt %d: failed to get repos: %v", i+1, err)
			sleeper.Sleep(2 * time.Second)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Printf("Attempt %d: unexpected status code: %d", i+1, res.StatusCode)
			sleeper.Sleep(2 * time.Second)
			continue
		}

		if remaining := res.Header.Get("X-Ratelimit-Remaining"); remaining == "0" {
			resetTimeStr := res.Header.Get("X-Ratelimit-Reset")
			resetTimeUnix, _ := strconv.ParseInt(resetTimeStr, 10, 64)
			resetTime := time.Unix(resetTimeUnix, 0)
			log.Printf("Attempt %d: rate limit exceeded, sleeping until %v", i+1, resetTime)
			sleeper.Sleep(time.Until(resetTime))
		} else {
			break
		}
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("failed to get repos after 3 attempts")
	}

	var repos []RepoModel
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		return "", err
	}

	filteredRepos := []RepoModel{}
	if isForked {
		for _, repo := range repos {
			if repo.Name != "" && repo.URL != "" {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	} else {
		for _, repo := range repos {
			if !repo.Fork && repo.Name != "" && repo.URL != "" {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	}

	// Marshal filteredRepos into JSON
	jsonData, err := json.MarshalIndent(filteredRepos, "", "  ")
	if err != nil {
		return "", err
	}

	if makeFile {
		// Open file to write
		file, err := os.Create("repos.json")
		if err != nil {
			return "", err
		}
		defer file.Close()

		// Write JSON to file
		_, err = file.Write(jsonData)
		if err != nil {
			return "", err
		}
	}

	return string(jsonData), nil
}

// Filters repositories from a JSON string by programming language and returns them as a JSON string.
func RepoByLanguage(jsonStr string, language string) (string, error) {
	var repos []RepoModel
	if err := json.Unmarshal([]byte(jsonStr), &repos); err != nil {
		return "", err
	}

	filteredRepos := []RepoModel{}
	for _, repo := range repos {
		if repo.Lang == language {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	jsonData, err := json.MarshalIndent(filteredRepos, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

type CommandExecutor interface {
	ExecuteCommand(name string, arg ...string) ([]byte, error)
	ExecuteCommandInDir(dir, name string, arg ...string) ([]byte, error)
}

type RealCommandExecutor struct{}

func (r RealCommandExecutor) ExecuteCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.CombinedOutput()
}

func (r RealCommandExecutor) ExecuteCommandInDir(dir, name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// Clones GitHub repositories from a JSON string concurrently, updates if they already exist, and logs progress.
func CloneReposFromJson(executor CommandExecutor, jsonStr, dirName string, concurrentLimit int8, verboseLog bool) error {
	// unmarshal the JSON string into a slice of Repo structs
	var repos []Repo
	if err := json.Unmarshal([]byte(jsonStr), &repos); err != nil {
		return err
	}

	// create a directory for repos if it doesn't already exist
	dir := dirName
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	// use a WaitGroup to wait for all clones to finish
	var wg sync.WaitGroup
	wg.Add(len(repos))

	var counter uint64 = 1
	sem := make(chan struct{}, concurrentLimit)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	errors := make(chan error) // Create a new error channel

	// clone each repo in a separate goroutine
	for _, repo := range repos {
		go func(repo Repo) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()

			repoPath := path.Join(dir, repo.Name)
			var cmdOut []byte
			var err error
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				// repo doesn't exist, clone it
				cloneCmd := []string{"git", "clone"}
				if isSSHURL(repo.URL) {
					cloneCmd = append(cloneCmd, "--ssh")
				}
				cloneCmd = append(cloneCmd, repo.URL, repoPath)
				cmdOut, err = executor.ExecuteCommand(cloneCmd[0], cloneCmd[1:]...)

				if err != nil {
					errors <- fmt.Errorf("error cloning repo: %w", err) // Send error to channel
					return
				}
			} else if err != nil {
				errors <- fmt.Errorf("error checking if repo exists: %w", err)
				return
			} else {
				// repo exists, pull latest changes
				cmdOut, err = executor.ExecuteCommandInDir(repoPath, "git", "pull")
				if err != nil {
					errors <- fmt.Errorf("error pulling latest changes: %w", err)
					return
				}
			}
			if err != nil {
				fmt.Printf("failed to clone %s: %s\n", repo.URL, string(cmdOut))
			}

			if verboseLog {
				// TODO: more succinct messaging
				log.Printf("Cloned %s/%s\n", dir, repo.Name)
			}

			atomic.AddUint64(&counter, 1)
			s.Suffix = fmt.Sprintf(" Cloning %d/%d repos\n", counter, len(repos))
		}(repo)
	}

	// wait for all clones to finish
	go func() {
		wg.Wait()
		close(errors) // Close the channel when all goroutines are done
	}()

	// Check for any errors from the goroutines
	for err := range errors {
		if err != nil {
			return err // Return the first error that occurred
		}
	}

	s.Stop()
	fmt.Printf("Cloned %d repos\n", len(repos))
	return nil
}

// Checks if the URL is using the SSH scheme
func isSSHURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme == "ssh"
}
