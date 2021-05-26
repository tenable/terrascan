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

package writer

import (
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/owenrumney/go-sarif/sarif"
	"io"
	"strings"
)

const (
	sarifFormat supportedFormat = "sarif"
)

func init() {
	RegisterWriter(sarifFormat, SarifWriter)
}

// SarifWriter writes sarif formatted violation results report
func SarifWriter(data interface{}, writer io.Writer) error {
	outputData := data.(policy.EngineOutput)
	//summary := outputData.Summary
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return err
	}

	run := sarif.NewRun("terrascan", "https://github.com/accurics/terrascan")
	// add a run to the report
	report.AddRun(run)

	for _, passedRule := range outputData.PassedRules {
		m := make(map[string]string)
		m["category"] = passedRule.Category
		m["severity"] = passedRule.Severity

		run.AddRule(string(passedRule.RuleID)).
			WithDescription(passedRule.Description).WithName(passedRule.RuleName).WithProperties(m)
	}

	// for each result add the rule, location and result to the report
	for _, violation := range outputData.Violations {
		m := make(map[string]string)
		m["category"] = violation.Category
		m["severity"] = violation.Severity

		rule := run.AddRule(string(violation.RuleID)).
			WithDescription(violation.Description).WithName(violation.RuleName).WithProperties(m)

		location := sarif.NewLocation().
			WithPhysicalLocation(sarif.NewPhysicalLocation().
				WithArtifactLocation(sarif.NewSimpleArtifactLocation(violation.File)).
				WithContextRegion(sarif.NewRegion().WithStartLine(violation.LineNumber)))

		if len(violation.ResourceType) > 0 && len(violation.ResourceName) > 0 {
			location.LogicalLocations = append(location.LogicalLocations, sarif.NewLogicalLocation().
				WithKind(violation.ResourceType).WithName(violation.ResourceName))
		}

		run.AddResult(rule.ID).
			WithMessage(sarif.NewTextMessage(violation.Description)).
			WithMessage(sarif.NewMarkdownMessage(violation.Description)).
			WithLevel(getSarifLevel(violation.Severity)).
			WithLocation(location)
	}

	// print the report to anything that implements `io.Writer`
	return report.PrettyWrite(writer)
}

func getSarifLevel(severity string) string {
	m := make(map[string]string, 3)
	m["low"] = "note"
	m["medium"] = "warning"
	m["high"] = "error"

	return m[strings.ToLower(severity)]
}
