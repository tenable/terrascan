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

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var src = map[string]interface{}{
	"floatValue":  float64(10),
	"stringValue": "random-string",
	"boolValue":   true,
	"mapValue":    map[string]interface{}{"key-one": "value", "key-two": 5},
	"sliceValue":  []interface{}{true, 3, "random-string"},
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		name, key string
		expected  float64
	}{
		{name: "valid", key: "floatValue", expected: 10},
		{name: "invalid", key: "stringValue", expected: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ToFloat64(src, test.key)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name, key string
		expected  string
	}{
		{name: "valid", key: "stringValue", expected: "random-string"},
		{name: "invalid", key: "floatValue", expected: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ToString(src, test.key)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name, key string
		expected  bool
	}{
		{name: "valid", key: "boolValue", expected: true},
		{name: "invalid", key: "floatValue", expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ToBool(src, test.key)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		name, key string
		expected  map[string]interface{}
	}{
		{
			name:     "valid",
			key:      "mapValue",
			expected: map[string]interface{}{"key-one": "value", "key-two": 5},
		},
		{
			name:     "invalid",
			key:      "floatValue",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ToMap(src, test.key)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name, key string
		expected  []interface{}
	}{
		{
			name:     "valid",
			key:      "sliceValue",
			expected: []interface{}{true, 3, "random-string"},
		},
		{
			name:     "invalid",
			key:      "floatValue",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ToSlice(src, test.key)
			assert.Equal(t, test.expected, v)
		})
	}
}
