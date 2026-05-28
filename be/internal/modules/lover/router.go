package lover

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/lover/personalities", ListPersonalitiesHandler)
	r.GET("/lover/speaking-styles", ListSpeakingStylesHandler)
	r.GET("/lover/voices", ListVoicesHandler)
	r.GET("/lover/name-suggestions", NameSuggestionsHandler)
	r.GET("/lover/profile/me", GetMyProfileHandler)
	r.POST("/lover/profile", CreateProfileHandler)
	r.POST("/lover/profile/quick", CreateQuickPresetHandler)
	r.POST("/lover/preview-voice", PreviewVoiceHandler)
}
