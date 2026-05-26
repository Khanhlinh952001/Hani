package auth

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/auth/register", RegisterHandler)
	r.POST("/auth/login", LoginHandler)

	authed := r.Group("/auth")
	authed.Use(RequireAuth())
	authed.GET("/me", MeHandler)
	authed.PATCH("/me", PatchMeHandler)
	authed.POST("/me/avatar", UploadAvatarHandler)
}
