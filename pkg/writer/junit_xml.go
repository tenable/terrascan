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
	"encoding/xml"
	"fmt"
	"io"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/version"
)

const (
	junitXMLFormat supportedFormat = "junit-xml"
	testSuiteName  string          = "TERRASCAN_POLICY_SUITE"
	testSuitesName string          = "TERRASCAN_POLICY_SUITES"
	testNameFormat string          = `[ERROR] resource: "%s" at line: %d, violates: RULE - %s`
)

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName  xml.Name `xml:"testsuites"`
	Tests    int      `xml:"tests,attr"`
	Name     string   `xml:"name,attr"`
	Failures int      `xml:"failures,attr"`
	Time     string   `xml:"time,attr"`
	Suites   []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Package    string          `xml:"package,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Severity  string   `xml:"severity,attr"`
	Category  string   `xml:"category,attr"`
	// omit empty time because today we do not have this data
	Time        string            `xml:"time,attr,omitempty"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func newJunitTestSuites(summary results.ScanSummary) JUnitTestSuites {
	return JUnitTestSuites{
		Tests:    summary.TotalPolicies,
		Name:     testSuitesName,
		Failures: summary.ViolatedPolicies,
		Time:     fmt.Sprint(summary.TotalTime),
	}
}

func newJunitTestSuite(summary results.ScanSummary) JUnitTestSuite {
	return JUnitTestSuite{
		Name:     testSuiteName,
		Tests:    summary.TotalPolicies,
		Time:     fmt.Sprint(summary.TotalTime),
		Failures: summary.ViolatedPolicies,
		Package:  summary.ResourcePath,
		Properties: []JUnitProperty{
			{
				Name:  "Terrascan Version",
				Value: version.Get(),
			},
		}}
}

func init() {
	RegisterWriter(junitXMLFormat, JUnitXMLWriter)
}

// JUnitXMLWriter writes scan summary in junit xml format
func JUnitXMLWriter(data interface{}, writer io.Writer) error {
	output, ok := data.(policy.EngineOutput)
	if !ok {
		return fmt.Errorf("incorrect input for JunitXML writer, supported type is policy.EngineOutput")
	}

	junitXMLOutput := convert(output)

	return XMLWriter(junitXMLOutput, writer)
}

// convert is helper func to convert engine output to JUnitTestSuites
func convert(output policy.EngineOutput) JUnitTestSuites {
	testSuites := newJunitTestSuites(output.Summary)
	// since we have a single suite for now, a suite will have same data as in root level element testsuites
	suite := newJunitTestSuite(output.Summary)

	tests := violationsToTestCases(output.ViolationStore.Violations, false)
	if tests != nil {
		suite.TestCases = append(suite.TestCases, tests...)
	}

	skippedTests := violationsToTestCases(output.ViolationStore.SkippedViolations, true)
	if skippedTests != nil {
		suite.TestCases = append(suite.TestCases, skippedTests...)
	}

	testSuites.Suites = append(testSuites.Suites, suite)

	return testSuites
}

// violationsToTestCases is helper func to convert scan violations to JunitTestCases
func violationsToTestCases(violations []*results.Violation, isSkipped bool) []JUnitTestCase {
	testCases := make([]JUnitTestCase, 0)
	for _, v := range violations {
		var testCase JUnitTestCase
		if isSkipped {
			testCase = JUnitTestCase{Failure: new(JUnitFailure), SkipMessage: new(JUnitSkipMessage)}
			testCase.SkipMessage.Message = v.Comment
		} else {
			testCase = JUnitTestCase{Failure: new(JUnitFailure)}
		}
		testCase.Classname = v.File
		testCase.Name = fmt.Sprintf(testNameFormat, v.ResourceName, v.LineNumber, v.RuleID)
		testCase.Severity = v.Severity
		testCase.Category = v.Category
		// since junitXML doesn't contain the attributes we want to show as violations
		// we would add details of violations in the failure message
		testCase.Failure.Message = getViolationString(*v)
		testCases = append(testCases, testCase)
	}
	return testCases
}

// getViolationString is used to get violation details as string
func getViolationString(v results.Violation) string {
	resourceName := v.ResourceName
	if resourceName == "" {
		resourceName = `""`
	}

	out := fmt.Sprintf("%s: %s, %s: %s, %s: %d, %s: %s, %s: %s, %s: %s, %s: %s, %s: %s, %s: %s",
		"Description", v.Description,
		"File", v.File,
		"Line", v.LineNumber,
		"Severity", v.Severity,
		"Rule Name", v.RuleName,
		"Rule ID", v.RuleID,
		"Resource Name", resourceName,
		"Resource Type", v.ResourceType,
		"Category", v.Category)
	return out
}
