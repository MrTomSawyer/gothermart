package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitLogger(env string, level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	var cfg zap.Config

	switch env {
	case "production":
		cfg = zap.NewProductionConfig()
	case "development":
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = logger.Sugar()
	return nil
}
