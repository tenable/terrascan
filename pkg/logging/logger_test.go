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
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestGetLoggerLevel(t *testing.T) {
	table := []struct {
		name  string
		input string
		want  zapcore.Level
	}{
		{
			name:  "empty log level",
			input: "",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "invalid log level",
			input: "some log level",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "debug log level",
			input: "debug",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "panic log level",
			input: "panic",
			want:  zapcore.PanicLevel,
		},
	}

	for _, tt := range table {
		got := getLoggerLevel(tt.input)
		if got != tt.want {
			t.Errorf("got: '%v', want: '%v'", got, tt.want)
		}
	}
}

func TestGetLogger(t *testing.T) {

	table := []struct {
		name          string
		logLevel      string
		encoding      string
		encodingLevel func(zapcore.Level, zapcore.PrimitiveArrayEncoder)
	}{
		{
			name:          "check debug log level",
			logLevel:      "debug",
			encoding:      "json",
			encodingLevel: zapcore.LowercaseLevelEncoder,
		},
		{
			name:          "check log level",
			logLevel:      "panic",
			encoding:      "console",
			encodingLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	for _, tt := range table {
		got := GetLogger(tt.logLevel, tt.encoding, tt.encodingLevel)
		if ce := got.Check(getLoggerLevel(tt.logLevel), "testing"); ce == nil {
			t.Errorf("unexpected error")
		}
	}
}

func TestGetDefaultLogger(t *testing.T) {
	t.Run("json encoding", func(t *testing.T) {
		Init("json", "info")
		got := GetDefaultLogger()
		want := globalLogger
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
	t.Run("console encoding", func(t *testing.T) {
		Init("console", "info")
		got := GetDefaultLogger()
		want := globalLogger
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}
