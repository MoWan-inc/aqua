package handler

import (
	"fmt"
	"github.com/MoWan-inc/aqua/pkg/config"
	serviceutil "github.com/MoWan-inc/aqua/pkg/service/util"
	"github.com/MoWan-inc/aqua/pkg/util/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"go.uber.org/zap"
	"time"
)

func NewServer(injector *do.Injector, config *config.ApiConfig) (*gin.Engine, error) {
	// new engin
	engine := newServer(config)

	logger := log.GetDefaultLogger()
	baseLogger := logger.Base().WithOptions(zap.AddCallerSkip(-1))
	engine.Use(ginzap.Ginzap(baseLogger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(baseLogger, true))

	// token
	engine.Use(serviceutil.TokenAuthentication(
		serviceutil.GetTokenAuth(config.Tokens, []string{}))) // 白名单为空

	groupAPI := getGroupAPI(engine, config)
	handlers := []serviceutil.APIHandler{}
	for _, h := range handlers {
		h.RegisterTo(groupAPI)
	}

	return engine, nil
}

func newServer(config *config.ApiConfig) *gin.Engine {
	engine := gin.Default()
	logger := log.GetDefaultLogger()
	baseLogger := logger.Base().WithOptions(zap.AddCallerSkip(-1))

	middleWares := gin.HandlersChain{
		ginzap.RecoveryWithZap(baseLogger, true),
		location.Default(),
		corsMiddleware(engine),
	}

	engine.Use(middleWares...)

	if config.EnablePProf {
		pprof.Register(engine, "/debug/pprof")
	}

	return engine
}

func corsMiddleware(_ *gin.Engine) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // enable cookie
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}

func getGroupAPI(e *gin.Engine, config *config.ApiConfig) *gin.RouterGroup {
	prefix := config.Prefix
	path := "api"
	if prefix != "" {
		path = fmt.Sprintf("%s/%s", path, prefix)
	}
	return e.Group(path)
}
