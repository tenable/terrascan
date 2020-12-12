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
	{{- $showDetails := .ShowViolationDetails}}
    {{range $index, $element := .ViolationStore.Violations}}
	Description{{"\t"}}:{{"\t"}}{{$element.Description}}
	File{{"\t\t"}}:{{"\t"}}{{$element.File}}
	Line{{"\t\t"}}:{{"\t"}}{{$element.LineNumber}}
	Severity{{"\t"}}:{{"\t"}}{{$element.Severity}}
	{{if $showDetails -}}
	Rule Name{{"\t"}}:{{"\t"}}{{$element.RuleName}}
	Rule ID{{"\t\t"}}:{{"\t"}}{{$element.RuleID}}
	Resource Name{{"\t"}}:{{"\t"}}{{$element.ResourceName}}
	Resource Type{{"\t"}}:{{"\t"}}{{$element.ResourceType}}
	Category{{"\t"}}:{{"\t"}}{{$element.Category}}
	{{end}}
	-----------------------------------------------------------------------
	{{end}}
{{end}}	
Scan Summary -

	File/Folder{{"\t\t"}}:{{"\t"}}{{.ResourcePath}}
	IaC Type{{"\t\t"}}:{{"\t"}}{{.IacType}}
	Scanned At{{"\t\t"}}:{{"\t"}}{{.Timestamp}}
	Policies Validated{{"\t"}}:{{"\t"}}{{.TotalPolicies}}
	Violated Policies{{"\t"}}:{{"\t"}}{{.ViolationStore.Count.TotalCount}}
	Low{{"\t\t\t"}}:{{"\t"}}{{.ViolationStore.Count.LowCount}}
	Medium{{"\t\t\t"}}:{{"\t"}}{{.ViolationStore.Count.MediumCount}}
	High{{"\t\t\t"}}:{{"\t"}}{{.ViolationStore.Count.HighCount}}
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
