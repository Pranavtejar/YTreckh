package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("129301")

func CreateJWT(userID int64, username string, UUID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"username": username,
		"UUID": UUID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
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

		userID := int64(claims["user_id"].(float64))
		username := claims["username"].(string)
		uuid, _ := claims["UUID"].(string)
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("UUID", uuid)
		return next(c)
	}
}
