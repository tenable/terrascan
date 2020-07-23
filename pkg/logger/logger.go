package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// levelMap maps human readable log level to zapcore.Level
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLoggerLevel takes human readable log level (debug, info, warn, error...)
// as input and return zapcore.Level
func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// Init initializes global custom zap logger
func Init(encoding, level string) {

	encodingLevel := zapcore.LowercaseColorLevelEncoder
	if encoding == "json" {
		encodingLevel = zapcore.LowercaseLevelEncoder
	}

	// build zap config
	zapConfig := zap.Config{
		Encoding:    encoding,
		Level:       zap.NewAtomicLevelAt(getLoggerLevel(level)),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "file",
			MessageKey:   "msg",
			EncodeLevel:  encodingLevel,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := zapConfig.Build()

	// initialize global logger
	zap.ReplaceGlobals(logger)
}
