package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/recative/recative-service-sdk/pkg/logger"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				now := time.Now().Unix()
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Int64("timestamp", now),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", stack(3)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.Error("Recovery from panic",
					zap.Int64("timestamp", now),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", stack(3)),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
