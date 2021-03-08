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
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/config"
	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
	tfv14 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v14"
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/notifications/webhook"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

var (
	errMockLoadIacDir   = fmt.Errorf("mock LoadIacDir")
	errMockLoadIacFile  = fmt.Errorf("mock LoadIacFile")
	errMockPolicyEngine = fmt.Errorf("mock PolicyEngine")
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
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				policyPath: []string{"./testdata/testpolicies"},
			},
			wantErr:         nil,
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: "./testdata/webhook.toml",
				policyPath: []string{"./testdata/testpolicies"},
			},
			wantErr:         nil,
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: "testdata/invalid-notifier.toml",
			},
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: "./testdata/does-not-exist",
			},
			wantErr:         config.ErrNotPresent,
			wantIacProvider: &tfv14.TfV14{},
		},
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v14",
				configFile: "./testdata/webhook.toml",
				policyPath: []string{"./testdata/notthere"},
			},
			wantErr:         fmt.Errorf("failed to initialize OPA policy engine"),
			wantIacProvider: &tfv14.TfV14{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.Init()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.executor.iacProvider, tt.wantIacProvider) {
				t.Errorf("got: '%v', want: '%v'", tt.executor.iacProvider, tt.wantIacProvider)
			}
			for i, notifier := range tt.executor.notifiers {
				if !reflect.DeepEqual(reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i])) {
					t.Errorf("got: '%v', want: '%v'", reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i]))
				}
			}
		})
	}

	table = []struct {
		name            string
		executor        Executor
		wantErr         error
		wantIacProvider iacProvider.IacProvider
		wantNotifiers   []notifications.Notifier
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				policyPath: []string{"./testdata/testpolicies"},
			},
			wantErr:         nil,
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: "./testdata/webhook.toml",
				policyPath: []string{"./testdata/testpolicies"},
			},
			wantErr:         nil,
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: "testdata/invalid-notifier.toml",
			},
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: "./testdata/does-not-exist",
			},
			wantErr:         config.ErrNotPresent,
			wantIacProvider: &tfv12.TfV12{},
		},
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
				configFile: "./testdata/webhook.toml",
				policyPath: []string{"./testdata/notthere"},
			},
			wantErr:         fmt.Errorf("failed to initialize OPA policy engine"),
			wantIacProvider: &tfv12.TfV12{},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.Init()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.executor.iacProvider, tt.wantIacProvider) {
				t.Errorf("got: '%v', want: '%v'", tt.executor.iacProvider, tt.wantIacProvider)
			}
			for i, notifier := range tt.executor.notifiers {
				if !reflect.DeepEqual(reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i])) {
					t.Errorf("got: '%v', want: '%v'", reflect.TypeOf(notifier), reflect.TypeOf(tt.wantNotifiers[i]))
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
		name          string
		wantErr       error
		configfile    string
		flags         flagSet
		wantScanRules []string
		wantSkipRules []string
		wantSeverity  string
	}{
		{
			name:       "values passed through flag should override configfile value",
			configfile: "./testdata/scan-skip-rules-low-severity.toml",
			wantErr:    nil,
			flags: flagSet{
				severity:   "high",
				scanRules:  []string{"AWS.S3Bucket.DS.High.1043"},
				skipRules:  []string{"accurics.kubernetes.IAM.109"},
				dirPath:    "./testdata/testdir",
				policyPath: []string{"./testdata/testpolicies"},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
			},
			wantSkipRules: []string{
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity: "high",
		},
		{
			name:       "skipRules passed through flag should override configfile value",
			configfile: "./testdata/scan-skip-rules-low-severity.toml",
			wantErr:    nil,
			flags: flagSet{
				skipRules:  []string{"accurics.kubernetes.IAM.109"},
				dirPath:    "./testdata/testdir",
				policyPath: []string{"./testdata/testpolicies"},
				cloudType:  []string{"aws"},
			},
			wantScanRules: []string{
				"AWS.S3Bucket.DS.High.1043",
				"accurics.kubernetes.IAM.107",
			},
			wantSkipRules: []string{
				"accurics.kubernetes.IAM.109",
			},
			wantSeverity: "low",
		},
		{
			name:       "scanRules passed through flag should override configfile value",
			configfile: "./testdata/scan-skip-rules-low-severity.toml",
			wantErr:    nil,
			flags: flagSet{
				scanRules:  []string{"AWS.S3Bucket.DS.High.1043"},
				dirPath:    "./testdata/testdir",
				policyPath: []string{"./testdata/testpolicies"},
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
			wantSeverity: "low",
		},
		{
			name:       "severity passed through flag should override configfile value",
			configfile: "./testdata/scan-skip-rules-low-severity.toml",
			wantErr:    nil,
			flags: flagSet{
				severity:   "medium",
				dirPath:    "./testdata/testdir",
				policyPath: []string{"./testdata/testpolicies"},
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
			wantSeverity: "medium",
		},
		{
			name:       "configfile value will be used if no flags are passed",
			configfile: "./testdata/scan-skip-rules-low-severity.toml",
			wantErr:    nil,
			flags: flagSet{
				dirPath:    "./testdata/testdir",
				policyPath: []string{"./testdata/testpolicies"},
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
			wantSeverity: "low",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotExecutor, gotErr := NewExecutor(tt.flags.iacType, tt.flags.iacVersion, tt.flags.cloudType, tt.flags.filePath, tt.flags.dirPath, tt.configfile, tt.flags.policyPath, tt.flags.scanRules, tt.flags.skipRules, tt.flags.categories, tt.flags.severity)

			if !reflect.DeepEqual(tt.wantErr, gotErr) {
				t.Errorf("Mismatch in error => got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
			if utils.IsSliceEqual(gotExecutor.scanRules, tt.wantScanRules) && utils.IsSliceEqual(gotExecutor.skipRules, tt.wantSkipRules) && gotExecutor.severity != tt.wantSeverity {
				t.Errorf("got: 'scanRules = %v, skipRules = %v, severity = %s', want: 'scanRules = %v, skipRules = %v, severity = %s'", gotExecutor.scanRules, gotExecutor.skipRules, gotExecutor.severity, tt.wantScanRules, tt.wantSkipRules, tt.wantSeverity)
			}
		})
	}
}
