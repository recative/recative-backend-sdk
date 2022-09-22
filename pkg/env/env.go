package env

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type EnvironmentType string

const (
	Debug EnvironmentType = "DEBUG"
	Test  EnvironmentType = "TEST"
	Prod  EnvironmentType = "PROD"
)

type EnvironmentConfig struct {
	Environment EnvironmentType `env:"ENVIRONMENT"`
}

func ForceParse(structPointer any) {
	err := Parse(structPointer)
	if err != nil {
		panic(err)
	}
}

func Parse(structPointer any) error {
	err := env.ParseWithFuncs(structPointer, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf("string"): func(s string) (interface{}, error) {
			// This is fit docker-compose env file parse which will save \"xxx\" in .env file
			s = strings.Trim(s, `"`)
			return s, nil
		},
	}, env.Options{
		RequiredIfNoDef: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func Environment() EnvironmentType {
	return environmentConfig.Environment
}

var environmentConfig EnvironmentConfig

func init() {
	res, err := strconv.ParseBool(os.Getenv("PARSE_DOTENV"))
	if err != nil {
		panic(fmt.Sprintf("PARSE_DOTENV %s is not bool", os.Getenv("PARSE_DOTENV")))
	}
	if res != true {
		err := godotenv.Load(".env")
		if err != nil {
			panic(fmt.Errorf(
				"error loading .env file: %w when ENVIRONMENT is %s",
				err,
				os.Getenv("ENVIRONMENT"),
			))
		}
	}

	ForceParse(&environmentConfig)

	if environmentConfig.Environment != Debug &&
		environmentConfig.Environment != Test &&
		environmentConfig.Environment != Prod {
		panic("unknown ENVIRONMENT: " + environmentConfig.Environment)
	}
}

type RuntimeEnvironmentConfig struct {
	GoEnv string `env:"ENVIRONMENT"`
}
