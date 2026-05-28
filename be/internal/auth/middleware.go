package auth

import (
	"net/http"
	"strconv"
	"strings"

	"be/internal/billing"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ParseRequestToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims.Guest {
			SetUser(c, claims)
			c.Next()
			return
		}
		if claims.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if u, err := users.GetUserByIDService(strconv.Itoa(claims.UserID)); err == nil {
			claims.Plan = billing.PlanForUser(u)
		}
		SetUser(c, claims)
		c.Next()
	}
}

// RequireRegistered blocks guest tokens from registered-only routes.
func RequireRegistered() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsGuest(c) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "registration_required"})
			return
		}
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
