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
	userIDVal := c.Get("user_id")
	uuidValRaw := c.Get("UUID")

	userID, ok1 := userIDVal.(int64)
	uuidVal, ok2 := uuidValRaw.(string)

	if !ok1 || !ok2 {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

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

//MAKE ALL THE EXISTING PLAYLISTS REGISTERED TO THIS ACCOUNT SHOW UP ON THE PAGE 
func LibraryPage(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, name, uuid, songs FROM playlists")
	if err != nil{
		return c.Render(404, "error.html", nil)
	}
	defer rows.Close()
	var playlists []map[string]any

	for rows.Next(){
		var id int
		var uuid string
		var name string
		var songs string

		if err := rows.Scan(&id, &name, &uuid, &songs); err != nil{
			return c.String(500, err.Error())
		}

		playlists = append(playlists, map[string]any{
			"id":id,
			"uuid":uuid,
			"name":name,
			"songs":songs,
		})
	}
	
	return c.Render(http.StatusOK, "library.html",map[string]any{
		"Playlists": playlists,
	})
}
