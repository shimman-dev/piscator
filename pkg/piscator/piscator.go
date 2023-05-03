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
	"sync"
)

type RepoModel struct {
	Url     string `json:"html_url"`
	Name    string `json:"name"`
	Lang    string `json:"language"`
	Private bool   `json:"private"`
	Size    int    `json:"size"`
}

type RepoCollection struct {
	Repos []*RepoModel `json:"repos"`
}

func GetRepos(name string, isOrg, isPrivate bool) string {
	var githubURL string

	gh, err := url.Parse("https://api.github.com/")
	if err != nil {
		log.Fatal(err)
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

	res, err := http.Get(githubURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var repos []RepoModel
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		log.Fatal(err)
	}

	filteredRepos := []RepoModel{}
	for _, repo := range repos {
		if repo.Name != "" && repo.Url != "" {
			filteredRepo := RepoModel{
				Name:    repo.Name,
				Lang:    repo.Lang,
				Private: repo.Private,
				Size:    repo.Size,
				Url:     repo.Url,
			}
			filteredRepos = append(filteredRepos, filteredRepo)
		}
	}

	// Marshal filteredRepos into JSON
	jsonData, err := json.MarshalIndent(filteredRepos, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("repos.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonData)
}

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

func CloneReposFromJson(jsonStr, name string) error {
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

	// clone each repo in a separate goroutine
	for _, repo := range repos {
		go func(repo Repo) {
			defer wg.Done()

			cmd := exec.Command("git", "clone", repo.URL)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("failed to clone %s: %s\n", repo.URL, string(out))
			}

			// TODO: more succinct messaging
			log.Printf("Cloned %s into %s/%s\n", repo.URL, dir, repo.Name)
		}(repo)
	}

	// wait for all clones to finish
	wg.Wait()

	return nil
}
