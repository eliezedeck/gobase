package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	L *zap.Logger
)

func init() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan_2 15:04:05.000")
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	var err error
	if L, err = loggerConfig.Build(); err != nil {
		panic(err)
	}
}
