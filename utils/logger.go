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
		logLevel = zapcore.InfoLevel
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		logLevel = zapcore.DebugLevel
	}

	encoderConfig.EncodeCaller = nil
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimestampFormat)

	var encoder zapcore.Encoder
	if strings.ToLower(env) == "production" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	Logger = zap.New(core, zap.AddCaller())
}
