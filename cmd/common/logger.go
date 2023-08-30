package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(debug bool) *zap.Logger {
	cfg := BuildConfig(debug)
	cfg.DisableStacktrace = true
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return log
}

func BuildConfig(debug bool) zap.Config {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.Encoding = "console"
	return cfg
}
