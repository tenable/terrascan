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

package results

import (
	"time"
)

// Violation Contains data for each violation
type Violation struct {
	RuleName     string      `json:"rule_name" yaml:"rule_name" xml:"rule_name,attr"`
	Description  string      `json:"description" yaml:"description" xml:"description,attr"`
	RuleID       string      `json:"rule_id" yaml:"rule_id" xml:"rule_id,attr"`
	Severity     string      `json:"severity" yaml:"severity" xml:"severity,attr"`
	Category     string      `json:"category" yaml:"category" xml:"category,attr"`
	RuleFile     string      `json:"-" yaml:"-" xml:"-"`
	RuleData     interface{} `json:"-" yaml:"-" xml:"-"`
	Comment      string      `json:"skip_comment,omitempty" yaml:"skip_comment,omitempty" xml:"skip_comment,omitempty"`
	ResourceName string      `json:"resource_name" yaml:"resource_name" xml:"resource_name,attr"`
	ResourceType string      `json:"resource_type" yaml:"resource_type" xml:"resource_type,attr"`
	ResourceData interface{} `json:"-" yaml:"-" xml:"-"`
	File         string      `json:"file" yaml:"file" xml:"file,attr"`
	LineNumber   int         `json:"line" yaml:"line" xml:"line,attr"`
}

// ViolationStore Storage area for violation data
type ViolationStore struct {
	Violations        []*Violation `json:"violations" yaml:"violations" xml:"violations>violation"`
	SkippedViolations []*Violation `json:"skipped_violations" yaml:"skipped_violations" xml:"skipped_violations>violation"`
	Summary           ScanSummary  `json:"scan_summary" yaml:"scan_summary" xml:"scan_summary"`
}

// ScanSummary will hold the default scan summary data
type ScanSummary struct {
	ResourcePath         string `json:"file/folder" yaml:"file/folder" xml:"file/folder,attr"`
	IacType              string `json:"iac_type" yaml:"iac_type" xml:"iac_type,attr"`
	Timestamp            string `json:"scanned_at" yaml:"scanned_at" xml:"scanned_at,attr"`
	ShowViolationDetails bool   `json:"-" yaml:"-" xml:"-"`
	TotalPolicies        int    `json:"policies_validated" yaml:"policies_validated" xml:"policies_validated,attr"`
	ViolatedPolicies     int    `json:"violated_policies" yaml:"violated_policies" xml:"violated_policies,attr"`
	LowCount             int    `json:"low" yaml:"low" xml:"low,attr"`
	MediumCount          int    `json:"medium" yaml:"medium" xml:"medium,attr"`
	HighCount            int    `json:"high" yaml:"high" xml:"high,attr"`
	TotalTime            int64  `json:"-" yaml:"-" xml:"-"`
}

// Add adds two ViolationStores
func (vs ViolationStore) Add(extra ViolationStore) ViolationStore {
	// Just concatenate the slices, since order shouldn't be important
	vs.Violations = append(vs.Violations, extra.Violations...)
	vs.SkippedViolations = append(vs.SkippedViolations, extra.SkippedViolations...)

	// Add the scan summary
	vs.Summary.LowCount += extra.Summary.LowCount
	vs.Summary.MediumCount += extra.Summary.MediumCount
	vs.Summary.HighCount += extra.Summary.HighCount
	vs.Summary.ViolatedPolicies += extra.Summary.ViolatedPolicies
	vs.Summary.TotalPolicies += extra.Summary.TotalPolicies

	return vs
}

// AddSummary will update the summary with remaining details
func (vs *ViolationStore) AddSummary(iacType, iacResourcePath string) {

	vs.Summary.IacType = iacType
	vs.Summary.ResourcePath = iacResourcePath
	vs.Summary.Timestamp = time.Now().UTC().String()
}
