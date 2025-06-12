package github

import (
	"encoding/json"
	"fmt"
	"github.com/0xViva/webpage/components"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	repoWg      sync.WaitGroup
	mu          sync.Mutex
	githubRepos []components.GitHubRepo
)

func GetLatestRepos(token string) ([]components.GitHubRepo, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	author := "0xViva"
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.github.com/user/repos?sort=pushed&direction=desc&per_page=10&author=%s", author), nil)
	if err != nil {
		log.Printf("failed to create repos request: %v", err)
		return nil, fmt.Errorf("failed to create repos request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "personal-webpage")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to fetch repos: %v", err)
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("github api returned status code %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("github api returned status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var repos []struct {
		Name        string    `json:"name"`
		FullName    string    `json:"full_name"`
		HTMLURL     string    `json:"html_url"`
		Description string    `json:"description"`
		UpdatedAt   time.Time `json:"updated_at"`
		Visibility  string    `json:"visibility"`
		Owner       struct {
			Login string `json:"login"`
		} `json:"owner"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Printf("failed to decode repos response: %v", err)
		return nil, fmt.Errorf("failed to decode repos response: %w", err)
	}

	githubRepos := make([]components.GitHubRepo, 0, len(repos))
	for _, repo := range repos {
		repo := repo // capture loop variable
		repoWg.Add(1)
		go func() {
			defer repoWg.Done()

			// Fetch branches
			branchesReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", repo.Owner.Login, repo.Name), nil)
			branchesReq.Header.Set("Authorization", "Bearer "+token)

			branchesResp, err := client.Do(branchesReq)
			if err != nil {
				log.Printf("failed to fetch branches for %s: %v", repo.Name, err)
				return
			}
			defer branchesResp.Body.Close()

			var branches []struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(branchesResp.Body).Decode(&branches); err != nil {
				log.Printf("failed to decode branches for %s: %v", repo.Name, err)
				return
			}

			var (
				branchWg        sync.WaitGroup
				commitsMu       sync.Mutex
				enrichedCommits []components.GitHubCommit
				commitsSeen     = make(map[string]struct{})
			)

			for _, branch := range branches {
				branch := branch
				branchWg.Add(1)
				go func() {
					defer branchWg.Done()

					// Fetch commits
					commitReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&per_page=6&author=%s", repo.Owner.Login, repo.Name, branch.Name, author), nil)
					commitReq.Header.Set("Authorization", "Bearer "+token)

					commitResp, err := client.Do(commitReq)
					if err != nil {
						log.Printf("failed to fetch commits for %s on %s: %v", repo.Name, branch.Name, err)
						return
					}
					defer commitResp.Body.Close()

					var baseCommits []struct {
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
					if err := json.NewDecoder(commitResp.Body).Decode(&baseCommits); err != nil {
						log.Printf("failed to decode commits for %s on %s: %v", repo.Name, branch.Name, err)
						return
					}

					for _, base := range baseCommits {
						if base.Commit.Author.Name != author {
							continue
						}

						// Fetch commit details
						detailsReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", repo.Owner.Login, repo.Name, base.SHA), nil)
						detailsReq.Header.Set("Authorization", "Bearer "+token)

						detailsResp, err := client.Do(detailsReq)
						if err != nil {
							log.Printf("failed to fetch commit detail for %s: %v", base.SHA, err)
							continue
						}
						defer detailsResp.Body.Close()

						var detail struct {
							Stats struct {
								Additions int `json:"additions"`
								Deletions int `json:"deletions"`
							} `json:"stats"`
						}
						if err := json.NewDecoder(detailsResp.Body).Decode(&detail); err != nil {
							log.Printf("failed to decode commit detail for %s: %v", base.SHA, err)
							continue
						}

						commit := components.GitHubCommit{
							SHA:       base.SHA,
							HTMLURL:   base.HTMLURL,
							Message:   base.Commit.Message,
							Additions: detail.Stats.Additions,
							Deletions: detail.Stats.Deletions,
							Author:    base.Commit.Author,
						}

						commitsMu.Lock()
						if _, exists := commitsSeen[base.SHA]; !exists {
							enrichedCommits = append(enrichedCommits, commit)
							commitsSeen[base.SHA] = struct{}{}
						}
						commitsMu.Unlock()
					}
				}()
			}

			branchWg.Wait()

			if len(enrichedCommits) == 0 {
				return
			}

			// Update UpdatedAt from commits
			for _, c := range enrichedCommits {
				if c.Author.Date.After(repo.UpdatedAt) {
					repo.UpdatedAt = c.Author.Date
				}
			}

			sort.Slice(enrichedCommits, func(i, j int) bool {
				return enrichedCommits[i].Author.Date.After(enrichedCommits[j].Author.Date)
			})

			if len(enrichedCommits) > 6 {
				enrichedCommits = enrichedCommits[:6]
			}

			mu.Lock()
			githubRepos = append(githubRepos, components.GitHubRepo{
				Name:        repo.Name,
				HTMLURL:     repo.HTMLURL,
				Description: repo.Description,
				UpdatedAt:   repo.UpdatedAt,
				Visibility:  repo.Visibility,
				Owner:       repo.Owner,
				Commits:     enrichedCommits,
			})
			mu.Unlock()
		}()
	}

	repoWg.Wait()
	sort.Slice(githubRepos, func(i, j int) bool {
		return githubRepos[i].UpdatedAt.After(githubRepos[j].UpdatedAt)
	})
	if len(githubRepos) > 3 {
		githubRepos = githubRepos[:3]
	}
	return githubRepos, nil
}
