package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-synthesizer-api/docs"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/config"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/gateway"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/handler"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/service"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/usecase"
	"github.com/kujilabo/cocotola-synthesizer-api/src/lib/handler/middleware"
)

// @securityDefinitions.basic BasicAuth
func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	if len(*env) == 0 {
		appEnv := os.Getenv("APP_ENV")
		if len(appEnv) == 0 {
			*env = "local"
		} else {
			*env = appEnv
		}
	}

	logrus.Infof("env: %s", *env)

	go func() {
		sig := <-sigs
		logrus.Info()
		logrus.Info(sig)
		done <- true
	}()

	cfg, db, sqlDB, router, tp, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	authMiddleware := gin.BasicAuth(gin.Accounts{
		cfg.Auth.Username: cfg.Auth.Password,
	})

	rf, err := gateway.NewRepositoryFactory(ctx, db, cfg.DB.DriverName)
	if err != nil {
		panic(err)
	}
	rfFunc := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, db, cfg.DB.DriverName)
	}

	synthesizerClient := gateway.NewSynthesizerClient(cfg.Synthesizer.Key, time.Duration(cfg.Synthesizer.TimeoutSec)*time.Second)

	adminUsecase := usecase.NewAdminUsecase(rf)
	userUsecase := usecase.NewUserUsecase(db, rfFunc, synthesizerClient)

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(cfg.App.Name))
		v1.Use(middleware.NewTraceLogMiddleware(cfg.App.Name))
		v1.Use(authMiddleware)
		{
			admin := v1.Group("admin")
			adminHandler := handler.NewAdminHandler(adminUsecase)
			admin.POST("find", adminHandler.FindAudio)
		}
		{
			user := v1.Group("user")
			userHandler := handler.NewUserHandler(userUsecase)
			user.POST("synthesize", userHandler.Synthesize)
			user.GET("audio/:audioID", userHandler.FindAudioByAudioID)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = cfg.App.Name
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Swagger.Host
	docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}

	gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			done <- true
		}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exiting")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Infof("Server forced to shutdown. err: %v", err)
	}
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *gin.Engine, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// init log
	if err := config.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// tracer
	tp, err := config.InitTracerProvider(cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, xerrors.Errorf("failed to InitTracerProvider. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := config.InitDB(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, nil, xerrors.Errorf("failed to InitDB. err: %w", err)
	}

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if cfg.Debug.GinMode {
		logrus.Info("GinMode Debug is enabled")
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	} else {
		logrus.Info("GinMode Debug is disabled")
	}

	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	return cfg, db, sqlDB, router, tp, nil
}
