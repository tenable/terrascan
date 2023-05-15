package writer

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
)

// test data
var (
	summaryWithNoViolations = results.ScanSummary{
		ResourcePath:     "test",
		IacType:          "terraform",
		Timestamp:        "2020-12-12 11:21:29.902796 +0000 UTC",
		TotalPolicies:    566,
		LowCount:         0,
		MediumCount:      0,
		HighCount:        1,
		ViolatedPolicies: 1,
	}

	outputWithPassedRules = policy.EngineOutput{
		ViolationStore: &results.ViolationStore{
			PassedRules: []*results.PassedRule{
				{
					RuleName:    "s3EnforceUserACL",
					Description: "S3 bucket Access is allowed to all AWS Account Users.",
					RuleID:      "AWS.S3Bucket.DS.High.1043",
					Severity:    "HIGH",
					Category:    "S3",
				},
			},
			Summary: summaryWithNoViolations,
		},
	}
	outputWithDirScanErrors = policy.EngineOutput{
		ViolationStore: &results.ViolationStore{
			DirScanErrors: []results.DirScanErr{
				{
					IacType:    "kustomize",
					Directory:  "test/e2e/test_data/iac/aws/aws_db_instance_violation",
					ErrMessage: "kustomization.y(a)ml file not found in the directory test/e2e/test_data/iac/aws/aws_db_instance_violation",
				},
				{
					IacType:    "helm",
					Directory:  "test/e2e/test_data/iac/aws/aws_db_instance_violation",
					ErrMessage: "no helm charts found in directory test/e2e/test_data/iac/aws/aws_db_instance_violation",
				},
			},
			PassedRules: []*results.PassedRule{
				{
					RuleName:    "s3EnforceUserACL",
					Description: "S3 bucket Access is allowed to all AWS Account Users.",
					RuleID:      "AWS.S3Bucket.DS.High.1043",
					Severity:    "HIGH",
					Category:    "S3",
				},
			},
			Summary: summaryWithNoViolations,
		},
	}
	vulnerabilitiesInputHumanReadable = policy.EngineOutput{
		ViolationStore: &results.ViolationStore{
			Vulnerabilities: []*results.Vulnerability{
				{
					Image:           "test",
					Container:       "test",
					VulnerabilityID: "CVE-2019-18276",
					PrimaryURL:      "http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-18276",
					Description:     "GNU Bash. Bash is the GNU Project's shell",
					Severity:        "HIGH",
				},
			},
			Summary: results.ScanSummary{
				ResourcePath:     "test",
				IacType:          "terraform",
				Timestamp:        "2020-12-12 11:21:29.902796 +0000 UTC",
				TotalPolicies:    566,
				LowCount:         0,
				MediumCount:      0,
				HighCount:        1,
				ViolatedPolicies: 1,
				Vulnerabilities:  &summaryWithNoViolations.ViolatedPolicies,
			},
		},
	}
	summaryWithRepoURLRepoRef = results.ScanSummary{
		ResourcePath:     "https://github.com/user/repository.git",
		Branch:           "main",
		IacType:          "terraform",
		Timestamp:        "2020-12-12 11:21:29.902796 +0000 UTC",
		TotalPolicies:    566,
		LowCount:         0,
		MediumCount:      0,
		HighCount:        1,
		ViolatedPolicies: 1,
	}
)

