package main

import (
	"html/template"
	"io"

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
	e.GET("/signup", handlers.SignupPage)
	e.POST("/login", handlers.Login)
	e.POST("/signup", handlers.Signup)

	g := e.Group("/dashboard")
	g.Use(auth.AuthMiddleware)
	g.GET("", handlers.Dashboard)

	e.Start(":8080")
}
