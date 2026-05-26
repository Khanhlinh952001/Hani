package routes

import (
	"be/internal/admin"
	"be/internal/auth"
	"be/internal/modules/memories"
	"be/internal/modules/messages"
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

	// Public
	auth.SetupRoutes(api)

	// Protected
	protected := api.Group("")
	protected.Use(auth.RequireAuth())
	sessions.SetupRoutes(protected)
	messages.SetupRoutes(protected)
	memories.SetupRoutes(protected)
	websocket.SetupRoutes(protected)
	admin.SetupRoutes(protected)
	protected.POST("/soniox/temporary-key", stt.TemporaryKeyHandler)
}
