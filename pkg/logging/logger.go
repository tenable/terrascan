/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.SugaredLogger

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

	// select encoding level
	encodingLevel := zapcore.LowercaseColorLevelEncoder
	if encoding == "json" {
		encodingLevel = zapcore.LowercaseLevelEncoder
	}

	// get logger
	logger := GetLogger(level, encoding, encodingLevel)

	// set global Logger as well
	globalLogger = logger.Sugar()

	// initialize global logger
	zap.ReplaceGlobals(logger)
}

// GetLogger creates a customer zap logger
func GetLogger(logLevel, encoding string, encodingLevel func(zapcore.Level, zapcore.PrimitiveArrayEncoder)) *zap.Logger {

	// build zap config
	zapConfig := zap.Config{
		Encoding:    encoding,
		Level:       zap.NewAtomicLevelAt(getLoggerLevel(logLevel)),
		OutputPaths: []string{"stderr"},
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

	// create zap logger
	logger, _ := zapConfig.Build()

	return logger
}

// GetDefaultLogger returns the globalLogger
func GetDefaultLogger() *zap.SugaredLogger {
	return globalLogger
}
