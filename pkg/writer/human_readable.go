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
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/accurics/terrascan/pkg/results"
	"go.uber.org/zap"
)

const (
	humanReadbleFormat supportedFormat = "human"

	defaultTemplate string = `
{{if (gt (len .ViolationStore.Violations) 0) }}
Violation Details - 
	{{- $showDetails := .ViolationStore.Summary.ShowViolationDetails}}
    {{range $index, $element := .ViolationStore.Violations}}
	{{defaultViolations $element false | printf "%s"}}
	{{- if $showDetails -}}
	{{detailedViolations $element | printf "%s"}}
	{{- end}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}
{{- if (gt (len .ViolationStore.SkippedViolations) 0) }}
Skipped Violations - 
	{{- $showDetails := .ViolationStore.Summary.ShowViolationDetails}}
	{{range $index, $element := .ViolationStore.SkippedViolations}}
	{{defaultViolations $element true | printf "%s"}}
	{{- if $showDetails -}}
	{{detailedViolations $element | printf "%s"}}
	{{- end}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}
Scan Summary -

	{{scanSummary .ViolationStore.Summary | printf "%s"}}
`
)

func init() {
	RegisterWriter(humanReadbleFormat, HumanReadbleWriter)
}

// HumanReadbleWriter display scan summary in human readable format
func HumanReadbleWriter(data interface{}, writer io.Writer) error {
	tmpl, err := template.New("Report").Funcs(template.FuncMap{
		"defaultViolations":  defaultViolations,
		"detailedViolations": detailedViolations,
		"scanSummary":        scanSummary,
	}).Parse(defaultTemplate)
	if err != nil {
		zap.S().Errorf("failed to write human readable output. error: '%v'", err)
		return err
	}

	buffer := bytes.Buffer{}
	tmpl.Execute(&buffer, data)

	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	writer.Write([]byte{'\n'})
	return nil
}

func defaultViolations(v results.Violation, isSkipped bool) string {
	out := fmt.Sprintf("%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%d\n\t%-15v:\t%s\n\t",
		"Description", v.Description,
		"File", v.File,
		"Line", v.LineNumber,
		"Severity", v.Severity)
	if isSkipped {
		skipComment := fmt.Sprintf("%-15v:\t%s\n\t", "Skip Comment", v.Comment)
		out = out + skipComment
	}
	return out
}

func detailedViolations(v results.Violation) string {
	resourceName := v.ResourceName
	// print "" when as resource name value when it is empty
	if resourceName == "" {
		resourceName = `""`
	}
	out := fmt.Sprintf("%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t",
		"Rule Name", v.RuleName,
		"Rule ID", v.RuleID,
		"Resource Name", resourceName,
		"Resource Type", v.ResourceType,
		"Category", v.Category)
	return out
}

func scanSummary(s results.ScanSummary) string {
	out := fmt.Sprintf("%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t",
		"File/Folder", s.ResourcePath,
		"IaC Type", s.IacType,
		"Scanned At", s.Timestamp,
		"Policies Validated", s.TotalPolicies,
		"Violated Policies", s.ViolatedPolicies,
		"Low", s.LowCount,
		"Medium", s.MediumCount,
		"High", s.HighCount)
	return out
}
