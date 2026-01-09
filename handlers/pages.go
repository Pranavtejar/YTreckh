package handlers

import "github.com/labstack/echo/v4"

func Home(c echo.Context) error {
	return c.Render(200, "login.html", nil)
}

func LoginPage(c echo.Context) error {
	return c.Render(200, "login.html", nil)
}

func SignupPage(c echo.Context) error {
	return c.Render(200, "signup.html", nil)
}

func Dashboard(c echo.Context) error {
	userID := c.Get("username").(string)
	uuid := c.Get("UUID").(string)
	return c.Render(200, "dashboard.html", map[string]any{
		"UserID" : userID,
		"UUID" : uuid,
	})
}
