package main

import (
	"github.com/0xViva/portfolio-webpage/views"
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
	e.Logger.Fatal(e.Start(":3000"))
}

func homeHandler(c echo.Context) error {

	return render(c, views.Home("0xViva's Webpage"))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
