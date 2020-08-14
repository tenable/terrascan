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
	"io"
	"text/template"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/policy"
)

const (
	consoleFormat                 supportedFormat = "console"
	defaultViolationFormat                        = "{{ .Severity }} | {{ .RuleID }} | {{ .Category }} | {{ .RuleName }} | {{ .Description }}"
	defaultViolationSummaryFormat                 = "Violation Count | [HIGH] {{ .HighCount }} | [MEDIUM] {{ .MediumCount }} | [LOW] {{ .LowCount }} | Total: {{ .TotalCount }}"
)

func init() {
	RegisterWriter(consoleFormat, ConsoleWriter)
}

// ConsoleWriter prints data in a custom format
func ConsoleWriter(data policy.EngineOutput, writer io.Writer) error {
	violationFormat := defaultViolationFormat
	violationSummaryFormat := defaultViolationSummaryFormat

	violationTemplate, err := template.New("violation").Parse(violationFormat)
	if err != nil {
		zap.S().Error("Failed to create violation template",
			zap.String("violation format", violationFormat), zap.Error(err))
	}

	violationSummaryTemplate, err := template.New("violation-summary").Parse(defaultViolationSummaryFormat)
	if err != nil {
		zap.S().Error("Failed to create violation summary template",
			zap.String("violation summary format", violationSummaryFormat), zap.Error(err))
	}

	for i := range data.ViolationStore.Violations {
		err = violationTemplate.Execute(writer, data.ViolationStore.Violations[i])
		if err != nil {
			zap.S().Error("Failed to print violation",
				zap.String("violation format", violationFormat), zap.Error(err))
		}
	}

	err = violationSummaryTemplate.Execute(writer, data.ViolationStore.Stats)
	if err != nil {
		zap.S().Error("Failed to print violation summary",
			zap.String("violation summary format", violationSummaryFormat), zap.Error(err))
	}

	writer.Write([]byte{'\n'})
	return nil
}
