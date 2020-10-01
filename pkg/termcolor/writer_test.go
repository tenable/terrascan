package termcolor

import (
	"regexp"
	"strings"
	"testing"

	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/writer"
)

var (
	jsonData, yamlData strings.Builder
)

func buildStore() *results.ViolationStore {
	res := results.NewViolationStore()

	res.AddResult(&results.Violation{
		RuleName:     "rule name",
		Description:  "description",
		RuleID:       "rule id",
		Severity:     "severity",
		Category:     "category",
		ResourceName: "resource name",
		ResourceType: "resource type",
		File:         "file",
		LineNumber:   1,
	})

	return res
}

func init() {
	res := buildStore()

	w := NewColorizedWriter(&jsonData)

	err := writer.Write("json", res, w)
	if err != nil {
		panic(err)
	}

	w = NewColorizedWriter(&yamlData)

	err = writer.Write("yaml", res, w)
	if err != nil {
		panic(err)
	}
}

func verifyLineWithStringIsColorized(s string, buf string, t *testing.T) {
	re := regexp.MustCompile(`(?m)^(.*` + s + `.*)$`)
	m := re.FindString(buf)
	if !strings.Contains(m, ColorPrefix) {
		t.Errorf("%s not colorized [%v]\n%s", s, m, buf)
	}
}

///////////  YAML

func TestYAMLBogusSeverityIsNotColorized(t *testing.T) {
	re := regexp.MustCompile(`(?m)^(.*severity.*)$`)
	m := re.FindString(yamlData.String())
	if strings.Contains(m, ColorPrefix) {
		t.Errorf("severity is colorized [%v]\n%s", m, yamlData.String())
	}
}

func TestYAMLRuleNameIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("rule_name", yamlData.String(), t)
}
func TestYAMLDescriptionIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("description", yamlData.String(), t)
}
func TestYAMLSeverityIsColorized(t *testing.T) {

	res := buildStore()
	yw := &strings.Builder{}
	w := NewColorizedWriter(yw)

	// HIGH, MEDIUM, LOW
	testSeverity := func(sev string) {
		res.Violations[0].Severity = sev
		yw.Reset()
		err := writer.Write("yaml", res, w)
		if err != nil {
			panic(err)
		}
		verifyLineWithStringIsColorized("severity", yw.String(), t)
	}
	testSeverity("HIGH")
	testSeverity("MEDIUM")
	testSeverity("LOW")
}
func TestYAMLResourceNameIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("resource_name", yamlData.String(), t)
}
func TestYAMLResourceTypeIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("resource_type", yamlData.String(), t)
}
func TestYAMLFileIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("file", yamlData.String(), t)
}

func TestYAMLCountIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("count", yamlData.String(), t)
}
func TestYAMLCountLowIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("low", yamlData.String(), t)
}
func TestYAMLCountMediumIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("medium", yamlData.String(), t)
}
func TestYAMLCountHighIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("high", yamlData.String(), t)
}

///////////  JSON

func TestJSONBogusSeverityIsNotColorized(t *testing.T) {
	re := regexp.MustCompile(`(?m)^(.*severity.*)$`)
	m := re.FindString(jsonData.String())
	if strings.Contains(m, ColorPrefix) {
		t.Errorf("severity is colorized [%v]\n%s", m, jsonData.String())
	}
}

func TestJSONRuleNameIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("rule_name", jsonData.String(), t)
}
func TestJSONDescriptionIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("description", jsonData.String(), t)
}
func TestJSONSeverityIsColorized(t *testing.T) {

	res := buildStore()
	yw := &strings.Builder{}
	w := NewColorizedWriter(yw)

	// HIGH, MEDIUM, LOW
	testSeverity := func(sev string) {
		res.Violations[0].Severity = sev
		yw.Reset()
		err := writer.Write("json", res, w)
		if err != nil {
			panic(err)
		}
		verifyLineWithStringIsColorized("severity", yw.String(), t)
	}
	testSeverity("HIGH")
	testSeverity("MEDIUM")
	testSeverity("LOW")
}
func TestJSONResourceNameIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("resource_name", jsonData.String(), t)
}
func TestJSONResourceTypeIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("resource_type", jsonData.String(), t)
}
func TestJSONFileIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("file", jsonData.String(), t)
}

func TestJSONCountIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("count", jsonData.String(), t)
}
func TestJSONCountLowIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("low", jsonData.String(), t)
}
func TestJSONCountMediumIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("medium", jsonData.String(), t)
}
func TestJSONCountHighIsColorized(t *testing.T) {
	verifyLineWithStringIsColorized("high", jsonData.String(), t)
}
