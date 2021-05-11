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

package functions

import (
	"strings"
)

// ResourceIDs is a map[ARMResource.Type]output.ResourceConfig.ID
// required for resolving the resourceId function calls in ARM templates.
var ResourceIDs = map[string]string{}

// LookUp function looks for different keywords in str
// and accordingly selects a function to call.
func LookUp(vars, params map[string]interface{}, key string) interface{} {
	switch true {
	case strings.Contains(key, "concat"):
		return Concat(vars, params, key)
	case strings.Contains(key, "toLower"):
		return ToLower(vars, params, key)
	case strings.Contains(key, "resourceId"):
		return ResourceID(vars, params, key)
	case strings.Contains(key, "parameters"):
		s := Parameters(key)
		return LookUp(vars, params, s)
	case strings.Contains(key, "variables"):
		s := Variables(key)
		return LookUp(vars, params, s)
	case strings.Contains(key, "uniqueString"):
		return UniqueString()
	default:
		if v, ok := params[key]; ok {
			if res, ok := v.(string); ok {
				return LookUp(vars, params, res)
			}
			return v
		}
		if v, ok := vars[key]; ok {
			if res, ok := v.(string); ok {
				return LookUp(vars, params, res)
			}
			return v
		}
		return key
	}
}
