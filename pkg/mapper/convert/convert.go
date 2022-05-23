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

// ToFloat64 looks for key in src and converts the value to a string type.
// Returns empty string otherwise.
func ToFloat64(src map[string]interface{}, key string) float64 {
	if v, ok := src[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

// ToString looks for key in src and converts the value to a string type.
// Returns empty string otherwise.
func ToString(src map[string]interface{}, key string) string {
	if v, ok := src[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ToBool looks for key in src and converts the value to a bool type.
// Returns false otherwise.
func ToBool(src map[string]interface{}, key string) bool {
	if v, ok := src[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// ToMap looks for key in src and converts the value to a map type.
// Returns nil otherwise.
func ToMap(src map[string]interface{}, key string) map[string]interface{} {
	if v, ok := src[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// ToSlice looks for key in src and converts the value to a slice type.
// Returns nil otherwise.
func ToSlice(src map[string]interface{}, key string) []interface{} {
	if v, ok := src[key]; ok {
		if s, ok := v.([]interface{}); ok {
			return s
		}
	}
	return nil
}
