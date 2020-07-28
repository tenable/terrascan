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

package utils

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	validJSONFile = "./testdata/valid.json"
)

var (
	validJSON      []byte
	validJSONInput = map[string]int{"apple": 5, "lettuce": 7}
)

func init() {
	validJSON, _ = ioutil.ReadFile(validJSONFile)

}

func TestPrintJSON(t *testing.T) {

	table := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "empty JSON",
			input: make(map[string]interface{}),
			want:  "{}",
		},
		{
			name:  "valid JSON",
			input: validJSONInput,
			want:  string(validJSON),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := &bytes.Buffer{}
			PrintJSON(tt.input, got)
			if strings.TrimSpace(got.String()) != strings.TrimSpace(tt.want) {
				t.Errorf("got:\n'%v'\n, want:\n'%v'\n", got, tt.want)
			}
		})
	}
}
