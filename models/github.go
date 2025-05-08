package models

import "time"

type GitHubRepo struct {
	Name        string         `json:"name"`
	HTMLURL     string         `json:"html_url"`
	Description string         `json:"description"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Visibility  string         `json:"visibility"`
	Commits     []GitHubCommit `json:"commits"`
}

type GitHubCommit struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Message string `json:"message"`
		Author  struct {
			Name string    `json:"name"`
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}
