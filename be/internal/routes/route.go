package routes

import (
	"be/internal/admin"
	"be/internal/auth"
	"be/internal/billing"
	"be/internal/modules/characters"
	"be/internal/modules/lover"
	"be/internal/modules/memories"
	"be/internal/modules/messages"
	"be/internal/modules/push"
	"be/internal/modules/sessions"
	"be/internal/stt"
	"be/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")

	auth.SetupRoutes(api)
	billing.SetupPublicRoutes(api)

	protected := api.Group("")
	protected.Use(auth.RequireAuth())

	// WebSocket + guest preview chat
	websocket.SetupRoutes(protected)

	registered := protected.Group("")
	registered.Use(auth.RequireRegistered())
	characters.SetupRoutes(registered)
	lover.SetupRoutes(registered)
	sessions.SetupRoutes(registered)
	messages.SetupRoutes(registered)
	memories.SetupRoutes(registered)
	push.SetupRoutes(registered)
	admin.SetupRoutes(registered)
	registered.POST("/soniox/temporary-key", stt.TemporaryKeyHandler)
}
