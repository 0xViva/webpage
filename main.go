package main

import (
	"encoding/json"
	"fmt"
	"github.com/0xViva/webpage/components"
	"github.com/0xViva/webpage/views"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	toEmail       string
	EmailPassword string
	githubToken   string
)

func main() {
	e := echo.New()

	godotenv.Load()

	toEmail = os.Getenv("TO_EMAIL")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
	githubToken = os.Getenv("GITHUB_TOKEN")

	e.Use(middleware.Logger())
	e.Static("/style", "style")
	e.Static("/assets", "assets")

	e.GET("/", homeHandler)
	e.GET("/form", formHandler)
	e.POST("/contact", contactHandler)
	e.GET("/fetch-repos", fetchReposHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
	host := c.Request().Host
	title := getNameFromDomain(host) + "'s Website"
	name := "August"

	return render(c, views.Home(title, name))

}
func fetchReposHandler(c echo.Context) error {
	repos, err := getLatestRepos(githubToken)
	if err != nil {
		c.Logger().Errorf("Failed to fetch GitHub repos: %v", err)
		return render(c, components.RepoContainer(nil))
	}
	return render(c, components.RepoContainer(repos))
}
func formHandler(c echo.Context) error {
	return render(c, components.Form())
}

func contactHandler(c echo.Context) error {
	host := c.Request().Host
	name := c.FormValue("name")
	email := c.FormValue("email")
	company := c.FormValue("company")
	projectType := c.FormValue("project-type")
	budget := c.FormValue("budget")
	timeline := c.FormValue("timeline")
	message := c.FormValue("message")

	emailBody := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
            <h2 style="color: #333; border-bottom: 2px solid #eee; padding-bottom: 10px;">Dev Request</h2>
            
            <div style="background: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
                <h3 style="color: #2563eb; margin-top: 0;">Contact Information</h3>
                <p><strong>Name:</strong> %s</p>
                <p><strong>Email:</strong> %s</p>
                <p><strong>Company:</strong> %s</p>
            </div>

            <div style="background: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
                <h3 style="color: #2563eb; margin-top: 0;">Project Details</h3>
                <p><strong>Type:</strong> %s</p>
                <p><strong>Budget Range:</strong> %s</p>
                <p><strong>Timeline:</strong> %s</p>
            </div>

            <div style="background: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
                <h3 style="color: #2563eb; margin-top: 0;">Project Description</h3>
                <p style="white-space: pre-wrap;">%s</p>
            </div>

            <div style="color: #666; font-size: 12px; margin-top: 20px; padding-top: 10px; border-top: 1px solid #eee;">
                <p>Submitted via %s</p>
                <p>Generated on %s</p>
            </div>
        </div>
    `, name, email, company, projectType, budget, timeline, message, host, time.Now().Format("January 2, 2006 at 15:04 MST"))

	m := gomail.NewMessage()
	hostname := getNameFromDomain(host)
	if hostname == "August" {
		hostname = "augustg"
	}
	m.SetHeader("From", fmt.Sprintf("%s <%s>", hostname+".dev", toEmail))
	m.SetHeader("To", fmt.Sprintf("%s <%s>", hostname+".dev", toEmail))
	m.SetAddressHeader("Cc", email, name)
	m.SetHeader("Subject", fmt.Sprintf("Project Proposal from %s - %s", name, projectType))

	domain := strings.Split(toEmail, "@")[1]
	messageID := fmt.Sprintf("<%d.project-request@%s>", time.Now().UnixNano(), domain)
	m.SetHeader("Message-ID", messageID)

	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, toEmail, EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return render(c, components.Failed("Sorry, there was an error sending your request. Please try again later."))
	}

	return render(c, components.Submitted(name, email, company, projectType, budget, timeline, message))
}

func getLatestRepos(token string) ([]components.GitHubRepo, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET",
		"https://api.github.com/user/repos?affiliation=owner,collaborator,organization_member&sort=updated&direction=desc&per_page=3", nil)
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
		// Fetch all branches
		branchesReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", repo.Owner.Login, repo.Name), nil)
		branchesReq.Header.Set("Authorization", "Bearer "+token)

		branchesResp, err := client.Do(branchesReq)
		if err != nil {
			log.Printf("failed to fetch branches for %s: %v", repo.Name, err)
			continue
		}
		defer branchesResp.Body.Close()

		var branches []struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(branchesResp.Body).Decode(&branches); err != nil {
			log.Printf("failed to decode branches for %s: %v", repo.Name, err)
			continue
		}

		var enrichedCommits []components.GitHubCommit

		for _, branch := range branches {
			commitReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&per_page=3", repo.Owner.Login, repo.Name, branch.Name), nil)
			if err != nil {
				log.Printf("failed to create commits request for %s on branch %s: %v", repo.Name, branch.Name, err)
				continue
			}
			commitReq.Header.Set("Authorization", "Bearer "+token)

			commitResp, err := client.Do(commitReq)
			if err != nil {
				log.Printf("failed to fetch commits for %s on branch %s: %v", repo.Name, branch.Name, err)
				continue
			}
			defer commitResp.Body.Close()

			if commitResp.StatusCode != http.StatusOK {
				log.Printf("github api returned status code %d for commits of %s on branch %s", commitResp.StatusCode, repo.Name, branch.Name)
				continue
			}

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
				log.Printf("failed to decode commits response for %s on branch %s: %v", repo.Name, branch.Name, err)
				continue
			}

			for _, base := range baseCommits {
				if base.Commit.Author.Name != "0xViva" {
					continue
				}

				// Fetch commit stats
				commitDetailsReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", repo.Owner.Login, repo.Name, base.SHA), nil)
				if err != nil {
					log.Printf("failed to create commit detail request for %s: %v", base.SHA, err)
					continue
				}
				commitDetailsReq.Header.Set("Authorization", "Bearer "+token)

				commitDetailsResp, err := client.Do(commitDetailsReq)
				if err != nil {
					log.Printf("failed to fetch commit detail for %s: %v", base.SHA, err)
					continue
				}
				defer commitDetailsResp.Body.Close()

				if commitDetailsResp.StatusCode != http.StatusOK {
					log.Printf("github api returned status code %d for commit %s", commitDetailsResp.StatusCode, base.SHA)
					continue
				}

				var detail struct {
					Stats struct {
						Additions int `json:"additions"`
						Deletions int `json:"deletions"`
					} `json:"stats"`
				}
				if err := json.NewDecoder(commitDetailsResp.Body).Decode(&detail); err != nil {
					log.Printf("failed to decode commit detail for %s: %v", base.SHA, err)
					continue
				}

				enrichedCommits = append(enrichedCommits, components.GitHubCommit{
					SHA:       base.SHA,
					HTMLURL:   base.HTMLURL,
					Message:   base.Commit.Message,
					Additions: detail.Stats.Additions,
					Deletions: detail.Stats.Deletions,
					Author:    base.Commit.Author,
				})
			}
		}
		for _, c := range enrichedCommits {
			if c.Author.Date.After(repo.UpdatedAt) {
				repo.UpdatedAt = c.Author.Date
			}
		}

		githubRepos = append(githubRepos, components.GitHubRepo{
			Name:        repo.Name,
			HTMLURL:     repo.HTMLURL,
			Description: repo.Description,
			UpdatedAt:   repo.UpdatedAt,
			Visibility:  repo.Visibility,
			Commits:     enrichedCommits,
		})
	}
	sort.Slice(githubRepos, func(i, j int) bool {
		return githubRepos[i].UpdatedAt.After(githubRepos[j].UpdatedAt)
	})
	return githubRepos, nil
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}

func getNameFromDomain(domain string) string {

	mapping := map[string]string{
		"augustg.dev":    "August",
		"0xviva.dev":     "0xViva",
		"localhost:8080": "Localhost",
	}

	if title, exists := mapping[domain]; exists {
		return title
	}

	host := strings.Split(domain, ":")[0]
	if title, exists := mapping[host]; exists {
		return title
	}
	fmt.Printf("No title found for domain: %s\n", domain)
	return "Default Title"
}
