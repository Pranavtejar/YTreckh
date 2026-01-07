package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"youtubevid/auth"
	"youtubevid/db"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var userID int64
	var hash string

	err := db.DB.QueryRow("SELECT id, password FROM users WHERE name=?", username).Scan(&userID, &hash)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]any{"Error": "invalid credentials"})
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]any{"Error": "invalid credentials"})
	}

	token, _ := auth.CreateJWT(userID)
	c.SetCookie(auth.CreateCookie(token))
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}

func Signup(c echo.Context) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)

	_, err := db.DB.Exec("INSERT INTO users(name, password) VALUES(?, ?)", c.FormValue("username"), hash)
	if err != nil {
		return c.Render(http.StatusBadRequest, "signup.html", map[string]any{"Error": "username already exists"})
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}