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

import (
	"encoding/json"
)

// GetVal function pointer of any type, validates and returns value
func GetVal[T any](ptr *T) T {
	if ptr == nil {
		var retval T
		return retval
	}
	return *ptr
}

// interfaceToStruct takes the input data and populates the input struct.
func interfaceToStruct(sourceData, structPtr interface{}) error {
	byteData, err := json.Marshal(sourceData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, structPtr)
	if err != nil {
		return err
	}

	return nil
}

// PatchAWSTags returns tags in map[string]string format instead key="" and value= ""
func PatchAWSTags(sourceData interface{}) interface{} {

	type awsTags struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var tagsList []awsTags
	err := interfaceToStruct(sourceData, &tagsList)
	if err != nil {
		if valueRawMap, ok := sourceData.(map[string]interface{}); ok {
			return valueRawMap
		}
		return sourceData
	}

	// build the tags map to mimic terraform's aws tag format
	tfTagsMap := make(map[string]string)
	for i := range tagsList {
		if len(tagsList[i].Key) != 0 {
			tfTagsMap[tagsList[i].Key] = tagsList[i].Value
		}
	}

	return tfTagsMap
}

// GetBoolValueFromString returns boolean value from string
func GetBoolValueFromString(inputStr string) bool {
	if inputStr == "true" {
		return true
	} else if inputStr == "false" {
		return false
	}
	return false

}
