package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config {
		Level: zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		DisableCaller: false,
		DisableStacktrace: false,
		Sampling: nil,
		Encoding: "json",
		EncoderConfig: encoderCfg,
		OutputPaths: []string {
			"stderr",
		},
		ErrorOutputPaths: []string {
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getegid(),
		},
	}

	logger = zap.Must(config.Build())

	zap.ReplaceGlobals(logger)
}
