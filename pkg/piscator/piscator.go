package piscator

import (
	"log"
	"net/url"
	"path"
)

// make a post request to the GitHub API where you pass in an organization or user name

type RepoModel struct {
	Url      string
	Name     string
	Language string
	Private  bool
	Size     byte
}

var RepoCollection []RepoModel

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
	return githubURL
}
