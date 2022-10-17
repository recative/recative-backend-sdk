package http_engine

import (
	"github.com/gin-gonic/gin"
	"github.com/recative/recative-backend-sdk/pkg/config"
	"github.com/recative/recative-backend-sdk/pkg/http_engine/middleware"
)

type Config struct {
	ServerHost string `env:"SERVER_HOST"`
	ListenAddr string `env:"LISTEN_ADDR"`
}

func Default() *gin.Engine {
	if config.Environment() == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(middleware.Logger(), middleware.Recovery())

	return app
}
