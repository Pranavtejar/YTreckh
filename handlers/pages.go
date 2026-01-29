package handlers

import (
	"net/http"
	"youtubevid/db"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.Render(200, "login.html", nil)
}

func LoginPage(c echo.Context) error {
	return c.Render(200, "login.html", nil)
}

func SignupPage(c echo.Context) error {
	return c.Render(200, "signup.html", nil)
}

func Error(c echo.Context) error {
	return c.Render(404, "error.html", nil)
}

func Homepage(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok || username == "" {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	userID, _ := c.Get("user_id").(int64)
	uuidVal, _ := c.Get("UUID").(string)

	data := map[string]any{
		"Username": username,
		"UserID":   userID,
		"UUID":     uuidVal,
	}

	return c.Render(http.StatusOK, "dashboard.html", data)
}

func DisProfile(c echo.Context) error {
    uuid := c.Param("uuid")
    details, err := db.GetDetails(uuid)
    if err != nil {
        return c.Render(404, "error.html", map[string]any{"UUID": uuid})
    }
    return c.Render(200, "profile.html", details)
}
