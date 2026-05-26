package memories

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/memories", CreateMemoryHandler)
	r.POST("/memories/search", SearchMemoriesHandler)
	r.GET("/memories", GetMemoriesHandler)
	r.GET("/memories/:id", GetMemoryByIDHandler)
	r.PUT("/memories/:id", UpdateMemoryHandler)
	r.DELETE("/memories/:id", DeleteMemoryHandler)
}
