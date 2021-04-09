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

package runtime

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
	tfv14 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v14"
	"github.com/accurics/terrascan/pkg/notifications/webhook"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

var (
	errMockLoadIacDir   = fmt.Errorf("mock LoadIacDir")
	errMockLoadIacFile  = fmt.Errorf("mock LoadIacFile")
	errMockPolicyEngine = fmt.Errorf("mock PolicyEngine")

	testDataDir     = "testdata"
	testDir         = filepath.Join(testDataDir, "testdir")
	testPoliciesDir = filepath.Join(testDataDir, "testpolicies")
)

// MockIacProvider mocks IacProvider interface
type MockIacProvider struct {
	output output.AllResourceConfigs
	err    error
}

func (m MockIacProvider) LoadIacDir(dir string) (output.AllResourceConfigs, error) {
	return m.output, m.err
}

func (m MockIacProvider) LoadIacFile(file string) (output.AllResourceConfigs, error) {
	return m.output, m.err
}

// mock policy engine
type MockPolicyEngine struct {
	err error
}

func (m MockPolicyEngine) Init(input string, scanRules, skipRules, categories []string, severity string) error {
	return m.err
}

func (m MockPolicyEngine) FilterRules(input string, scanRules, skipRules, categories []string, severity string) {
	/*
		This method does nothing. Required to fullfil the Engine interface contract
	*/
}

func (m MockPolicyEngine) Configure() error {
	return m.err
}

func (m MockPolicyEngine) Evaluate(input policy.EngineInput) (out policy.EngineOutput, err error) {
	return out, m.err
}

func (m MockPolicyEngine) GetResults() (out policy.EngineOutput) {
	return out
}

func (m MockPolicyEngine) Release() error {
	return m.err
}

func TestExecute(t *testing.T) {

	// TODO: add tests to validate output of Execute()
	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "test LoadIacDir error",
			executor: Executor{
				dirPath:     "./testdata/testdir",
				iacProvider: MockIacProvider{err: errMockLoadIacDir},
			},
			wantErr: errMockLoadIacDir,
		},
		{
			name: "test LoadIacDir no error",
			executor: Executor{
				dirPath:       "./testdata/testdir",
				iacProvider:   MockIacProvider{err: nil},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test LoadIacFile error",
			executor: Executor{
				filePath:    "./testdata/testfile",
				iacProvider: MockIacProvider{err: errMockLoadIacFile},
			},
			wantErr: errMockLoadIacFile,
		},
		{
			name: "test LoadIacFile no error",
			executor: Executor{
				filePath:      "./testdata/testfile",
				iacProvider:   MockIacProvider{err: nil},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test SendNofitications no error",
			executor: Executor{
				iacProvider:   MockIacProvider{err: nil},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test SendNofitications mock error",
			executor: Executor{
				iacProvider:   MockIacProvider{err: nil},
				notifiers:     []notifications.Notifier{&MockNotifier{err: errMockNotifier}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: errMockNotifier,
		},
		{
			name: "test policy enginer no error",
			executor: Executor{
				iacProvider:   MockIacProvider{err: nil},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test policy engine error",
			executor: Executor{
				iacProvider:   MockIacProvider{err: nil},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: errMockPolicyEngine}},
			},
			wantErr: errMockPolicyEngine,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.executor.Execute()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func TestInit(t *testing.T) {
	table := []struct {
		name            string
		executor        Executor
		wantErr         error
		wantIacProvider iacProvider.IacProvider
		wantNotifiers   []notifications.Notifier
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				policyPath: []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: filepath.Join(testDataDir, "webhook.toml"),
				policyPath: []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: filepath.Join(testDataDir, "invalid-notifier.toml"),
			},
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: filepath.Join(testDataDir, "does-not-exist"),
			},
			wantErr:         config.ErrNotPresent,
			wantIacProvider: &tfv14.TfV14{},
		},
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: filepath.Join(testDataDir, "webhook.toml"),
				policyPath: []string{filepath.Join(testDataDir, "notthere")},
			},
			wantErr:         fmt.Errorf("failed to initialize OPA policy engine"),
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config file with invalid category",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: filepath.Join(testDataDir, "invalid-category.toml"),
				policyPath: []string{filepath.Join(testDataDir, "notthere")},
			},
			wantErr:         fmt.Errorf("(3, 5): no value can start with c"),
			wantIacProvider: &tfv14.TfV14{},
		},
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				policyPath: []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: filepath.Join(testDataDir, "webhook.toml"),
				policyPath: []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: filepath.Join(testDataDir, "invalid-notifier.toml"),
			},
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: filepath.Join(testDataDir, "does-not-exist"),
			},
			wantErr:         config.ErrNotPresent,
			wantIacProvider: &tfv12.TfV12{},
		},
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:   filepath.Join(testDataDir, "testfile"),
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: filepath.Join(testDataDir, "webhook.toml"),
				policyPath: []string{filepath.Join(testDataDir, "notthere")},
			},
			wantErr:         fmt.Errorf("failed to initialize OPA policy engine"),
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
	}

	for _, tt := range table {

		t.Run(tt.name, func(t *testing.T) {
			configErr := config.LoadGlobalConfig(tt.executor.configFile)
			if configErr != nil {
				if !reflect.DeepEqual(configErr, tt.wantErr) {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", configErr, tt.wantErr)
				}
			} else {
				gotErr := tt.executor.Init()
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
				if !reflect.DeepEqual(tt.executor.iacProvider, tt.wantIacProvider) {
					t.Errorf("got: '%v', want: '%v'", tt.executor.iacProvider, tt.wantIacProvider)
				}
				if len(tt.wantNotifiers) > 0 {
					for i, notifier := range tt.executor.notifiers {
						if !reflect.DeepEqual(reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i])) {
							t.Errorf("got: '%v', want: '%v'", reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i]))
						}
					}
				}
			}
		})
	}
}

