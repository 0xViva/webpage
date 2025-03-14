package main

import (
	"strings"

	"github.com/0xViva/webpage/views"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Serve static files
	e.Static("/style", "style")
	e.Static("/assets", "assets")

	// Routes
	e.GET("/", homeHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
	host := c.Request().Host
	title := getTitleFromDomain(host)

	return render(c, views.Home(title))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
func getTitleFromDomain(domain string) string {

	mapping := map[string]string{
		"augustg.dev": "August's webpage",
		"0xviva.dev":  "0xViva's webpage",
		"localhost":   "Localhost",
	}

	host := strings.Split(domain, ":")[0]

	if title, exists := mapping[host]; exists {
		return title
	}
	return "Default Title"
}
