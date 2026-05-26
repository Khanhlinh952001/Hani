package messages

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/messages", CreateMessageHandler)
	r.GET("/messages", GetMessagesHandler)
	r.GET("/messages/:id", GetMessageByIDHandler)
	r.DELETE("/messages/:id", DeleteMessageHandler)
}
