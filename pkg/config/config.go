package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Environment  EnvironmentType
	ConfigFile   ConfigFile
	AutoParseEnv bool
	// WeakMatchName ignore difference between camelCase, snake_case, etc.
	WeakMatchName bool
}

type ConfigFile struct {
	Name string
	Type string
	Path string
}

var _config Config

func init() {
	viper.AutomaticEnv()
}

func defaultConfig(opts ...ConfigOption) *Config {
	viper.SetDefault("ENVIRONMENT", string(Debug))
	viper.SetDefault("CONFIG_FILE_NAME", "config")
	viper.SetDefault("CONFIG_FILE_TYPE", "yaml")
	viper.SetDefault("CONFIG_FILE_PATH", ".")
	viper.SetDefault("AUTO_PARSE_ENV", true)
	viper.SetDefault("WEAK_MATCH_NAME", true)

	config := &Config{
		Environment: EnvironmentType(viper.GetString("ENVIRONMENT")),
		ConfigFile: ConfigFile{
			Name: viper.GetString("CONFIG_FILE_NAME"),
			Type: viper.GetString("CONFIG_FILE_TYPE"),
			Path: viper.GetString("CONFIG_FILE_PATH"),
		},
		AutoParseEnv:  viper.GetBool("AUTO_PARSE_ENV"),
		WeakMatchName: viper.GetBool("WEAK_MATCH_NAME"),
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

type ConfigOption func(*Config)

func Init(opts ...ConfigOption) error {
	config := defaultConfig(opts...)

	viper.SetConfigName(config.ConfigFile.Name)
	viper.SetConfigType(config.ConfigFile.Type)
	viper.AddConfigPath(config.ConfigFile.Path)

	if config.AutoParseEnv {
		for _, v := range os.Environ() {
			res := strings.Split(v, "=")
			err := viper.BindEnv(res[0])
			if err != nil {
				return err
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if config.Environment != Debug && config.Environment != Test && config.Environment != Prod {
		panic("invalid environment: " + string(config.Environment))
	}

	_config = *config

	return nil
}

func Parse(structPointer any) error {
	if err := viper.Unmarshal(structPointer, defaultDecoderConfig); err != nil {
		return err
	}

	return nil
}

func ForceParse(structPointer any) {
	err := Parse(structPointer)
	if err != nil {
		panic(err)
	}
}

func ParseByKey(key string, structPointer any) error {
	if err := viper.UnmarshalKey(key, structPointer, defaultDecoderConfig); err != nil {
		return err
	}

	return nil
}

func ForceParseByKey(key string, structPointer any) {
	err := ParseByKey(key, structPointer)
	if err != nil {
		panic(err)
	}
}

func defaultDecoderConfig(config *mapstructure.DecoderConfig) {
	if _config.WeakMatchName {
		config.MatchName = func(mapKey, fieldName string) bool {
			mapKey = strings.ReplaceAll(mapKey, "_", "")
			fieldName = strings.ReplaceAll(fieldName, "_", "")
			return strings.EqualFold(mapKey, fieldName)
		}
	}
	config.Squash = true
	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		stringTrimHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
}

func stringTrimHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t.Kind() == reflect.String {
			return strings.Trim(data.(string), `"`), nil
		}

		return data, nil
	}
}
