package push

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/devices", RegisterDeviceHandler)
	r.POST("/devices/heartbeat", HeartbeatHandler)
	r.POST("/push/test", TestPushHandler)
	r.DELETE("/devices/:token", RevokeDeviceHandler)
}
