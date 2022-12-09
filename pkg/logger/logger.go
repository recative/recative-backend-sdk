package logger

import (
	"github.com/recative/recative-service-sdk/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = initLogger(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func RawLogger() *zap.Logger {
	return logger
}

func initEncoderConfig() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.TimeKey = "timestamp"
	encoderConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(t.UnixMilli())
	}
	encoderConf.EncodeDuration = zapcore.MillisDurationEncoder

	return encoderConf
}

func initLogger(opts ...zap.Option) (*zap.Logger, error) {
	switch config.Environment() {
	case config.Debug:
		config := zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		config.EncoderConfig = zap.NewProductionEncoderConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return config.Build(opts...)
	case config.Test:
		config := zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		config.EncoderConfig = zap.NewProductionEncoderConfig()
		config.Encoding = "json"
		return config.Build(opts...)
	case config.Prod:
		fallthrough
	default:
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.EncoderConfig = initEncoderConfig()
		return config.Build(opts...)
	}
}
