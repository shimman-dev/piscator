package piscator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

// RepoModel is the struct for a GitHub repository
type RepoModel struct {
	Url     string `json:"html_url"`
	Name    string `json:"name"`
	Lang    string `json:"language"`
	Fork    bool   `json:"fork"`
	Private bool   `json:"private"`
	Size    int    `json:"size"`
}

// RepoCollection is a collection of RepoModel structs
type RepoCollection struct {
	Repos []*RepoModel `json:"repos"`
}

// Takes a GitHub username and returns a JSON string of their repos
func GetRepos(name string, isOrg, isPrivate, isForked, makeFile bool) (string, error) {
	var githubURL string

	gh, err := url.Parse("https://api.github.com/")
	if err != nil {
		return "", err
	}

	if isOrg {
		gh.Path = path.Join("orgs", name, "repos")
	} else {
		gh.Path = path.Join("users", name, "repos")
	}

	params := url.Values{}
	params.Add("per_page", "1000")
	if isPrivate {
		params.Add("private", "true")
	}
	githubURL = gh.String() + "?" + params.Encode()

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = http.Get(githubURL)
		if err != nil {
			log.Printf("Attempt %d: failed to get repos: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Attempt %d: unexpected status code: %d", i+1, res.StatusCode)
			time.Sleep(2 * time.Second)
			continue
		}
		defer res.Body.Close()
		break
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("failed to get repos after 3 attempts")
	}
	if remaining := res.Header.Get("X-Ratelimit-Remaining"); remaining == "0" {
		resetTimeStr := res.Header.Get("X-Ratelimit-Reset")
		resetTimeUnix, _ := strconv.ParseInt(resetTimeStr, 10, 64)
		resetTime := time.Unix(resetTimeUnix, 0)
		time.Sleep(time.Until(resetTime))
	}

	var repos []RepoModel
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		return "", err
	}

	filteredRepos := []RepoModel{}
	if isForked {
		for _, repo := range repos {
			if repo.Name != "" && repo.Url != "" {
				filteredRepo := RepoModel{
					Name:    repo.Name,
					Lang:    repo.Lang,
					Private: repo.Private,
					Fork:    repo.Fork,
					Size:    repo.Size,
					Url:     repo.Url,
				}
				filteredRepos = append(filteredRepos, filteredRepo)
			}
		}
	} else {
		for _, repo := range repos {
			if !repo.Fork && repo.Name != "" && repo.Url != "" {
				filteredRepo := RepoModel{
					Name:    repo.Name,
					Lang:    repo.Lang,
					Private: repo.Private,
					Fork:    repo.Fork,
					Size:    repo.Size,
					Url:     repo.Url,
				}
				filteredRepos = append(filteredRepos, filteredRepo)
			}
		}
	}

	// Marshal filteredRepos into JSON
	jsonData, err := json.MarshalIndent(filteredRepos, "", "  ")
	if err != nil {
		return "", err
	}

	if makeFile {
		err = ioutil.WriteFile("repos.json", jsonData, 0644)
		if err != nil {
			return "", err
		}
		log.Print("repos.json created")
	}

	return string(jsonData), nil
}

// Repo is a struct for a GitHub repository
type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

// Filter repos by language
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

// Takes JSON from GetRepos and git clones each repo
func CloneReposFromJson(jsonStr, name string, concurrentLimit int8, verboseLog bool) error {
	// unmarshal the JSON string into a slice of Repo structs
	var repos []Repo
	if err := json.Unmarshal([]byte(jsonStr), &repos); err != nil {
		return err
	}

	// create a directory for repos if it doesn't already exist
	dir := name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	// use a WaitGroup to wait for all clones to finish
	var wg sync.WaitGroup
	wg.Add(len(repos))

	var counter uint8 = 1
	sem := make(chan struct{}, concurrentLimit)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	// clone each repo in a separate goroutine
	for _, repo := range repos {
		go func(repo Repo) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()

			repoPath := path.Join(dir, repo.Name)
			var cmd *exec.Cmd
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				// repo doesn't exist, clone it
				cmd = exec.Command("git", "clone", repo.URL)
				cmd.Dir = dir
			} else {
				// repo exists, pull latest changes
				cmd = exec.Command("git", "pull")
				cmd.Dir = repoPath
			}
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("failed to clone %s: %s\n", repo.URL, string(out))
			}

			if verboseLog {
				// TODO: more succinct messaging
				log.Printf("Cloned %s/%s\n", dir, repo.Name)
			}

			s.Suffix = fmt.Sprintf(" Cloned %d/%d repos\n", counter, len(repos))
			counter += 1
		}(repo)
	}

	// wait for all clones to finish
	wg.Wait()
	s.Stop()
	fmt.Printf("Cloned %d repos\n", len(repos))

	return nil
}
