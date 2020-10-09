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
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

var (
	// reference patterns
	varRefPattern = regexp.MustCompile(`\$\{var\.\w*\}`)
)

// isVarRef returns true if the given string is a variable reference
func isVarRef(attrVal string) bool {
	return varRefPattern.MatchString(attrVal)
}

// getVarName returns the actual variable name as configured in IaC. It trims
// of "${var." prefix and "}" suffix and returns the variable name
func getVarName(varRef string) string {

	var (
		varPrefix = "${var."
		varSuffix = "}"
	)

	// 1. split at "${var.", remove everything before
	split := strings.Split(varRef, varPrefix)
	varName := split[1]

	// 2. split at "}", remove everything after
	split = strings.Split(varName, varSuffix)
	varName = split[0]

	zap.S().Debugf("extracted variable name %q from reference %q", varName, varRef)
	return varName
}

// ResolveVarRef returns the variable value as configured in IaC config in module
func (r *RefResolver) ResolveVarRef(varRef string) interface{} {

	// get variable name from varRef
	varName := getVarName(varRef)

	// check if variable name exists in the map of variables read from IaC
	hclVar, present := r.variables[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in variables", varName, varRef)
		return varRef
	}

	// return varRef if default value is not present, or value is a NilVal,
	// or if default value is not known
	if hclVar.Default.IsNull() || hclVar.Default.RawEquals(cty.NilVal) || !hclVar.Default.IsKnown() {
		return varRef
	}

	// default value is of cty.Value type, convert it to native golang type
	// based on cty.Type, determine golang type
	for _, converter := range ctyConverterFuncs {
		if val, err := converter(hclVar.Default); err == nil {
			zap.S().Debugf("resolved variable ref '%v', value: '%v'", varRef, val)

			// replace the variable reference string with actual value
			if reflect.TypeOf(val).Kind() == reflect.String {
				valStr := val.(string)
				resolvedVal := varRefPattern.ReplaceAll([]byte(varRef), []byte(valStr))
				return string(resolvedVal)
			}
			return val
		}
	}
	zap.S().Debugf("failed to convert cty.Value '%v' to golang native value", hclVar.Default.GoString())

	// return original reference if cty conversion attempts fail
	return varRef
}

// ResolveVarRefFromParentModuleCall returns the variable value as configured in
// ModuleCall from parent module. The resolved value can be an absolute value
// (string, int, bool etc.) or it can also be another reference, which may
// need further resolution
func (r *RefResolver) ResolveVarRefFromParentModuleCall(varRef string) interface{} {

	zap.S().Debugf("resolving variable ref %q in parent module call", varRef)

	// if module call struct is nil, nothing to process
	if r.parentModuleCall == nil {
		return varRef
	}

	// get variable name from varRef
	varName := getVarName(varRef)

	// get initialized variables from module call
	parentModuleCallBody, ok := r.parentModuleCall.Config.(*hclsyntax.Body)
	if !ok {
		return varRef
	}

	// get varName from module call, if present
	varAttr, present := parentModuleCallBody.Attributes[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in parent module call", varName, varRef)
		return varRef
	}

	// read source file
	fileBytes, err := ioutil.ReadFile(r.parentModuleCall.SourceAddrRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", r.parentModuleCall.SourceAddr, err)
		return varRef
	}

	// extract values from attribute expressions as golang interface{}
	c := converter{bytes: fileBytes}
	val, err := c.convertExpression(varAttr.Expr)
	if err != nil {
		zap.S().Errorf("failed to convert expression '%v', ref: '%v'", varAttr.Expr, varRef)
		return varRef
	}

	// replace the variable reference string with actual value
	if reflect.TypeOf(val).Kind() == reflect.String {
		valStr := val.(string)
		resolvedVal := varRefPattern.ReplaceAll([]byte(varRef), []byte(valStr))
		zap.S().Debugf("resolved str variable ref: '%v', value: '%v'", varRef, string(resolvedVal))
		return string(resolvedVal)
	}

	// return extracted value
	zap.S().Debugf("resolved variable ref: '%v', value: '%v'", varRef, val)
	return val
}
