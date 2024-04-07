package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Environment       string `json:"environment" mapstructure:"environment"`
	ServiceName       string `json:"service_name" mapstructure:"service_name"`
	Level             string `json:"level" mapstructure:"level"`
	Encoding          string `json:"encoding" mapstructure:"encoding"`
	DisableStacktrace bool   `json:"disable_stacktrace" mapstructure:"disable_stacktrace"`
}

func InitLogger(conf *Config) error {
	// Custom time format
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level, err := zapcore.ParseLevel(conf.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = conf.Encoding
	loggerConfig.EncoderConfig = encoderConfig
	loggerConfig.DisableStacktrace = conf.DisableStacktrace
	loggerConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return fmt.Errorf("failed to init logger: %v", err)
	}

	logger = logger.With(zap.String("env", conf.Environment), zap.String("service_name", conf.ServiceName))
	zap.ReplaceGlobals(logger)
	return nil
}
