package http_engine

import (
	"github.com/gin-gonic/gin"
	"github.com/recative/recative-backend-sdk/pkg/config"
	"github.com/recative/recative-backend-sdk/pkg/http_engine/middleware"
	"github.com/recative/recative-backend-sdk/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
	IsLogRequestBody bool
	ServerHost       string `env:"SERVER_HOST"`
	ListenAddr       string `env:"LISTEN_ADDR"`
}

type CustomHttpEngine struct {
	*gin.Engine
	config Config
}

func Default(_config Config) *CustomHttpEngine {
	if config.Environment() == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(middleware.Logger(_config.IsLogRequestBody), middleware.Recovery())

	return &CustomHttpEngine{
		Engine: app,
		config: _config,
	}
}

func (e *CustomHttpEngine) AddPing() {
	e.GET("/ping", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
}

func (e *CustomHttpEngine) Start() {
	err := e.Run(e.config.ListenAddr)
	if err != nil {
		logger.Fatal("http engine run failed", zap.Error(err))
	}
}
