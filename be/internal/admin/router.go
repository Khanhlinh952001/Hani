package admin

import (
	"be/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	g := r.Group("/admin")
	g.Use(auth.RequireAuth(), RequireAdmin())

	g.GET("/stats", StatsHandler)
	g.GET("/users", ListUsersHandler)
	g.PATCH("/users/:id", PatchUserHandler)
	g.POST("/users/:id/reset-usage", ResetUserUsageHandler)
	g.DELETE("/users/:id", DeleteUserHandler)
	g.GET("/users/:id/sessions", ListUserSessionsHandler)
	g.GET("/users/:id/memories", ListUserMemoriesHandler)
	g.DELETE("/users/:id/memories", ClearUserMemoriesHandler)
	g.POST("/users/:id/clear-conversation", ClearUserConversationHandler)
	g.GET("/sessions/:sessionId/messages", ListSessionMessagesHandler)
}
