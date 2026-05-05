package routes

import (
	"os"
	"ps_portal/api/auth"
	"ps_portal/api/resource"
	"ps_portal/handles"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	appBasePath := os.Getenv("APP_BASE_PATH")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// router.Use(handles.RateLimit())
	router.Use(handles.CorsMiddleware())
	router.Use(handles.StrictOriginMiddleware())

	router.POST(appBasePath+"/auth/GLogin", auth.GoogleLogin)
	router.POST(appBasePath+"/auth/MSLogin", auth.MicrosoftLogin)

	router.Use(utils.JWTAuthMiddleware())
	router.GET(appBasePath+"/resources", resource.GetMyResources)
	router.GET(appBasePath+"/user/images/:userId", handles.GetProfileImage)
	router.GET(appBasePath+"/images/courses/:id", handles.GetCourseImage)
	router.GET(appBasePath+"/images/dept/:id", handles.GetDeptImage)

	router.Use(handles.ScopeMiddleware())
	// router.Use(handles.DecryptMiddleware())

	// user images
	router.GET(appBasePath+"/user/profile-image", handles.GetUserProfileImage)

	//activity
	router.GET(appBasePath+"/resources/activity", resource.GetMyActivity)
	router.GET(appBasePath+"/resources/presentation", resource.GetMyPresentationView)

	return router
}