type flagSet struct {
	iacType    string
	iacVersion string
	filePath   string
	dirPath    string
	policyPath []string
	cloudType  []string
	categories []string
	severity   string
	scanRules  []string
	skipRules  []string
}

func TestNewExecutor(t *testing.T) {
	table := []struct {
		name           string
		wantErr        error
		configfile     string
		flags          flagSet
		wantScanRules  []string
		wantSkipRules  []string
		wantSeverity   string
		wantCategories []string
	}{
		{
			name:       "values passed through flag should override configfile value",
			configfile: filepath.Join(testDataDir, "scan-skip-rules-low-severity.toml"),
			wantErr:    nil,
			flags: flagSet{
				severity:   "high",
				scanRules:  []string{"AWS.S3Bucket.DS.High.1043"},
				skipRules:  []string{"accurics.kubernetes.IAM.109"},
				dirPath:    testDir,
				policyPath: []string{testPoliciesDir},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
			},
			wantSkipRules: []string{
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity:   "high",
			wantCategories: []string{"IDENTITY AND ACCESS MANAGEMENT", "RESILIENCE"},
		},
		{
			name:       "skipRules passed through flag should override configfile value",
			configfile: filepath.Join(testDataDir, "scan-skip-rules-low-severity.toml"),
			wantErr:    nil,
			flags: flagSet{
				skipRules:  []string{"accurics.kubernetes.IAM.109"},
				dirPath:    testDir,
				policyPath: []string{testPoliciesDir},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
				"accurics.kubernetes.IAM.107",
			},
			wantSkipRules: []string{
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity:   "low",
			wantCategories: []string{"IDENTITY AND ACCESS MANAGEMENT", "RESILIENCE"},
		},
		{
			name:       "scanRules passed through flag should override configfile value",
			configfile: filepath.Join(testDataDir, "scan-skip-rules-low-severity.toml"),
			wantErr:    nil,
			flags: flagSet{
				scanRules:  []string{"AWS.S3Bucket.DS.High.1043"},
				dirPath:    testDir,
				policyPath: []string{testPoliciesDir},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
			},
			wantSkipRules: []string{
				"AWS.S3Bucket.IAM.High.0370",
				"accurics.kubernetes.IAM.5",
				"accurics.kubernetes.OPS.461",
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity:   "low",
			wantCategories: []string{"IDENTITY AND ACCESS MANAGEMENT", "RESILIENCE"},
		},
		{
			name:       "severity and categories passed through flag should override configfile value",
			configfile: filepath.Join(testDataDir, "scan-skip-rules-low-severity.toml"),
			wantErr:    nil,
			flags: flagSet{
				severity:   "medium",
				dirPath:    testDataDir,
				policyPath: []string{testPoliciesDir},
				cloudType:  []string{"aws"},
				categories: []string{"DATA PROTECTION"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
				"accurics.kubernetes.IAM.107",
			},
			wantSkipRules: []string{
				"AWS.S3Bucket.IAM.High.0370",
				"accurics.kubernetes.IAM.5",
				"accurics.kubernetes.OPS.461",
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity:   "medium",
			wantCategories: []string{"DATA PROTECTION"},
		},
		{
			name:       "configfile value will be used if no flags are passed",
			configfile: filepath.Join(testDataDir, "scan-skip-rules-low-severity.toml"),
			wantErr:    nil,
			flags: flagSet{
				dirPath:    testDir,
				policyPath: []string{testPoliciesDir},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
				"accurics.kubernetes.IAM.107",
			},
			wantSkipRules: []string{
				"AWS.S3Bucket.IAM.High.0370",
				"accurics.kubernetes.IAM.5",
				"accurics.kubernetes.OPS.461",
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity:   "low",
			wantCategories: []string{"IDENTITY AND ACCESS MANAGEMENT", "RESILIENCE"},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			config.LoadGlobalConfig(tt.configfile)

			gotExecutor, gotErr := NewExecutor(tt.flags.iacType, tt.flags.iacVersion, tt.flags.cloudType, tt.flags.filePath, tt.flags.dirPath, tt.configfile, tt.flags.policyPath, tt.flags.scanRules, tt.flags.skipRules, tt.flags.categories, tt.flags.severity)

			if !reflect.DeepEqual(tt.wantErr, gotErr) {
				t.Errorf("Mismatch in error => got: '%v', want: '%v'", gotErr, tt.wantErr)
				t.Errorf("\n\n")
			}

			if !utils.IsSliceEqual(gotExecutor.scanRules, tt.wantScanRules) {
				t.Errorf("Mismatch in scanRules => got: '%v', want: '%v'", gotExecutor.scanRules, tt.wantScanRules)
				t.Errorf("\n\n")
			}

			if !utils.IsSliceEqual(gotExecutor.skipRules, tt.wantSkipRules) {
				t.Errorf("Mismatch in skipRules => got: '%v', want: '%v'", gotExecutor.skipRules, tt.wantSkipRules)
				t.Errorf("\n\n")
			}

			if gotExecutor.severity != tt.wantSeverity {
				t.Errorf("Mismatch in severity => got: '%v', want: '%v'", gotExecutor.severity, tt.wantSeverity)
				t.Errorf("\n\n")
			}

			if !utils.IsSliceEqual(gotExecutor.categories, tt.wantCategories) {
				t.Errorf("Mismatch in categories => got: '%v', want: '%v'", gotExecutor.categories, tt.wantCategories)
				t.Errorf("\n\n")
			}
		})
	}
}
