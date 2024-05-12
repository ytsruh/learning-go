package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"ytsruh.com/frankenui/views"
)

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.GET("/", HomeHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.Home())
}
