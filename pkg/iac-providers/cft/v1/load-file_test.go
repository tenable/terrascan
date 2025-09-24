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
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {
	testDataDir := "testdata"
	testFile, _ := filepath.Abs(path.Join(testDataDir, "testfile"))
	invalidFile, _ := filepath.Abs(path.Join(testDataDir, "deploy.yaml"))
	validFile, _ := filepath.Abs(path.Join(testDataDir, "templates", "s3", "deploy.template"))
	nestedFile, _ := filepath.Abs(path.Join(testDataDir, "templates", "s3", "nested.template"))
	partialWrongFile, _ := filepath.Abs(path.Join(testDataDir, "someResourcesIncorrectCftTemplate.yml"))

	validFileConfig := make(map[string][]output.ResourceConfig, 2)
	validFileConfig["aws_s3_bucket_policy"] = []output.ResourceConfig{{
		ID: "aws_s3_bucket_policy.BucketPolicy",
	}}
	validFileConfig["aws_s3_bucket"] = []output.ResourceConfig{{
		ID: "aws_s3_bucket.S3Bucket",
	}}

	nestedFileConfig := make(map[string][]output.ResourceConfig, 2)
	nestedFileConfig["aws_cloudformation_stack"] = []output.ResourceConfig{{
		ID: "aws_cloudformation_stack.myStackWithParams",
	}}
	nestedFileConfig["aws_s3_bucket"] = []output.ResourceConfig{{
		ID: "aws_s3_bucket.S3Bucket",
	}}

	partialWrongConfig := make(map[string][]output.ResourceConfig, 3)
	partialWrongConfig["aws_cognito_user_pool"] = []output.ResourceConfig{{
		ID: "aws_cognito_user_pool.goodpool",
	}}
	partialWrongConfig["aws_kinesis_stream"] = []output.ResourceConfig{{
		ID: "aws_kinesis_stream.riverstream",
	}, {
		ID: "aws_kinesis_stream.livestream",
	}}

	table := []struct {
		wantErr  error
		want     output.AllResourceConfigs
		cftv1    CFTV1
		name     string
		filePath string
		typeOnly bool
		options  map[string]interface{}
	}{
		{
			wantErr:  fmt.Errorf("unsupported extension for file %s", testFile),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "invalid extension",
			filePath: testFile,
			typeOnly: false,
		}, {
			wantErr:  fmt.Errorf("unable to read file nonexistent.txt"),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "nonexistent file",
			filePath: "nonexistent.txt",
			typeOnly: false,
		}, {
			wantErr:  fmt.Errorf("error while unmarshalling yaml, error %w", fmt.Errorf("yaml: line 28: did not find expected alphabetic or numeric character")),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "invalid file",
			filePath: invalidFile,
			typeOnly: false,
		}, {
			wantErr:  nil,
			want:     validFileConfig,
			cftv1:    CFTV1{},
			name:     "valid file",
			filePath: validFile,
			typeOnly: false,
		}, {
			wantErr:  nil,
			want:     nestedFileConfig,
			cftv1:    CFTV1{},
			name:     "nested file",
			filePath: nestedFile,
			typeOnly: false,
		},
		{
			wantErr:  nil,
			want:     partialWrongConfig,
			cftv1:    CFTV1{},
			name:     "partially wrong cft file",
			filePath: partialWrongFile,
			typeOnly: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotErr := tt.cftv1.LoadIacFile(tt.filePath, tt.options)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%+v', wantErr: '%+v'", gotErr, tt.wantErr)
			}

			if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}

			if err := verifyResponseObject(tt.want, gotResponse); err != nil {
				t.Error(err)
			}
		})
	}
}

func verifyResponseObject(want, got output.AllResourceConfigs) error {
	if len(got) == 0 && len(want) == 0 {
		return nil
	}

	if len(want) != len(got) {
		return fmt.Errorf("incorrect response object array length; got: '%d', want: '%d'", len(got), len(want))
	}

	for wantKey := range want {

		if _, ok := got[wantKey]; !ok {
			return fmt.Errorf("wanted resource object not found in response object")
		}

		for wantIndex := range want[wantKey] {

			flag := false
			for gotIndex := range got[wantKey] {
				if got[wantKey][gotIndex].ID == want[wantKey][wantIndex].ID {
					flag = true
				}
			}
			if !flag {
				return fmt.Errorf("resource ID mismatch for key: '%s' at index: '%d'", wantKey, wantIndex)
			}

		}
	}

	return nil
}
