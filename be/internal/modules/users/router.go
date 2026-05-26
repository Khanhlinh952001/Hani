package users

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/users", CreateUserHandler)
	r.GET("/users", GetAllUsersHandler)
	r.GET("/users/:id", GetUserByIDHandler)
	r.PUT("/users/:id", UpdateUserHandler)
	r.DELETE("/users/:id", DeleteUserHandler)
}
