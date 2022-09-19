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

package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

const (
	// JSONDoc type for json files
	JSONDoc = "json"
)

// LoadJSON loads a JSON file into an IacDocument struct
func LoadJSON(filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 0, 1)

	// First pass determines line number data
	currentLineNumber := 1
	{ // Limit the scope for Close()
		file, err := os.Open(filePath)
		if err != nil {
			return iacDocumentList, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			currentLineNumber++
		}

		if err = scanner.Err(); err != nil {
			return iacDocumentList, err
		}
	}

	// Second pass extracts raw data
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return iacDocumentList, err
	}

	iacDocumentList = append(iacDocumentList, &IacDocument{
		Type:      JSONDoc,
		StartLine: 1,
		FilePath:  filePath,
		EndLine:   currentLineNumber,
		Data:      fileBytes,
	})

	return iacDocumentList, nil
}

// AreEqualJSON validate if two json strings are equal
func AreEqualJSON(s1, s2 string) (bool, error) {
	return AreEqualJSONBytes([]byte(s1), []byte(s2))
}

// AreEqualJSONBytes validate if two json byte arrays are equal
func AreEqualJSONBytes(b1, b2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	errmsg := "error json unmarshalling bytes: %s. error: %v"

	var err error
	err = json.Unmarshal(b1, &o1)
	if err != nil {
		return false, fmt.Errorf(errmsg, b1, err.Error())
	}
	err = json.Unmarshal(b2, &o2)
	if err != nil {
		return false, fmt.Errorf(errmsg, b2, err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
