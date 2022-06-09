/*
    Copyright (C) 2022 Tenable, Inc.

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
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFileName                   = "terrascan.log"
	terrascanZapSink              = "terrascan-winfile-sink"
	terrascanZapSinkWithSeparator = terrascanZapSink + ":///"
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
func Init(encoding, level, logDir string) {

	// select encoding level
	encodingLevel := zapcore.LowercaseColorLevelEncoder
	if encoding == "json" {
		encodingLevel = zapcore.LowercaseLevelEncoder
	}

	// get logger
	logger := GetLogger(level, encoding, logDir, encodingLevel)

	// set global Logger as well
	globalLogger = logger.Sugar()

	// initialize global logger
	zap.ReplaceGlobals(logger)
}

// GetLogger creates a customer zap logger
func GetLogger(logLevel, encoding, logDir string, encodingLevel func(zapcore.Level, zapcore.PrimitiveArrayEncoder)) *zap.Logger {

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

	if logDir != "" {
		logDirPath := filepath.Join(logDir, logFileName)
		// currently zap's default file sink do not support handling of windows files
		// so if we are on windows os register an sink which will handle the opening of file.
		if utils.IsWindowsPlatform() {
			zap.RegisterSink(terrascanZapSink, newTerrascanWinFileSink)
			logDirPath = terrascanZapSinkWithSeparator + logDirPath
		}
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, logDirPath)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal("failed to initialize logger %w", err)
	}

	return logger
}

// GetDefaultLogger returns the globalLogger
func GetDefaultLogger() *zap.SugaredLogger {
	return globalLogger
}

// newTerrascanWinFileSink custom sink to open file on windows
// https://github.com/uber-go/zap/issues/621#issue-350197709 - referred from this issue comment
func newTerrascanWinFileSink(u *url.URL) (zap.Sink, error) {
	// copy pasting this error from default FileSink of zap to keep the standard behaviour across all platforms
	// newFileSink()
	if u.User != nil {
		return nil, fmt.Errorf("user and password not allowed with file URLs: got %v", u)
	}
	if u.Fragment != "" {
		return nil, fmt.Errorf("fragments not allowed with file URLs: got %v", u)
	}
	if u.RawQuery != "" {
		return nil, fmt.Errorf("query parameters not allowed with file URLs: got %v", u)
	}
	// Error messages are better if we check hostname and port separately.
	if u.Port() != "" {
		return nil, fmt.Errorf("ports not allowed with file URLs: got %v", u)
	}
	if hn := u.Hostname(); hn != "" && hn != "localhost" {
		return nil, fmt.Errorf("file URLs must leave host empty or use localhost: got %v", u)
	}
	// after url.Parse() u.Path contains and extra leading slash
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}
