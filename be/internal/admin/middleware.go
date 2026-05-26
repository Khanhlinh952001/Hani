package admin

import (
	"net/http"
	"strconv"

	"be/internal/auth"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := auth.UserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user, err := users.GetUserByIDService(strconv.Itoa(userID))
		if err != nil || !users.IsAdmin(user.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}
