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
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/pkg/version"
	"github.com/go-errors/errors"
	"github.com/owenrumney/go-sarif/sarif"
	"go.uber.org/zap"
	"io"
	"path/filepath"
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
	return writeSarif(data, writer, false)
}

func writeSarif(data interface{}, writer io.Writer, forGithub bool) error {
	outputData := data.(policy.EngineOutput)
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return err
	}

	run := sarif.NewRun("terrascan", "https://github.com/accurics/terrascan")
	run.Tool.Driver.WithVersion(version.GetNumeric())
	// add a run to the report
	report.AddRun(run)

	for _, passedRule := range outputData.PassedRules {
		m := sarif.NewPropertyBag()
		m.Properties["category"] = passedRule.Category
		m.Properties["severity"] = passedRule.Severity

		run.AddRule(passedRule.RuleID).
			WithDescription(passedRule.Description).WithName(passedRule.RuleName).WithProperties(m.Properties)
	}

	// for each result add the rule, location and result to the report
	for _, violation := range outputData.Violations {
		m := sarif.NewPropertyBag()
		m.Properties["category"] = violation.Category
		m.Properties["severity"] = violation.Severity

		rule := run.AddRule(violation.RuleID).
			WithDescription(violation.Description).WithName(violation.RuleName).WithProperties(m.Properties)

		var artifactLocation *sarif.ArtifactLocation

		if forGithub {
			artifactLocation = sarif.NewSimpleArtifactLocation(violation.File).
				WithUriBaseId(outputData.Summary.ResourcePath)
		} else {
			absFilePath, err := getAbsoluteFilePath(outputData.Summary.ResourcePath, violation.File)
			if err != nil {
				return errors.Errorf("unable to create absolute path, error: %v", err)
			}
			uriFilePath, err := utils.GetFileURI(absFilePath)
			if err != nil {
				return errors.Errorf("unable to create uri path, error: %v", err)
			}
			artifactLocation = sarif.NewSimpleArtifactLocation(uriFilePath)
		}

		location := sarif.NewLocation().WithPhysicalLocation(sarif.NewPhysicalLocation().
			WithArtifactLocation(artifactLocation).WithRegion(sarif.NewRegion().WithStartLine(violation.LineNumber)))

		if len(violation.ResourceType) > 0 && len(violation.ResourceName) > 0 {
			location.LogicalLocations = append(location.LogicalLocations, sarif.NewLogicalLocation().
				WithKind(violation.ResourceType).WithName(violation.ResourceName))
		}

		run.AddResult(rule.ID).
			WithMessage(sarif.NewTextMessage(violation.Description)).
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

func getAbsoluteFilePath(resourcePath, filePath string) (string, error) {
	if !filepath.IsAbs(resourcePath) {
		resourcePath, err := filepath.Abs(resourcePath)
		if err != nil {
			zap.S().Errorf("unable to get absolute path for %s, error: %v", resourcePath, err)
			return "", err
		}
	}
	fileMode := utils.GetFileMode(resourcePath)
	if fileMode != nil && (*fileMode).IsDir() {
		return filepath.Join(resourcePath, filePath), nil
	}
	return resourcePath, nil
}
