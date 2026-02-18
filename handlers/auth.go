package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"youtubevid/auth"
	"youtubevid/db"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var userID int64
	var uuid string
	var hash string

	err := db.DB.QueryRow("SELECT id, uuid, password FROM users WHERE name=?", username).Scan(&userID, &uuid, &hash)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]any{"Error": "invalid credentials"})
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]any{"Error": "invalid credentials"})
	}

	token, _ := auth.CreateJWT(userID, username, uuid)
	c.SetCookie(auth.CreateCookie(token))
	c.Response().Header().Set("HX-Redirect", "/user/homepage")
	return c.NoContent(http.StatusOK)
}


func Signup(c echo.Context) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)

	_, err := db.DB.Exec(
			"INSERT INTO users(name, uuid, password, playlist) VALUES(?, ?, ?, ?)",
			c.FormValue("username"),
			uuid.New().String(),
			hash,
			"[]",
	)
	
	if err != nil {
		return c.Render(http.StatusBadRequest, "signup.html", map[string]any{"Error": err.Error()})
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func Query(c echo.Context) error {
	uuid := c.FormValue("input")
	return c.Redirect(303, "/user/"+uuid)
}

func Library(c echo.Context) error {
	name := c.FormValue("PlaylistName")
	if name == "" {
		return c.String(http.StatusBadRequest, "playlist name is required")
	}

	userUUID, ok := c.Get("UUID").(string)
	if !ok || userUUID == "" {
		return c.String(http.StatusUnauthorized, "user not authenticated")
	}

	id := uuid.New().String()

	// INSERT playlist
	_, err := db.DB.Exec(
		"INSERT INTO playlists(name, uuid, songs) VALUES(?, ?, ?)",
		name,
		id,
		"[]",
	)
	if err != nil {
		fmt.Println("INSERT PLAYLIST ERROR:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// UPDATE user
	_, err = db.DB.Exec(
		"UPDATE users SET playlist = ? WHERE uuid = ?",
		id,
		userUUID,
	)
	if err != nil {
		fmt.Println("UPDATE USER ERROR:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// SUCCESS
	c.Response().Header().Set("HX-Redirect", "/user/library")
	return c.NoContent(http.StatusOK)
}

