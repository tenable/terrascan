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

package tfv12

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"go.uber.org/zap"
)

var (
	skipRulesPattern = regexp.MustCompile(`#ts:skip=\s*(([A-Za-z0-9]+\.?){5})(\s*,\s*([A-Za-z0-9]+\.?){5})*`)
	skipRulesPrefix  = "#ts:skip="
)

// CreateResourceConfig creates output.ResourceConfig
func CreateResourceConfig(managedResource *hclConfigs.Resource) (resourceConfig output.ResourceConfig, err error) {

	// read source file
	fileBytes, err := ioutil.ReadFile(managedResource.DeclRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to read terraform file")
	}

	// convert resource config from hcl.Body to map[string]interface{}
	c := converter{bytes: fileBytes}
	hclBody := managedResource.Config.(*hclsyntax.Body)
	goOut, err := c.convertBody(hclBody)
	if err != nil {
		zap.S().Errorf("failed to convert hcl.Body to go struct; resource '%s', file: '%s'. error: '%v'",
			managedResource.Name, managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to convert hcl.Body to go struct")
	}

	// create a resource config
	resourceConfig = output.ResourceConfig{
		ID:        fmt.Sprintf("%s.%s", managedResource.Type, managedResource.Name),
		Name:      managedResource.Name,
		Type:      managedResource.Type,
		Source:    managedResource.DeclRange.Filename,
		Line:      managedResource.DeclRange.Start.Line,
		Config:    goOut,
		SkipRules: getSkipRules(c.rangeSource(hclBody.Range())),
	}

	// successful
	zap.S().Debugf("created resource config for resource '%s', file: '%s'", resourceConfig.Name, resourceConfig.Source)
	return resourceConfig, nil
}

// getSkipRules returns a list of rules to be skipped. The rules to be skipped
// can be set in terraform resource config with the following comma separated pattern:
// #ts:skip=AWS.S3Bucket.DS.High.1043, AWS.S3Bucket.DS.High.1044
func getSkipRules(body string) []string {

	var skipRules []string

	// check if any rules comments are present in body
	if !skipRulesPattern.MatchString(body) {
		return skipRules
	}

	// get all skip rule comments
	comments := skipRulesPattern.FindAllString(body, -1)

	// extract rule ids from comments
	for _, c := range comments {
		c = strings.TrimPrefix(c, skipRulesPrefix)
		c = strings.ReplaceAll(c, ",", " ")
		skipRules = append(skipRules, strings.Fields(c)...)
	}
	return skipRules
}
