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
	"fmt"
	"reflect"
	"testing"
)

func TestJQFilterWithQuery(t *testing.T) {

	table := []struct {
		name    string
		jqQuery string
		input   []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "invalid json",
			jqQuery: "",
			input:   []byte{},
			want:    []byte{},
			wantErr: fmt.Errorf("failed to decode input JSON. error: 'unexpected end of JSON input'"),
		},
		{
			name:    "invalid query",
			jqQuery: "am invalid",
			input:   []byte("{}"),
			want:    []byte(""),
			wantErr: fmt.Errorf("failed to parse jq query. error: 'unexpected token \"invalid\"'"),
		},
		{
			name:    "jq error",
			jqQuery: "def f: f; f, f",
			input:   []byte("{}"),
			want:    []byte(""),
			wantErr: nil,
		},
		{
			name:    "simple query",
			jqQuery: ".foo",
			input:   []byte("{\"foo\": 128}"),
			want:    []byte("128"),
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JQFilterWithQuery(tt.jqQuery, tt.input)

			if string(got) != string(tt.want) {
				t.Errorf("got: '%v', want: '%v'", string(got), string(tt.want))
			}

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", err, tt.wantErr)
			}
		})
	}
}
