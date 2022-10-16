package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola-synthesizer-api/src/app/config"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/usecase"
	"github.com/kujilabo/cocotola-synthesizer-api/src/lib/controller/middleware"
)

func NewRouter(adminUsecase usecase.AdminUsecase, userUsecase usecase.UserUsecase, corsConfig cors.Config, appConfig *config.AppConfig, authConfig *config.AuthConfig, debugConfig *config.DebugConfig) *gin.Engine {
	if !debugConfig.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	authMiddleware := gin.BasicAuth(gin.Accounts{
		authConfig.Username: authConfig.Password,
	})

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))
		v1.Use(authMiddleware)
		{
			admin := v1.Group("admin")
			adminHandler := NewAdminHandler(adminUsecase)
			admin.POST("find", adminHandler.FindAudio)
		}
		{
			user := v1.Group("user")
			userHandler := NewUserHandler(userUsecase)
			user.POST("synthesize", userHandler.Synthesize)
			user.GET("audio/:audioID", userHandler.FindAudioByAudioID)
		}
	}

	return router
}
