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

package runtime

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/go-multierror"
	tfv15 "github.com/tenable/terrascan/pkg/iac-providers/terraform/v15"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/vulnerability"

	iacProvider "github.com/tenable/terrascan/pkg/iac-providers"
	armv1 "github.com/tenable/terrascan/pkg/iac-providers/arm/v1"
	cftv1 "github.com/tenable/terrascan/pkg/iac-providers/cft/v1"
	dockerv1 "github.com/tenable/terrascan/pkg/iac-providers/docker/v1"
	helmv3 "github.com/tenable/terrascan/pkg/iac-providers/helm/v3"
	k8sv1 "github.com/tenable/terrascan/pkg/iac-providers/kubernetes/v1"
	kustomizev4 "github.com/tenable/terrascan/pkg/iac-providers/kustomize/v4"
	tfv12 "github.com/tenable/terrascan/pkg/iac-providers/terraform/v12"
	tfv14 "github.com/tenable/terrascan/pkg/iac-providers/terraform/v14"
	"github.com/tenable/terrascan/pkg/notifications/webhook"

	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/notifications"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/utils"
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

func (m MockIacProvider) LoadIacDir(dir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return m.output, m.err
}

func (m MockIacProvider) LoadIacFile(file string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return m.output, m.err
}

func (m MockIacProvider) Name() string {
	return "mock-iac"
}

// mock policy engine
type MockPolicyEngine struct {
	err error
}
type MockVulnerabiltyEngine struct {
	out vulnerability.EngineOutput
}

