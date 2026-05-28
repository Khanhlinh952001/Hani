package characters

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/characters", ListHandler)
	r.GET("/characters/:id/preview-voice", PreviewVoiceHandler)
	r.POST("/characters/select", SelectHandler)
}
