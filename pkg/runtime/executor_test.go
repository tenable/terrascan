package runtime

import (
	"fmt"
	"reflect"
	"testing"

	cloudProvider "github.com/accurics/terrascan/pkg/cloud-providers"
	awsProvider "github.com/accurics/terrascan/pkg/cloud-providers/aws"
	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
)

var (
	errMockLoadIacDir           = fmt.Errorf("mock LoadIacDir")
	errMockLoadIacFile          = fmt.Errorf("mock LoadIacFile")
	errMockCreateNormalizedJSON = fmt.Errorf("mock CreateNormalizedJSON")
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

// MockCloudProvider mocks CloudProvider interface
type MockCloudProvider struct {
	err error
}

func (m MockCloudProvider) CreateNormalizedJSON(data output.AllResourceConfigs) (mockInterface interface{}, err error) {
	return data, m.err
}

func TestExecute(t *testing.T) {

	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "test LoadIacDir",
			executor: Executor{
				dirPath:     "./testdata/testdir",
				iacProvider: MockIacProvider{err: errMockLoadIacDir},
			},
			wantErr: errMockLoadIacDir,
		},
		{
			name: "test LoadIacFile",
			executor: Executor{
				filePath:    "./testdata/testfile",
				iacProvider: MockIacProvider{err: errMockLoadIacFile},
			},
			wantErr: errMockLoadIacFile,
		},
		{
			name: "test CreateNormalizedJSON error",
			executor: Executor{
				filePath:      "./testdata/testfile",
				iacProvider:   MockIacProvider{err: nil},
				cloudProvider: MockCloudProvider{err: errMockCreateNormalizedJSON},
			},
			wantErr: errMockCreateNormalizedJSON,
		},
		{
			name: "test CreateNormalizedJSON",
			executor: Executor{
				filePath:      "./testdata/testfile",
				iacProvider:   MockIacProvider{err: nil},
				cloudProvider: MockCloudProvider{err: nil},
			},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.Execute()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func TestInit(t *testing.T) {

	table := []struct {
		name              string
		executor          Executor
		wantErr           error
		wantIacProvider   iacProvider.IacProvider
		wantCloudProvider cloudProvider.CloudProvider
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  "aws",
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr:           nil,
			wantIacProvider:   &tfv12.TfV12{},
			wantCloudProvider: &awsProvider.AWSProvider{},
		},
	}

	for _, tt := range table {
		gotErr := tt.executor.Init()
		if !reflect.DeepEqual(gotErr, tt.wantErr) {
			t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
		}
		if !reflect.DeepEqual(tt.executor.iacProvider, tt.wantIacProvider) {
			t.Errorf("got: '%v', want: '%v'", tt.executor.iacProvider, tt.wantIacProvider)
		}
		if !reflect.DeepEqual(tt.executor.cloudProvider, tt.wantCloudProvider) {
			t.Errorf("got: '%v', want: '%v'", tt.executor.cloudProvider, tt.wantCloudProvider)
		}
	}
}
