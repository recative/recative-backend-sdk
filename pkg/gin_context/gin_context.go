package gin_context

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/recative/recative-service-sdk/pkg/auth"
)

type CustomLogic = func(claims jwt.MapClaims, c *gin.Context) error

type ContextDependence struct {
	Auther      auth.Authable
	CustomLogic CustomLogic
}

var contextDependence ContextDependence

type InternalContextDependence struct {
	AuthorizationToken string
	CustomLogic        CustomLogic
}

var internalContextDependence InternalContextDependence

func Init(_contextDependence ContextDependence) {
	contextDependence = _contextDependence
}

func InitInternalContext(_internalContextDependence InternalContextDependence) {
	internalContextDependence = _internalContextDependence
}
