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
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"go.uber.org/zap"
)

var (
	// reference patterns
	localRefPattern = regexp.MustCompile(`\$\{local\..*\}`)
)

// isLocalRef returns true if the given string is a local value reference
func isLocalRef(attrVal string) bool {
	return localRefPattern.MatchString(attrVal)
}

// getLocalName returns the actual local value name as configured in IaC. It
// trims of "${local." prefix and "}" suffix and returns the local value name
func getLocalName(localRef string) string {

	var (
		localPrefix = "${local."
		localSuffix = "}"
	)

	// 1. split at "${local.", remove everything before
	split := strings.Split(localRef, localPrefix)
	localName := split[1]

	// 2. split at "}", remove everything after
	split = strings.Split(localName, localSuffix)
	localName = split[0]

	zap.S().Debugf("extracted local name %q from reference %q", localName, localRef)
	return localName
}

// ResolveLocalRef returns the local value as configured in IaC config in module
func (r *RefResolver) ResolveLocalRef(localRef string) interface{} {

	// get local name from localRef
	localName := getLocalName(localRef)

	// check if local name exists in the map of locals read from IaC
	localAttr, present := r.Config.Module.Locals[localName]
	if !present {
		zap.S().Debugf("local name: %q, ref: %q not present in locals", localName, localRef)
		return localRef
	}

	// read source file
	fileBytes, err := ioutil.ReadFile(localAttr.DeclRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", localAttr.DeclRange.Filename, err)
		return localRef
	}

	// extract values from attribute expressions as golang interface{}
	c := converter{bytes: fileBytes}
	val, err := c.convertExpression(localAttr.Expr.(hclsyntax.Expression))
	if err != nil {
		zap.S().Errorf("failed to convert expression '%v', ref: '%v'", localAttr.Expr, localRef)
		return localRef
	}

	// replace the local value reference string with actual value
	if reflect.TypeOf(val).Kind() == reflect.String {
		valStr := val.(string)
		resolvedVal := localRefPattern.ReplaceAll([]byte(localRef), []byte(valStr))
		zap.S().Debugf("resolved str local value ref: '%v', value: '%v'", localRef, string(resolvedVal))
		return r.ResolveStrRef(string(resolvedVal))
	}

	// return extracted value
	zap.S().Debugf("resolved local value ref: '%v', value: '%v'", localRef, val)
	return val
}
