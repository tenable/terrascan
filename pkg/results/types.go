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

// Violation Contains data for each violation
type Violation struct {
	RuleName     string      `json:"ruleName" yaml:"ruleName" xml:"ruleName,attr"`
	Description  string      `json:"description" yaml:"description" xml:"description,attr"`
	RuleID       string      `json:"rule" yaml:"rule" xml:"rule,attr"`
	Severity     string      `json:"severity" yaml:"severity" xml:"severity,attr"`
	Category     string      `json:"category" yaml:"category" xml:"category,attr"`
	RuleFile     string      `json:"ruleFile" yaml:"ruleFile" xml:"ruleFile,attr"`
	RuleData     interface{} `json:"-" yaml:"-" xml:"-"`
	ResourceName string      `json:"resourceName" yaml:"resourceName" xml:"resourceName,attr"`
	ResourceType string      `json:"resourceType" yaml:"resourceType" xml:"resourceType,attr"`
	ResourceData interface{} `json:"resourceData" yaml:"resourceData" xml:"resourceData,attr"`
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
	Violations []*Violation   `json:"violations" yaml:"violations" xml:"violations,attr"`
	Count      ViolationStats `json:"count" yaml:"count" xml:"count,attr"`
}
