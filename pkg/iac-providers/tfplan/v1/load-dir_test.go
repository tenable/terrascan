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

package tfplan

import (
	"reflect"
	"testing"
)

func TestLoadIacDir(t *testing.T) {

	t.Run("directory not supported", func(t *testing.T) {
		var (
			dirPath = "some-path"
			tfplan  = TFPlan{}
			wantErr = errIacDirNotSupport
			options = make(map[string]interface{})
		)
		_, err := tfplan.LoadIacDir(dirPath, options)
		if !reflect.DeepEqual(wantErr, err) {
			t.Errorf("error want: '%v', got: '%v'", wantErr, err)
		}
	})
}
