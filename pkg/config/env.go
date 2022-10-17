package config

type EnvironmentType string

const (
	Debug EnvironmentType = "DEBUG"
	Test  EnvironmentType = "TEST"
	Prod  EnvironmentType = "PROD"
)

type EnvironmentConfig struct {
	Environment EnvironmentType `env:"ENVIRONMENT" envDefault:"TEST"`
}

func Environment() EnvironmentType {
	return env
}

var env EnvironmentType
