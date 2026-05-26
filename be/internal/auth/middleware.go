package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ParseRequestToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		SetUser(c, claims.UserID, claims.Email, claims.Name)
		c.Next()
	}
}

func ParseRequestToken(c *gin.Context) (*Claims, error) {
	if t := c.Query("token"); t != "" {
		return ParseToken(t)
	}

	header := c.GetHeader("Authorization")
	if header == "" {
		return nil, errMissingToken
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errBadAuthHeader
	}

	return ParseToken(parts[1])
}

var (
	errMissingToken  = &authError{"missing token"}
	errBadAuthHeader = &authError{"Authorization format: Bearer <token>"}
)

type authError struct{ msg string }

func (e *authError) Error() string { return e.msg }
