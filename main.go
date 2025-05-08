package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/0xViva/webpage/components"
	"github.com/0xViva/webpage/models"
	"github.com/0xViva/webpage/views"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/gomail.v2"
	"io"
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

	e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
	host := c.Request().Host
	title := getNameFromDomain(host) + "'s Website"
	name := "August"

	repos, err := getLatestRepos("0xViva", githubToken)
	if err != nil {
		c.Logger().Errorf("Failed to fetch GitHub repos: %v", err)
		return render(c, views.Home(title, name, nil))
	}

	return render(c, views.Home(title, name, repos))
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

func getLatestRepos(username, token string) ([]models.GitHubRepo, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&direction=desc&per_page=3", nil)
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
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("github api returned status code %d, failed to read response body: %v", resp.StatusCode, readErr)
			return nil, fmt.Errorf("github api returned status code %d, failed to read response body: %w", resp.StatusCode, readErr)
		}
		log.Printf("github api returned status code %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("github api returned status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var repos []struct {
		Name        string    `json:"name"`
		HTMLURL     string    `json:"html_url"`
		Description string    `json:"description"`
		UpdatedAt   time.Time `json:"updated_at"`
		Visibility  string    `json:"visibility"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Printf("failed to decode repos response: %v", err)
		return nil, fmt.Errorf("failed to decode repos response: %w", err)
	}

	githubRepos := make([]models.GitHubRepo, 0, len(repos))
	for _, repo := range repos {
		commitReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=3", username, repo.Name), nil)
		if err != nil {
			log.Printf("failed to create commits request for %s: %v", repo.Name, err)
			return nil, fmt.Errorf("failed to create commits request for %s: %w", repo.Name, err)
		}
		commitReq.Header.Set("Authorization", "Bearer "+token)

		commitResp, err := client.Do(commitReq)
		if err != nil {
			log.Printf("failed to fetch commits for %s: %v", repo.Name, err)
			return nil, fmt.Errorf("failed to fetch commits for %s: %w", repo.Name, err)
		}
		defer commitResp.Body.Close()

		if commitResp.StatusCode != http.StatusOK {
			log.Printf("github api returned status code %d for commits of %s", commitResp.StatusCode, repo.Name)
			return nil, fmt.Errorf("github api returned status code %d for commits of %s", commitResp.StatusCode, repo.Name)
		}

		var commits []models.GitHubCommit
		if err := json.NewDecoder(commitResp.Body).Decode(&commits); err != nil {
			log.Printf("failed to decode commits response for %s: %v", repo.Name, err)
			return nil, fmt.Errorf("failed to decode commits response for %s: %w", repo.Name, err)
		}

		githubRepos = append(githubRepos, models.GitHubRepo{
			Name:        repo.Name,
			HTMLURL:     repo.HTMLURL,
			Description: repo.Description,
			UpdatedAt:   repo.UpdatedAt,
			Visibility:  repo.Visibility,
			Commits:     commits,
		})
	}

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
