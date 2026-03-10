// REFERENCE: userID is an int64, username is a string. 
// Fixed: .Error() conversion, variable typos, unassigned Split result, and safe float64-to-int64 type assertion for JWT claims.

package auth

import (
	"net/http"
	"strings"
	"time"

	"youtubevid/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("129301")

func CreateJWT(userID int64, username string, UUID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"UUID":     UUID,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func CreateCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7,
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		claims, ok := ParseJWT(cookie.Value)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// JWT unmarshals JSON numbers as float64. 
		// We must assert to float64 first, then convert to int64 to avoid a panic.
		floatID, ok := claims["user_id"].(float64)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		userID := int64(floatID)

		username, _ := claims["username"].(string)
		uuid, _ := claims["UUID"].(string)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("UUID", uuid)
		return next(c)
	}
}

func PlaylistView(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		if uuid == "" {
			return c.String(http.StatusBadRequest, "missing uuid")
		}

		// Extracting the username string set by AuthMiddleware
		user, ok := c.Get("username").(string)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		var playlistsRaw string
		err := db.DB.QueryRow("SELECT playlist FROM users WHERE name=?", user).Scan(&playlistsRaw)
		if err != nil {
			// FIXED: err.Error() instead of string(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// FIXED: Assigned result to 'playlists' and corrected 'playlistsRaw' typo
		playlists := strings.Split(playlistsRaw, ",")

		// FIXED: Logic to compare requested uuid against the user's database list
		authorized := false
		for _, p := range playlists {
			if strings.TrimSpace(p) == uuid {
				authorized = true
				break
			}
		}

		if !authorized {
			return c.String(http.StatusForbidden, "Unauthorized")
		}

		return next(c)
	}
}
