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

package functions

// GetString function validates and returns string from string pointer
func GetString(sptr *string) string {
	if sptr == nil {
		return ""
	}
	return *sptr
}

// GetBool function validates and returns bool from bool pointer
func GetBool(bptr *bool) bool {
	if bptr == nil {
		return false
	}
	return *bptr
}

// GetNum function validates and returns a number from int | float32 | float64 pointer
func GetNum[T int | float32 | float64](nptr *T) T {
	if nptr == nil {
		return 0
	}
	return *nptr
}
