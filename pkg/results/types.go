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

	"github.com/accurics/terrascan/pkg/utils"
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
	ResourceName string      `json:"resource_name" yaml:"resource_name" xml:"resource_name,attr"`
	ResourceType string      `json:"resource_type" yaml:"resource_type" xml:"resource_type,attr"`
	ResourceData interface{} `json:"-" yaml:"-" xml:"-"`
	File         string      `json:"file" yaml:"file" xml:"file,attr"`
	LineNumber   int         `json:"line" yaml:"line" xml:"line,attr"`
}

// ViolationStats Contains stats related to the violation data
type ViolationStats struct {
	LowCount    int `json:"low" yaml:"low" xml:"low,attr"`
	MediumCount int `json:"medium" yaml:"medium" xml:"medium,attr"`
	HighCount   int `json:"high" yaml:"high" xml:"high,attr"`
	TotalCount  int `json:"total" yaml:"total" xml:"total,attr"`
}

// ViolationStore Storage area for violation data
type ViolationStore struct {
	Violations []*Violation   `json:"violations" yaml:"violations" xml:"violations>violation"`
	Count      ViolationStats `json:"count" yaml:"count" xml:"count"`
}

// DefaultScanResult will hold the default scan summary data
type DefaultScanResult struct {
	IacType              string
	ResourcePath         string
	Timestamp            string
	TotalPolicies        int
	ShowViolationDetails bool
	ViolationStore
}

// Add adds two ViolationStores
func (vs ViolationStore) Add(extra ViolationStore) ViolationStore {
	// Just concatenate the slices, since order shouldn't be important
	vs.Violations = append(vs.Violations, extra.Violations...)

	// Add the counts
	vs.Count.LowCount += extra.Count.LowCount
	vs.Count.MediumCount += extra.Count.MediumCount
	vs.Count.HighCount += extra.Count.HighCount
	vs.Count.TotalCount += extra.Count.TotalCount

	return vs
}

// NewDefaultScanResult will initialize DefaultScanResult
func NewDefaultScanResult(iacType, iacFilePath, iacDirPath string, totalPolicyCount int, verbose bool, violationStore ViolationStore) *DefaultScanResult {
	sr := new(DefaultScanResult)

	if iacType == "" {
		// the default scan type is terraform
		sr.IacType = "terraform"
	} else {
		sr.IacType = iacType
	}

	if iacFilePath != "" {
		// can skip the error as the file validation is already done
		// while executor is initialized
		filePath, _ := utils.GetAbsPath(iacFilePath)
		sr.ResourcePath = filePath
	} else {
		// can skip the error as the directory validation is already done
		// while executor is initialized
		dirPath, _ := utils.GetAbsPath(iacDirPath)
		sr.ResourcePath = dirPath
	}
	sr.ShowViolationDetails = verbose
	sr.TotalPolicies = totalPolicyCount
	sr.ViolationStore = violationStore
	// set current time as scan time
	sr.Timestamp = time.Now().UTC().String()

	return sr
}
