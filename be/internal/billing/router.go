package billing

import "github.com/gin-gonic/gin"

func SetupPublicRoutes(r *gin.RouterGroup) {
	r.GET("/billing/plans", PlansHandler)
}
