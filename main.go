package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"youtubevid/auth"
	"youtubevid/db"
	"youtubevid/handlers"
)

type Template struct {
	t *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.t.ExecuteTemplate(w, name, data)
}

func main() {
	db.Init()

	e := echo.New()
	e.Renderer = &Template{
		t: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.GET("/", handlers.Home)
	e.GET("/login", handlers.LoginPage)
	e.GET("/error", handlers.Error)
	e.GET("/signup", handlers.SignupPage)

	e.POST("/login", handlers.Login)
	e.POST("/signup", handlers.Signup)

	e.GET("/logout", func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		return c.Redirect(http.StatusSeeOther, "/login")
	})

	g := e.Group("/user")
	g.Use(auth.AuthMiddleware)

	// FIXED ORDER (important)
	g.GET("/homepage", handlers.Homepage)
	g.GET("/library", handlers.LibraryPage)
	g.GET("/:uuid", handlers.DisProfile)

	g.POST("/homepage", handlers.Query)
	g.POST("/library", handlers.Library)

	e.Logger.Fatal(e.Start(":8080"))
}