const (
	expectedOutput1 = `Violation Details -
    
	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	File           :	modules/m1/main.tf
	Line           :	20
	Severity       :	HIGH
	Rule Name      :	s3EnforceUserACL
	Rule ID        :	AWS.S3Bucket.DS.High.1043
	Resource Name  :	bucket
	Resource Type  :	aws_s3_bucket
	Category       :	S3
	
	-----------------------------------------------------------------------
	

Skipped Violations -
	
	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	File           :	modules/m1/main.tf
	Line           :	20
	Severity       :	HIGH
	Skip Comment   :	
	Rule Name      :	s3EnforceUserACL
	Rule ID        :	AWS.S3Bucket.DS.High.1043
	Resource Name  :	bucket
	Resource Type  :	aws_s3_bucket
	Category       :	S3
	
	-----------------------------------------------------------------------
	

Scan Summary -

	File/Folder         :	test
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1`

	expectedOutput2 = `Scan Summary -

	File/Folder         :	test
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1`

	expectedOutput3 = `Passed Rules -
    
	Rule ID        :	AWS.S3Bucket.DS.High.1043
	Rule Name      :	s3EnforceUserACL
	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	Severity       :	HIGH
	Category       :	S3
	
	-----------------------------------------------------------------------
	

Scan Summary -

	File/Folder         :	test
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1`

	expectedOutputWithDirScanError = `Scan Errors -

	IaC Type            :	kustomize
	Directory           :	test/e2e/test_data/iac/aws/aws_db_instance_violation
	Error Message       :	kustomization.y(a)ml file not found in the directory test/e2e/test_data/iac/aws/aws_db_instance_violation
	
	-----------------------------------------------------------------------
	
	IaC Type            :	helm
	Directory           :	test/e2e/test_data/iac/aws/aws_db_instance_violation
	Error Message       :	no helm charts found in directory test/e2e/test_data/iac/aws/aws_db_instance_violation
	
	-----------------------------------------------------------------------
	


Passed Rules -
    
	Rule ID        :	AWS.S3Bucket.DS.High.1043
	Rule Name      :	s3EnforceUserACL
	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	Severity       :	HIGH
	Category       :	S3
	
	-----------------------------------------------------------------------
	

Scan Summary -

	File/Folder         :	test
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1`

	vulnerabilityScanOutputHumanReadable = `Vulnerabilities Details -
    
	Description         :	GNU Bash. Bash is the GNU Project's shell
	Vulnerability ID    :	CVE-2019-18276
	Resource Name       :	""
	Resource Type       :	
	Image               :	test
	Package             :	
	Line                :	0
	Primary URL         :	http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-18276
	Primary URL         :	HIGH
	
	-----------------------------------------------------------------------
	

Scan Summary -

	File/Folder         :	test
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1
	Vulnerabilities     :	1`

	expectedOutput4 = `Scan Summary -

	File/Folder         :	https://github.com/user/repository.git
	Branch              :	main
	IaC Type            :	terraform
	Scanned At          :	2020-12-12 11:21:29.902796 +0000 UTC
	Policies Validated  :	566
	Violated Policies   :	1
	Low                 :	0
	Medium              :	0
	High                :	1`
)

func TestHumanReadableWriter(t *testing.T) {

	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Human Readable Writer: Violations",
			input:          violationsInput,
			expectedOutput: expectedOutput1,
		},
		{
			name: "Human Readable Writer: No Violations",
			input: policy.EngineOutput{
				ViolationStore: &results.ViolationStore{
					Summary: summaryWithNoViolations,
				},
			},
			expectedOutput: expectedOutput2,
		},
		{
			name:           "Human Readable Writer: With PassedRules",
			input:          outputWithPassedRules,
			expectedOutput: expectedOutput3,
		},
		{
			name:           "Human Readable Writer: With Vulnerabilities",
			input:          vulnerabilitiesInputHumanReadable,
			expectedOutput: vulnerabilityScanOutputHumanReadable,
		},
		{
			name: "Human Readable Writer: with repository url and branch",
			input: policy.EngineOutput{
				ViolationStore: &results.ViolationStore{
					Summary: summaryWithRepoURLRepoRef,
				},
			},
			expectedOutput: expectedOutput4,
		},
		{
			name:           "Human Readable Writer: with directory scan error",
			input:          outputWithDirScanErrors,
			expectedOutput: expectedOutputWithDirScanError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			if err := HumanReadableWriter(tt.input, w); (err != nil) != tt.expectedError {
				t.Errorf("HumanReadableWriter() error = gotErr: %v, wantErr: %v", err, tt.expectedError)
			}
			outputBytes := bf.Bytes()
			gotOutput := string(bytes.TrimSpace(outputBytes))
			if !strings.EqualFold(gotOutput, strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("HumanReadableWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
