package sessions

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/sessions", CreateSessionHandler)
	r.GET("/sessions", GetSessionsHandler)
	r.POST("/sessions/current/clear", ClearCurrentSessionHandler)
	r.GET("/sessions/:id", GetSessionByIDHandler)
	r.PATCH("/sessions/:id/end", EndSessionHandler)
	r.DELETE("/sessions/:id", DeleteSessionHandler)
}
