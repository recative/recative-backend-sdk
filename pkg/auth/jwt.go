package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/recative/recative-service-sdk/pkg/http_engine/http_err"
	"github.com/recative/recative-service-sdk/pkg/logger"
	"go.uber.org/zap"
)

type Authable interface {
	GenJwt(mapClaims jwt.MapClaims) string
	ParseJwt(tokenStr string, structPointer any) error
	ParseJwtToMap(tokenStr string) (jwt.MapClaims, error)
}

type authable struct {
	JwtSecret string
}

type Config struct {
	JwtSecret string `env:"JWT_SECRET"`
}

func New(config Config) Authable {
	return authable{
		JwtSecret: config.JwtSecret,
	}
}

func (a authable) GenJwt(mapClaims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		logger.Panic("error when sign jwt", zap.Error(err))
	}
	return tokenString
}

// ParseJwt returns user id or error if any occurs.
func (a authable) ParseJwtToMap(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, http_err.Unauthorized.New("invalid token header")
	}
	token, err := (&jwt.Parser{UseJSONNumber: true}).Parse(tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JwtSecret), nil
		})
	if err != nil {
		return nil, http_err.Unauthorized.New("invalid JWT: parsing failed" + err.Error())
	}

	return token.Claims.(jwt.MapClaims), nil
}

func (a authable) ParseJwt(tokenStr string, structPointer any) error {
	if tokenStr == "" {
		return http_err.Unauthorized.New("invalid token header")
	}
	token, err := (&jwt.Parser{UseJSONNumber: true}).Parse(tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JwtSecret), nil
		})
	if err != nil {
		return http_err.Unauthorized.New("invalid JWT: parsing failed" + err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := mapstructure.Decode(claims, structPointer); err != nil {
			return http_err.Unauthorized.New("invalid JWT: parsing failed" + err.Error())
		}
		return nil
	}

	return http_err.Unauthorized.New("invalid JWT: parsing failed")
}
