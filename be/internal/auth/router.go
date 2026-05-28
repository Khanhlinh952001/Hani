package auth

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/auth/register", RegisterHandler)
	r.POST("/auth/login", LoginHandler)
	r.POST("/auth/refresh", RefreshHandler)
	r.POST("/auth/guest", GuestHandler)

	authed := r.Group("/auth")
	authed.Use(RequireAuth())
	authed.GET("/me", MeHandler)
	authed.POST("/logout", LogoutHandler)

	registered := authed.Group("")
	registered.Use(RequireRegistered())
	registered.PATCH("/me", PatchMeHandler)
	registered.POST("/me/avatar", UploadAvatarHandler)
	registered.GET("/billing/usage", UsageHandler)
}
