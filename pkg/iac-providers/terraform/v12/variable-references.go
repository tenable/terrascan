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
	varRefPattern = regexp.MustCompile(`(\$\{)?var\.(?P<name>\w*)(\})?`)
)

// isVarRef returns true if the given string is a variable reference
func isVarRef(attrVal string) bool {
	return varRefPattern.MatchString(attrVal)
}

// getVarName returns the actual variable name as configured in IaC. It trims
// of "${var." prefix and "}" suffix and returns the variable name
func getVarName(varRef string) (string, string) {

	// 1. extract the exact variable reference from the string
	varExpr := varRefPattern.FindString(varRef)

	// 2. extract variable name from variable reference
	match := varRefPattern.FindStringSubmatch(varRef)
	result := make(map[string]string)
	for i, name := range varRefPattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	varName := result["name"]

	zap.S().Debugf("extracted variable name %q from reference %q", varName, varRef)
	return varName, varExpr
}

// ResolveVarRef returns the variable value as configured in IaC config in module
func (r *RefResolver) ResolveVarRef(varRef string) interface{} {

	// get variable name from varRef
	varName, varExpr := getVarName(varRef)

	// check if variable name exists in the map of variables read from IaC
	hclVar, present := r.Config.Module.Variables[varName]
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
				resolvedVal := strings.Replace(varRef, varExpr, valStr, 1)
				if varRef == resolvedVal {
					zap.S().Debugf("resolved str variable ref refers to self: '%v'", varRef)
					return varRef
				}
				zap.S().Debugf("resolved str variable ref: '%v', value: '%v'", varRef, resolvedVal)
				return r.ResolveStrRef(resolvedVal)
			}
			zap.S().Debugf("resolved variable ref: '%v', value: '%v'", varRef, val)
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
	if r.ParentModuleCall == nil {
		return varRef
	}

	// get variable name from varRef
	varName, varExpr := getVarName(varRef)

	// get initialized variables from module call
	ParentModuleCallBody, ok := r.ParentModuleCall.Config.(*hclsyntax.Body)
	if !ok {
		return varRef
	}

	// get varName from module call, if present
	varAttr, present := ParentModuleCallBody.Attributes[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in parent module call", varName, varRef)
		return varRef
	}

	// read source file
	fileBytes, err := ioutil.ReadFile(r.ParentModuleCall.SourceAddrRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", r.ParentModuleCall.SourceAddr, err)
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
		resolvedVal := strings.Replace(varRef, varExpr, valStr, 1)
		if varRef == resolvedVal {
			zap.S().Debugf("resolved str variable ref refers to self: '%v'", varRef)
			return varRef
		}
		zap.S().Debugf("resolved str variable ref: '%v', value: '%v'", varRef, string(resolvedVal))
		return r.ResolveStrRef(resolvedVal)
	}

	// return extracted value
	zap.S().Debugf("resolved variable ref: '%v', value: '%v'", varRef, val)
	return val
}
