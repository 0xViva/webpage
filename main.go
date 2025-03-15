package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0xViva/webpage/components"
	"github.com/0xViva/webpage/views"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/gomail.v2"
)

var (
	toEmail       string
	EmailPassword string
)

func main() {
	e := echo.New()

	godotenv.Load()

	toEmail = os.Getenv("TO_EMAIL")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")

	// Middleware
	e.Use(middleware.Logger())
	// Serve static files
	e.Static("/style", "style")
	e.Static("/assets", "assets")

	// Routes
	e.GET("/", homeHandler)
	e.GET("/form", formHandler)
	e.POST("/contact", contactHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
	host := c.Request().Host
	title := getNameFromDomain(host) + "'s Website"
	name := "August"

	return render(c, views.Home(title, name))
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

	// Create a professional HTML email template
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

	// Configure email message
	m := gomail.NewMessage()
	hostname := getNameFromDomain(host)
	fmt.Print(hostname)
	m.SetHeader("From", fmt.Sprintf("%s <%s>", hostname, toEmail))
	m.SetHeader("To", fmt.Sprintf("%s <%s>", hostname, toEmail))
	m.SetAddressHeader("Cc", email, name)
	m.SetHeader("Subject", fmt.Sprintf("Project Proposal from %s - %s", name, projectType))
	m.SetBody("text/html", emailBody)

	// Send email
	d := gomail.NewDialer("smtp.gmail.com", 587, toEmail, EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return render(c, components.Failed("Sorry, there was an error sending your request. Please try again later."))
	}

	return render(c, components.Submitted(name, email, company, projectType, budget, timeline, message))
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