func (m MockPolicyEngine) Init(input string, filter policy.PreLoadFilter) error {
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

func (m MockPolicyEngine) Evaluate(input policy.EngineInput, filter policy.PreScanFilter) (out policy.EngineOutput, err error) {
	return out, m.err
}

func (m MockVulnerabiltyEngine) ReportVulnerability(input vulnerability.EngineInput, options map[string]interface{}) (out vulnerability.EngineOutput) {
	return m.out
}

func (m MockVulnerabiltyEngine) FetchVulnerabilities(input output.AllResourceConfigs, options map[string]interface{}) (out output.AllResourceConfigs) {
	return out
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
		name            string
		configOnly      bool
		configWithError bool
		executor        Executor
		wantErr         error
	}{
		{
			name: "test LoadIacDir error",
			executor: Executor{
				dirPath:      testDir,
				iacProviders: []iacProvider.IacProvider{MockIacProvider{err: errMockLoadIacDir}},
			},
			wantErr: multierror.Append(errMockLoadIacDir),
		},
		{
			name: "test LoadIacDir no error",
			executor: Executor{
				dirPath:       testDir,
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test LoadIacFile error",
			executor: Executor{
				filePath:     filepath.Join(testDataDir, "testfile"),
				iacProviders: []iacProvider.IacProvider{MockIacProvider{err: errMockLoadIacFile}},
			},
			// iac file load doesn't return go-multierror
			wantErr: errMockLoadIacFile,
		},
		{
			name: "test LoadIacFile no error",
			executor: Executor{
				filePath:      filepath.Join(testDataDir, "testfile"),
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test SendNofitications no error",
			executor: Executor{
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test policy enginer no error",
			executor: Executor{
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "test policy engine error",
			executor: Executor{
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: errMockPolicyEngine}},
			},
			wantErr: errMockPolicyEngine,
		},
		{
			name: "test find vulnerability engine",
			executor: Executor{
				iacProviders:  []iacProvider.IacProvider{MockIacProvider{err: nil}},
				notifiers:     []notifications.Notifier{&MockNotifier{err: nil}},
				policyEngines: []policy.Engine{MockPolicyEngine{err: nil}},
				vulnerabilityEngine: MockVulnerabiltyEngine{
					out: vulnerability.EngineOutput{
						XMLName:        xml.Name{},
						ViolationStore: results.NewViolationStore(),
					},
				},
				findVulnerabilities: true,
			},
			wantErr: nil,
		},
		{
			name: "has scan errors with all the iac providers",
			executor: Executor{
				dirPath:      testDir,
				iacType:      "all",
				iacProviders: []iacProvider.IacProvider{MockIacProvider{err: errMockLoadIacDir}},
			},
			wantErr: nil,
		},
		{
			name:            "test config with error",
			configWithError: true,
			executor: Executor{
				dirPath:      testDir,
				iacType:      "terraform",
				iacProviders: []iacProvider.IacProvider{MockIacProvider{err: errMockLoadIacDir}},
			},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.executor.Execute(tt.configOnly, tt.configWithError)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func TestInitPolicyEngine(t *testing.T) {
	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
				policyPath:  []string{filepath.Join(testDataDir, "notthere")},
			},
			wantErr: fmt.Errorf("failed to initialize OPA policy engine"),
		},
		{
			name: "invalid policy path",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v12",
				policyPath:  []string{filepath.Join(testDataDir, "notthere")},
			},
			wantErr: fmt.Errorf("failed to initialize OPA policy engine"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.initPolicyEngines()
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
		configFile      string
		wantIacProvider []iacProvider.IacProvider
		wantNotifiers   []notifications.Notifier
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
				policyPath:  []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv14.TfV14{}},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "empty iac type with -d flag",
			executor: Executor{
				dirPath:     testDataDir,
				policyTypes: []string{"aws"},
				policyPath:  []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&armv1.ARMV1{}, &cftv1.CFTV1{}, &dockerv1.DockerV1{}, &helmv3.HelmV3{}, &k8sv1.K8sV1{}, &kustomizev4.KustomizeV4{}, &tfv15.TfV15{}},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "empty iac type with -f flag",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				policyTypes: []string{"aws"},
				policyPath:  []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv15.TfV15{}},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
				policyPath:  []string{testPoliciesDir},
			},
			configFile:      filepath.Join(testDataDir, "webhook.toml"),
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv14.TfV14{}},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
			},
			configFile:      filepath.Join(testDataDir, "invalid-notifier.toml"),
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: []iacProvider.IacProvider{&tfv14.TfV14{}},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
			},
			configFile:      filepath.Join(testDataDir, "does-not-exist"),
			wantErr:         config.ErrNotPresent,
			wantIacProvider: []iacProvider.IacProvider{&tfv14.TfV14{}},
		},
		{
			name: "config file with invalid category",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v14",
				policyPath:  []string{filepath.Join(testDataDir, "notthere")},
			},
			configFile:      filepath.Join(testDataDir, "invalid-category.toml"),
			wantErr:         fmt.Errorf("(3, 5): no value can start with c"),
			wantIacProvider: []iacProvider.IacProvider{&tfv14.TfV14{}},
		},
		{
			name: "valid filePath",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v12",
				policyPath:  []string{testPoliciesDir},
			},
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv12.TfV12{}},
			wantNotifiers:   []notifications.Notifier{},
		},
		{
			name: "valid notifier",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v12",
				policyPath:  []string{testPoliciesDir},
			},
			configFile:      filepath.Join(testDataDir, "webhook.toml"),
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv12.TfV12{}},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "invalid notifier",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v12",
			},
			configFile:      filepath.Join(testDataDir, "invalid-notifier.toml"),
			wantErr:         fmt.Errorf("notifier not supported"),
			wantIacProvider: []iacProvider.IacProvider{&tfv12.TfV12{}},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
		{
			name: "config not present",
			executor: Executor{
				filePath:    filepath.Join(testDataDir, "testfile"),
				dirPath:     "",
				policyTypes: []string{"aws"},
				iacType:     "terraform",
				iacVersion:  "v12",
			},
			configFile:      filepath.Join(testDataDir, "does-not-exist"),
			wantErr:         config.ErrNotPresent,
			wantIacProvider: []iacProvider.IacProvider{&tfv12.TfV12{}},
		},
		{
			name: "notification webhook configs passed as CLI args",
			executor: Executor{
				filePath:                 filepath.Join(testDataDir, "testfile"),
				dirPath:                  "",
				policyTypes:              []string{"aws"},
				iacType:                  "terraform",
				iacVersion:               "v12",
				notificationWebhookURL:   "http://some-host.url",
				notificationWebhookToken: "token",
			},
			configFile:      filepath.Join(testDataDir, "webhook.toml"),
			wantErr:         nil,
			wantIacProvider: []iacProvider.IacProvider{&tfv12.TfV12{}},
			wantNotifiers:   []notifications.Notifier{&webhook.Webhook{}},
		},
	}

	for _, tt := range table {

		t.Run(tt.name, func(t *testing.T) {
			configErr := config.LoadGlobalConfig(tt.configFile)
			if configErr != nil {
				if !reflect.DeepEqual(configErr, tt.wantErr) {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", configErr, tt.wantErr)
				}
			} else {
				gotErr := tt.executor.Init()
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
				if !reflect.DeepEqual(tt.executor.iacProviders, tt.wantIacProvider) {
					t.Errorf("got: '%v', want: '%v'", tt.executor.iacProviders, tt.wantIacProvider)
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
	iacType                  string
	iacVersion               string
	filePath                 string
	dirPath                  string
	policyPath               []string
	policyTypes              []string
	categories               []string
	severity                 string
	scanRules                []string
	skipRules                []string
	notificationWebhookURL   string
	notificationWebhookToken string
	repoURL                  string
	repoRef                  string
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
				severity:    "high",
				scanRules:   []string{"AWS.S3Bucket.DS.High.1043"},
				skipRules:   []string{"accurics.kubernetes.IAM.109"},
				dirPath:     testDir,
				policyPath:  []string{testPoliciesDir},
				policyTypes: []string{"aws"},
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
				skipRules:   []string{"accurics.kubernetes.IAM.109"},
				dirPath:     testDir,
				policyPath:  []string{testPoliciesDir},
				policyTypes: []string{"aws"},
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
				scanRules:   []string{"AWS.S3Bucket.DS.High.1043"},
				dirPath:     testDir,
				policyPath:  []string{testPoliciesDir},
				policyTypes: []string{"aws"},
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
				severity:    "medium",
				dirPath:     testDataDir,
				policyPath:  []string{testPoliciesDir},
				policyTypes: []string{"aws"},
				categories:  []string{"DATA PROTECTION"},
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
				dirPath:     testDir,
				policyPath:  []string{testPoliciesDir},
				policyTypes: []string{"aws"},
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

			gotExecutor, gotErr := NewExecutor(tt.flags.iacType, tt.flags.iacVersion, tt.flags.policyTypes, tt.flags.filePath, tt.flags.dirPath, tt.flags.policyPath, tt.flags.scanRules, tt.flags.skipRules, tt.flags.categories, tt.flags.severity, false, false, false, tt.flags.notificationWebhookURL, tt.flags.notificationWebhookToken, tt.flags.repoURL, tt.flags.repoRef, []string{})

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
