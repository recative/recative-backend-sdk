package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/recative/recative-backend-sdk/util/is_zero"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

type Config struct {
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

func Init(config Config) error {
	if is_zero.CheckComparable(config) {
		viper.SetConfigName(config.ConfigFile.Name)
		viper.SetConfigType(config.ConfigFile.Type)
		viper.AddConfigPath(config.ConfigFile.Path)
	}

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

	return nil
}

func Parse(structPointer any) error {
	if err := viper.Unmarshal(structPointer, func(config *mapstructure.DecoderConfig) {
		config.MatchName = func(mapKey, fieldName string) bool {
			mapKey = strings.ReplaceAll(mapKey, "_", "")
			fieldName = strings.ReplaceAll(fieldName, "_", "")
			return strings.EqualFold(mapKey, fieldName)
		}
		config.Squash = true
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			stringTrimHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	}); err != nil {
		return err
	}

	return nil
}

func stringTrimHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t.Kind() == reflect.String {
			return strings.Trim(data.(string), `"`), nil
		}

		// Convert it by parsing
		return data, nil
	}
}
