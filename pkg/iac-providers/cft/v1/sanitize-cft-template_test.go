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

package cftv1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/awslabs/goformation/v7"
)

func TestCFTV1_sanitizeCftTemplate(t *testing.T) {
	type args struct {
		isYAML bool
	}
	tests := []struct {
		name      string
		inputFile string
		args      args
		wantErr   bool
	}{
		{
			name:      "input file with incorrect values in parameters",
			inputFile: filepath.Join("testdata", "incorrectTypesInParamsCftTemplate.yml"),
			args: args{
				isYAML: true,
			},
			wantErr: false,
		},
		{
			name:      "input file with incorrect values in parameters",
			inputFile: filepath.Join("testdata", "incorrectTypesInResourcesCftTemplate.yml"),
			args: args{
				isYAML: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CFTV1{}
			data, err := os.ReadFile(tt.inputFile)
			if err != nil {
				t.Error(err)
			}

			_, err = goformation.Open(tt.inputFile)
			if err == nil {
				t.Error("CFTV1.sanitizeCftTemplate() got no error, expected parsing error")
			}

			templateMap, err := a.sanitizeCftTemplate(tt.inputFile, data, tt.args.isYAML)
			if (err != nil) != tt.wantErr {
				t.Errorf("CFTV1.sanitizeCftTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			resData, err := json.Marshal(templateMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceMap marshalling error; error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_, err = goformation.ParseJSON(resData)
			if err != nil {
				t.Errorf("CFTV1.sanitizeCftTemplate() parsing error; error = %v, wantError: %v", err, nil)
			}
		})
	}
}

func Test_fixWithType(t *testing.T) {
	intVar := 1
	boolVar := true
	floatVar := 1.0
	stringVar := "test"
	stringSliceVar := []string{}
	pointerToStringSliceVar := []*string{}
	boolSliceVar := []bool{}

	type Address struct {
		City string
		PIN  *string
	}

	type Employee struct {
		Name      *string
		Age       float64
		Skills    []string
		Addr      *Address `json:"emp_address"`
		IsManager *bool
	}

	type Department struct {
		Name      string
		Count     *float64
		Employees []Employee `json:"dept_employees"`
	}

	dept := Department{}

	var invalidDeptData map[string]interface{}
	var validDeptData map[string]interface{}

	invalidDeptDataStr := []byte(`{
		"Name": "Engineering",
		"Count": "100",
		"dept_employees": [
			{
				"Name": "emp1",
				"Age": 25,
				"Skills": ["skill1", 2, 3],
				"IsManager": "true",
				"emp_address": {
					"City": "Xandar",
					"PIN": 111111
				}
			},
			{
				"Name": "emp2",
				"Age": "35",
				"Skills": ["skill1", "skill2", 3],
				"IsManager": false,
				"emp_address": {
					"City": 123,
					"PIN": "222222"
				}
			}
		],
		"Rank": 1
	}`)

	err := json.Unmarshal(invalidDeptDataStr, &invalidDeptData)
	if err != nil {
		t.Error(err)
	}

	validDeptDataStr := []byte(`{
		"Name": "Engineering",
		"Count": 100,
		"dept_employees": [
			{
				"Name": "emp1",
				"Age": 25,
				"Skills": ["skill1", "2", "3"],
				"IsManager": true,
				"emp_address": {
					"City": "Xandar",
					"PIN": "111111"
				}
			},
			{
				"Name": "emp2",
				"Age": 35,
				"Skills": ["skill1", "skill2", "3"],
				"IsManager": false,
				"emp_address": {
					"City": "123",
					"PIN": "222222"
				}
			}
		],
		"Rank": 1
	}`)

	json.Unmarshal(validDeptDataStr, &validDeptData)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		data interface{}
		r    reflect.Type
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "type of data matches expected type: int",
			args: args{
				data: 13,
				r:    reflect.TypeOf(intVar),
			},
			want: nil,
		},
		{
			name: "type of data matches expected type: string",
			args: args{
				data: "1",
				r:    reflect.TypeOf(stringVar),
			},
			want: nil,
		},
		{
			name: "type of data matches expected type: bool",
			args: args{
				data: false,
				r:    reflect.TypeOf(boolVar),
			},
			want: nil,
		},
		{
			name: "type of data matches expected type: float",
			args: args{
				data: 1.0,
				r:    reflect.TypeOf(floatVar),
			},
			want: nil,
		},
		// want int against input data
		{
			name: "want int and original data is string, data can be converted to int",
			args: args{
				data: "1",
				r:    reflect.TypeOf(intVar),
			},
			want: int64(1),
		},
		{
			name: "want int and original data is float",
			args: args{
				data: 2.0,
				r:    reflect.TypeOf(intVar),
			},
			want: 2,
		},
		{
			name: "want int and original data is string, data cannot be converted to int",
			args: args{
				data: "someValue",
				r:    reflect.TypeOf(intVar),
			},
			// we don't modify the value if it can't be converted
			want: nil,
		},
		// want float against input data
		{
			name: "want float and original data is string, data can be converted to float",
			args: args{
				data: "3.3",
				r:    reflect.TypeOf(floatVar),
			},
			want: 3.3,
		},
		{
			name: "want float and original data is string, data cannot be converted to float",
			args: args{
				data: "someStringValue",
				r:    reflect.TypeOf(floatVar),
			},
			want: nil,
		},
		{
			name: "want float and original data is int",
			args: args{
				data: 4,
				r:    reflect.TypeOf(floatVar),
			},
			want: 4.0,
		},
		// want string against input data
		{
			name: "want string and original data is int",
			args: args{
				data: 4,
				r:    reflect.TypeOf(stringVar),
			},
			want: "4",
		},
		{
			name: "want string and original data is float",
			args: args{
				data: 3.141,
				r:    reflect.TypeOf(stringVar),
			},
			want: "3.141",
		},
		{
			name: "want string and original data is boolean",
			args: args{
				data: false,
				r:    reflect.TypeOf(stringVar),
			},
			want: "false",
		},
		// want bool against input data
		{
			name: "want bool and original data is string",
			args: args{
				data: "false",
				r:    reflect.TypeOf(boolVar),
			},
			want: false,
		},
		{
			name: "want bool and original data is int",
			args: args{
				data: 3,
				r:    reflect.TypeOf(boolVar),
			},
			want: nil,
		},
		// tests for array and objects
		{
			name: "want array of string and input is array of integers",
			args: args{
				data: []interface{}{1, 2, 3},
				r:    reflect.TypeOf(stringSliceVar),
			},
			want: []interface{}{"1", "2", "3"},
		},
		{
			name: "want array of string and input is array of integers",
			args: args{
				data: []interface{}{1, 2, 3},
				r:    reflect.TypeOf(pointerToStringSliceVar),
			},
			want: []interface{}{"1", "2", "3"},
		},
		{
			name: "want array of bools and input is array of strings",
			args: args{
				data: []interface{}{"false", "true"},
				r:    reflect.TypeOf(boolSliceVar),
			},
			want: []interface{}{false, true},
		},
		{
			name: "input is map[string]interface{} with invalid data w.r.t struct fields",
			args: args{
				data: invalidDeptData,
				r:    reflect.TypeOf(dept),
			},
			want: validDeptData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixWithType(tt.args.data, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fixWithType() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_examineStruct(t *testing.T) {
	nonStructVar := "test"

	type structWithoutJSONTags struct {
		One   string
		Two   int
		Three interface{}
	}

	type structWithJSONTags struct {
		One   string      `json:"one"`
		Two   int         `json:"t,omitempty"`
		Three interface{} `json:"third_tag,omitempty"`
		Four  float64
	}

	structVar1 := structWithoutJSONTags{}
	structVar2 := structWithJSONTags{}

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name     string
		args     args
		want     map[string]reflect.StructField
		length   int
		wantKeys []string
	}{
		{
			name: "input type is not a struct",
			args: args{
				t: reflect.TypeOf(nonStructVar),
			},
			want: nil,
		},
		{
			name: "input type is a struct, struct fields don't have json tags",
			args: args{
				t: reflect.TypeOf(structVar1),
			},
			want:     nil,
			length:   3,
			wantKeys: []string{"One", "Two", "Three"},
		},
		{
			name: "input type is a struct, struct fields have json tags",
			args: args{
				t: reflect.TypeOf(structVar2),
			},
			want:     nil,
			length:   4,
			wantKeys: []string{"one", "t", "third_tag", "Four"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := examineStruct(tt.args.t)
			if got != nil {
				if tt.length != len(got) {
					t.Errorf("examineStruct() = returned map doesn't have correct length, expected %d, got %d", tt.length, len(got))
				}

				for _, key := range tt.wantKeys {
					_, ok := got[key]
					if !ok {
						t.Errorf("examineStruct() = returned map doesn't have an expected key %s", key)
					}
				}
			}
		})
	}
}
