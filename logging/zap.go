package logging

import (
	"os"
	"strings"

	"github.com/eliezedeck/gozap2seq"
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

	// Automatically use SEQ if the env var `SEQ_URL` is present
	sequrl := strings.TrimSpace(os.Getenv("SEQ_URL"))
	if sequrl != "" {
		seqtoken := strings.TrimSpace(os.Getenv("SEQ_API_TOKEN"))
		injector, err := gozap2seq.NewLogInjector(sequrl, seqtoken)
		if err != nil {
			panic(err)
		}
		loggerseq := injector.Build(loggerConfig)

		// Tee both loggers
		if loggerconsole, err := loggerConfig.Build(); err != nil {
			panic(err)
		} else {
			L = zap.New(zapcore.NewTee(loggerconsole.Core(), loggerseq.Core()),
				zap.AddCaller(),
				zap.AddStacktrace(zapcore.ErrorLevel))
		}
	} else {
		// Use normal Zap logging
		loggerconsole, err := loggerConfig.Build()
		if err != nil {
			panic(err)
		}
		L = zap.New(loggerconsole.Core(),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel))
		L.Warn("Running normally without SEQ /")
	}
}
