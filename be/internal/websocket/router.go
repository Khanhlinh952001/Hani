package websocket

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/ws/chat", HandleChat)
}
