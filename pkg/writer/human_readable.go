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
	"io"
	"text/template"

	"go.uber.org/zap"
)

const (
	humanReadbleFormat supportedFormat = "human"

	defaultTemplate string = `
{{if (gt (len .ViolationStore.Violations) 0) }}
Violation Details - 
	{{- $showDetails := .ViolationStore.Summary.ShowViolationDetails}}
    {{range $index, $element := .ViolationStore.Violations}}
	{{printf "%-15v" "Description"}}:{{"\t"}}{{$element.Description}}
	{{printf "%-15v" "File"}}:{{"\t"}}{{$element.File}}
	{{printf "%-15v" "Line"}}:{{"\t"}}{{$element.LineNumber}}
	{{printf "%-15v" "Severity"}}:{{"\t"}}{{$element.Severity}}
	{{if $showDetails -}}
	{{printf "%-15v" "Rule Name"}}:{{"\t"}}{{$element.RuleName}}
	{{printf "%-15v" "Rule ID"}}:{{"\t"}}{{$element.RuleID}}
	{{printf "%-15v" "Resource Name"}}:{{"\t"}}{{if $element.ResourceName}}{{$element.ResourceName}}{{else}}""{{end}}
	{{printf "%-15v" "Resource Type"}}:{{"\t"}}{{$element.ResourceType}}
	{{printf "%-15v" "Category"}}:{{"\t"}}{{$element.Category}}
	{{end}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}	
Scan Summary -

	{{printf "%-20v" "File/Folder"}}:{{"\t"}}{{.ViolationStore.Summary.ResourcePath}}
	{{printf "%-20v" "IaC Type"}}:{{"\t"}}{{.ViolationStore.Summary.IacType}}
	{{printf "%-20v" "Scanned At"}}:{{"\t"}}{{.ViolationStore.Summary.Timestamp}}
	{{printf "%-20v" "Policies Validated"}}:{{"\t"}}{{.ViolationStore.Summary.TotalPolicies}}
	{{printf "%-20v" "Violated Policies"}}:{{"\t"}}{{.ViolationStore.Summary.ViolatedPolicies}}
	{{printf "%-20v" "Low"}}:{{"\t"}}{{.ViolationStore.Summary.LowCount}}
	{{printf "%-20v" "Medium"}}:{{"\t"}}{{.ViolationStore.Summary.MediumCount}}
	{{printf "%-20v" "High"}}:{{"\t"}}{{.ViolationStore.Summary.HighCount}}
`
)

func init() {
	RegisterWriter(humanReadbleFormat, HumanReadbleWriter)
}

// HumanReadbleWriter display scan summary in human readable format
func HumanReadbleWriter(data interface{}, writer io.Writer) error {
	tmpl, err := template.New("Report").Parse(defaultTemplate)
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
