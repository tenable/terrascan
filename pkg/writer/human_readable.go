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

package writer

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/tenable/terrascan/pkg/results"
	"go.uber.org/zap"
)

const (
	oneLineCommonFormat = "%-20v:\t%s\n\t"
)

const (
	humanReadableFormat supportedFormat = "human"

	defaultTemplate string = `
{{if (gt (len .ViolationStore.DirScanErrors) 0)}}
Scan Errors -
{{range $index, $element := .ViolationStore.DirScanErrors}}
	{{dirScanErrors $element | printf "%s"}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}
{{if (gt (len .ViolationStore.PassedRules) 0) }}
Passed Rules -
    {{range $index, $element := .ViolationStore.PassedRules}}
	{{passedRules $element | printf "%s"}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}
{{- if (gt (len .ViolationStore.Violations) 0) }}
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
{{- if (gt (len .ViolationStore.Vulnerabilities) 0) }}
Vulnerabilities Details -
    {{range $index, $element := .ViolationStore.Vulnerabilities}}
	{{defaultVulnerabilities $element | printf "%s"}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}
Scan Summary -

	{{scanSummary .ViolationStore.Summary | printf "%s"}}
`
)

func init() {
	RegisterWriter(humanReadableFormat, HumanReadableWriter)
}

// HumanReadableWriter display scan summary in human readable format
func HumanReadableWriter(data interface{}, writers []io.Writer) error {
	tmpl, err := template.New("Report").Funcs(template.FuncMap{
		"defaultViolations":      defaultViolations,
		"detailedViolations":     detailedViolations,
		"scanSummary":            scanSummary,
		"passedRules":            passedRules,
		"defaultVulnerabilities": defaultVulnerabilities,
		"dirScanErrors":          dirScanErrors,
	}).Parse(defaultTemplate)
	if err != nil {
		zap.S().Errorf("failed to write human readable output. error: '%v'", err)
		return err
	}

	buffer := bytes.Buffer{}
	tmpl.Execute(&buffer, data)

	for _, writer := range writers {
		_, err = writer.Write(buffer.Bytes())
		if err != nil {
			return err
		}
		writer.Write([]byte{'\n'})
	}

	return nil
}

func defaultViolations(v results.Violation, isSkipped bool) string {
	part := fmt.Sprintf("%-15v:\t%s\n\t%-15v:\t%s\n\t",
		"Description", v.Description,
		"File", v.File)

	if v.ModuleName != "" {
		moduleName := fmt.Sprintf("%-15v:\t%s\n\t", "Module Name", v.ModuleName)
		part = part + moduleName
	}

	if v.PlanRoot != "" {
		planRoot := fmt.Sprintf("%-15v:\t%s\n\t", "Plan Root", v.PlanRoot)
		part = part + planRoot
	}
	out := fmt.Sprintf("%-15v:\t%d\n\t%-15v:\t%s\n\t",
		"Line", v.LineNumber,
		"Severity", v.Severity)
	if isSkipped {
		skipComment := fmt.Sprintf("%-15v:\t%s\n\t", "Skip Comment", v.Comment)
		out = out + skipComment
	}
	out = part + out

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

	out := fmt.Sprintf(oneLineCommonFormat,
		"File/Folder", s.ResourcePath)

	if s.Branch != "" {
		out += fmt.Sprintf(oneLineCommonFormat, "Branch", s.Branch)
	}

	out += fmt.Sprintf("%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t%-20v:\t%d\n\t",
		"IaC Type", s.IacType,
		"Scanned At", s.Timestamp,
		"Policies Validated", s.TotalPolicies,
		"Violated Policies", s.ViolatedPolicies,
		"Low", s.LowCount,
		"Medium", s.MediumCount,
		"High", s.HighCount,
	)

	if s.Vulnerabilities != nil {
		out += fmt.Sprintf("%-20v:\t%d\n\t", "Vulnerabilities", *s.Vulnerabilities)
	}

	return out
}

func passedRules(v results.PassedRule) string {
	out := fmt.Sprintf("%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t%-15v:\t%s\n\t",
		"Rule ID", v.RuleID,
		"Rule Name", v.RuleName,
		"Description", v.Description,
		"Severity", v.Severity,
		"Category", v.Category)
	return out
}

func defaultVulnerabilities(v results.Vulnerability) string {
	resourceName := v.ResourceName
	// print "" when as resource name value when it is empty
	if resourceName == "" {
		resourceName = `""`
	}
	out := fmt.Sprintf("%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%d\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t",
		"Description", v.Description,
		"Vulnerability ID", v.VulnerabilityID,
		"Resource Name", resourceName,
		"Resource Type", v.ResourceType,
		"Image", v.Image,
		"Package", v.Package,
		"Line", v.LineNumber,
		"Primary URL", v.PrimaryURL,
		"Primary URL", v.Severity)
	return out
}

func dirScanErrors(d results.DirScanErr) string {
	out := fmt.Sprintf("%-20v:\t%s\n\t%-20v:\t%s\n\t%-20v:\t%s\n\t",
		"IaC Type", d.IacType,
		"Directory", d.Directory,
		"Error Message", d.ErrMessage)
	return out
}
