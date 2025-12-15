package utils

import (
	"backend-path/constants"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func ZapLogger(env string) {
	var logLevel zapcore.LevelEnabler
	var encoderConfig zapcore.EncoderConfig

	if strings.ToLower(env) == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		logLevel = zapcore.ErrorLevel
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		logLevel = zapcore.DebugLevel
	}

	encoderConfig.EncodeCaller = nil
	encoderConfig.EncodeLevel = nil
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimestampFormat)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	Logger = logger
}