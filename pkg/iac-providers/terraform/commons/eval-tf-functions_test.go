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

package commons

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestEvaluateTemplatefile(t *testing.T) {
	/*
		working - path
		1 - templatefile(filepath)
		2 - templatefile(path.root+filepath)
		3 - templatefile(path.module+filepath)
		4 - templatefile(path.cwd+filepath)

		working - path - variables
		5 - templatefile(filepath, variable)
		6 - templatefile(filepath, variable+extra_variable)

		7 - templatefile(path.root+filepath, variable)
		8 - templatefile(path.root+filepath, variable+extra_variable)

		9  - templatefile(path.module+filepath, variable)
		10 - templatefile(path.module+filepath, variable+extra_variable)

		11 - templatefile(path.cwd+filepath, variable)
		12 - templatefile(path.cwd+filepath, variable+extra_variable)

		error - wrong path
		13 - templatefile(wrong_filepath)
		14 - templatefile(path.root+wrong_filepath)
		15 - templatefile(path.module+wrong_filepath)
	*/

	const templateData = "template_data"

	withoutVariablesPath := filepath.Join(testDataDir, templateData, "withoutVariables.json")
	withoutVariablesOut, _ := ioutil.ReadFile(withoutVariablesPath)

	withVariablesPath := filepath.Join(testDataDir, templateData, "withVariables.json")
	withVariablesOutPath := filepath.Join(testDataDir, templateData, "withVariables_output.json")
	withVariablesOut, _ := ioutil.ReadFile(withVariablesOutPath)

	const (
		requiredVariable = "{ image_version \t\t = \n\n \"1.2.3.4\" }"
		extraVariables   = `{ image_version = "1.2.3.4", test = "test" }`
	)

	tests := []struct {
		name       string
		exprValue  string
		modfiledir string
		wantConfig string
		wantErr    error
	}{
		// without variables
		{
			name:       "valid file path | without variables",
			exprValue:  fmt.Sprintf(`templatefile("%s")`, withoutVariablesPath),
			modfiledir: "",
			wantConfig: string(withoutVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.root | without variables",
			exprValue:  `${templatefile("${path.root}/withoutVariables.json")}`,
			modfiledir: filepath.Dir(withoutVariablesPath),
			wantConfig: string(withoutVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.module | without variables",
			exprValue:  `templatefile("${path.module}/withoutVariables.json")`,
			modfiledir: filepath.Dir(withoutVariablesPath),
			wantConfig: string(withoutVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.cwd | with variables",
			exprValue:  `${templatefile("${path.cwd}/testdata/template_data/withoutVariables.json")}`,
			modfiledir: "",
			wantConfig: string(withoutVariablesOut),
			wantErr:    nil,
		},

		// with variables
		{
			name:       "valid file path | with required variables",
			exprValue:  fmt.Sprintf(`templatefile("%s", %s)`, withVariablesPath, requiredVariable),
			modfiledir: "",
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path | with extra variables",
			exprValue:  fmt.Sprintf(`${templatefile("%s", %s)}`, withVariablesPath, extraVariables),
			modfiledir: "",
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.root | with required variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.root}/withVariables.json", %s)`, requiredVariable),
			modfiledir: filepath.Dir(withVariablesPath),
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.root | with extra variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.root}/withVariables.json", %s)`, extraVariables),
			modfiledir: filepath.Dir(withVariablesPath),
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.module | with required variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.module}/withVariables.json", %s)`, requiredVariable),
			modfiledir: filepath.Dir(withVariablesPath),
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.module | with extra variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.module}/withVariables.json", %s)`, extraVariables),
			modfiledir: filepath.Dir(withVariablesPath),
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.cwd | with required variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.cwd}/testdata/template_data/withVariables.json", %s)`, requiredVariable),
			modfiledir: "",
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},
		{
			name:       "valid file path with path.cwd | with extra variables",
			exprValue:  fmt.Sprintf(`templatefile("${path.cwd}/testdata/template_data/withVariables.json", %s)`, extraVariables),
			modfiledir: "",
			wantConfig: string(withVariablesOut),
			wantErr:    nil,
		},

		// invalid file path
		{
			name:       "invalid file path | without variables",
			exprValue:  `templatefile("/foo/bar")`,
			modfiledir: "",
			wantConfig: "",
			wantErr:    fmt.Errorf("failed to read template file: open /foo/bar: no such file or directory"),
		},
		{
			name:       "invalid file path with path.root | without variables",
			exprValue:  `templatefile("${path.root}/foo/bar")`,
			modfiledir: "tim/blin",
			wantConfig: "",
			wantErr:    fmt.Errorf("failed to read template file: open tim/blin/foo/bar: no such file or directory"),
		},
		{
			name:       "invalid file path with path.module | without variables",
			exprValue:  `templatefile("${path.module}/foo/bar")`,
			modfiledir: "tim/blin",
			wantConfig: "",
			wantErr:    fmt.Errorf("failed to read template file: open tim/blin/foo/bar: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, gotErr := evalTemplatefileFunc(tt.exprValue, tt.modfiledir)
			if gotErr != nil && tt.wantErr == nil {
				t.Errorf("evalTemplatefileFunc() unexpected error | got = %v", gotErr)
			}

			if (gotErr != nil && tt.wantErr != nil) && (gotErr.Error() != tt.wantErr.Error()) {
				t.Errorf("evalTemplatefileFunc() error mismatch | got = %v, want = %v", gotErr, tt.wantErr)
			}

			if gotConfig != tt.wantConfig {
				t.Errorf("evalTemplatefileFunc() config mismatch | got = %v, want = %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
