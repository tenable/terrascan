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
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/notifications/webhook"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/pelletier/go-toml"
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

func (m MockPolicyEngine) Init(input string, scanRules, skipRules []string) error {
	return m.err
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
				dirPath:      "./testdata/testdir",
				iacProvider:  MockIacProvider{err: nil},
				policyEngine: []policy.Engine{MockPolicyEngine{err: nil}},
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
				filePath:     "./testdata/testfile",
				iacProvider:  MockIacProvider{err: nil},
				policyEngine: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test SendNofitications no error",
			executor: Executor{
				iacProvider:  MockIacProvider{err: nil},
				notifiers:    []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngine: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test SendNofitications mock error",
			executor: Executor{
				iacProvider:  MockIacProvider{err: nil},
				notifiers:    []notifications.Notifier{&MockNotifier{err: errMockNotifier}},
				policyEngine: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: errMockNotifier,
		},
		{
			name: "test policy enginer no error",
			executor: Executor{
				iacProvider:  MockIacProvider{err: nil},
				notifiers:    []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngine: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test policy engine error",
			executor: Executor{
				iacProvider:  MockIacProvider{err: nil},
				notifiers:    []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngine: []policy.Engine{MockPolicyEngine{err: errMockPolicyEngine}},
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
			wantIacProvider: nil,
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

func TestGetRulesInTomlTree(t *testing.T) {
	// test data
	fileLoadErr := "error while loading toml file"
	emptyTomlTree, err := config.LoadConfig("testdata/empty.toml")
	if err != nil {
		t.Fatalf(fileLoadErr, err)
	}

	configFileData, err := config.LoadConfig("testdata/scan-skip-rules.toml")
	if err != nil {
		t.Fatalf(fileLoadErr, err)
	}
	validRulesFormat := (configFileData.Get(rulesKey)).(*toml.Tree)

	configFileData, err = config.LoadConfig("testdata/invalid-scan-skip-rules.toml")
	if err != nil {
		t.Fatalf(fileLoadErr, err)
	}
	invalidRulesFormat := (configFileData.Get(rulesKey)).(*toml.Tree)

	type args struct {
		tree       *toml.Tree
		configFile string
		key        string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "empty toml file",
			args: args{
				tree:       emptyTomlTree,
				configFile: "",
				key:        scanRulesKey,
			},
			want: []string{},
		},
		{
			name: "get scan rules - valid data",
			args: args{
				tree:       validRulesFormat,
				configFile: "",
				key:        scanRulesKey,
			},
			want: []string{"AWS.S3Bucket.DS.High.1043", "accurics.kubernetes.IAM.107"},
		},
		{
			name: "get skip rules - valid data",
			args: args{
				tree:       validRulesFormat,
				configFile: "",
				key:        skipRulesKey,
			},
			want: []string{"AWS.S3Bucket.IAM.High.0370", "accurics.kubernetes.IAM.5",
				"accurics.kubernetes.OPS.461", "accurics.kubernetes.IAM.109"},
		},
		{
			name: "get scan rules - invalid data",
			args: args{
				tree:       invalidRulesFormat,
				configFile: "",
				key:        scanRulesKey,
			},
			wantErr: true,
		},
		{
			name: "get skip rules - invalid data",
			args: args{
				tree:       invalidRulesFormat,
				configFile: "",
				key:        skipRulesKey,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRulesInTomlTree(tt.args.tree, tt.args.configFile, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRulesInTomlTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRulesInTomlTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutorInitScanAndSkipRules(t *testing.T) {
	assertionTestName := "valid config file with scan and skip rules"

	type fields struct {
		configFile string
		scanRules  []string
		skipRules  []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:   "no config file",
			fields: fields{},
		},
		{
			name: "config file doesn't exist",
			fields: fields{
				configFile: "testdata/test.toml",
			},
			wantErr: true,
		},
		{
			name: "empty config file",
			fields: fields{
				configFile: "testdata/empty.toml",
			},
		},
		{
			name: "config file with empty rules",
			fields: fields{
				configFile: "testdata/webhook.toml",
			},
		},
		{
			name: assertionTestName,
			fields: fields{
				configFile: "testdata/scan-skip-rules.toml",
				scanRules:  []string{"testRuleA", "testRuleB"},
				skipRules:  []string{"testRuleC"},
			},
		},
		{
			name: "valid config file with invalid scan rules",
			fields: fields{
				configFile: "testdata/invalid-scan-skip-rules.toml",
			},
			wantErr: true,
		},
		{
			name: "valid config file with invalid skip rules",
			fields: fields{
				configFile: "testdata/invalid-skip-rules.toml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Executor{
				configFile: tt.fields.configFile,
				scanRules:  tt.fields.scanRules,
				skipRules:  tt.fields.skipRules,
			}
			if err := e.initScanAndSkipRules(); (err != nil) != tt.wantErr {
				t.Errorf("Executor.initScanAndSkipRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == assertionTestName {
				if len(e.scanRules) != 4 && len(e.skipRules) != 5 {
					t.Errorf("Expected scanRules: %d and skipRules: %d, got scanRules: %d and skipRules: %d", 4, 5, len(e.scanRules), len(e.skipRules))
				}
			}
		})
	}
}
