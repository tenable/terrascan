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
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

var (
	// variable reference pattern
	refPattern    = regexp.MustCompile(`\$\{.*\}`)
	varRefPattern = regexp.MustCompile(`\$\{var\..*\}`)
	varPrefix     = "${var."
	varSuffix     = "}"
)

// RefResolver tries to resolve all the references in the given terraform
// config
type RefResolver struct {
	variables  map[string]*hclConfigs.Variable
	moduleCall *hclConfigs.ModuleCall
}

// NewRefResolver returns a new RefResolver struct
func NewRefResolver(variables map[string]*hclConfigs.Variable, moduleCall *hclConfigs.ModuleCall) *RefResolver {
	return &RefResolver{
		variables:  variables,
		moduleCall: moduleCall,
	}
}

// ResolveRefs figures out all the variable references in the resource config
// and tries to replace the variable references, if possible,
// with actual value
func (r *RefResolver) ResolveRefs(config jsonObj) jsonObj {

	// iterate over every attribute in the config
	for k, v := range config {

		var (
			valKind = reflect.TypeOf(v).Kind()
			valType = reflect.TypeOf(v).String()
		)

		if valKind == reflect.String {
			// case 1: config value is a string; in resource config, refs
			// are of the type string
			if isVarRef(v.(string)) {

				// variable values initialized in the parent module call take
				// precedence over values initialized in the same module,
				// hence the following steps:

				// 1. resolve variables initialized in parent module call
				var resolved bool
				config[k], resolved = r.getVarValueFromParentModuleCall(v.(string))
				if resolved {
					continue
				}

				// 2. resolve variables initialized in same module
				config[k], _ = r.getVarValue(v.(string))
			}
		} else if valType == "tfv12.jsonObj" && valKind == reflect.Map {
			// case 2: config value is of type jsonObj
			config[k] = r.ResolveRefs(v.(jsonObj))
		} else if valType == "[]tfv12.jsonObj" && valKind == reflect.Slice {
			// case 3: config value is of type []jsonObj

			// type assert interface{} -> []jsonObj
			sConfig, ok := v.([]jsonObj)
			if !ok {
				continue
			}

			// iterate over the []jsonObj to resolve refs
			for i, c := range sConfig {
				sConfig[i] = r.ResolveRefs(c)
			}
			config[k] = sConfig
		}
	}

	return config
}

// getVarValueFromParentModuleCall returns the variable value as configured in
// ModuleCall from parent module
func (r *RefResolver) getVarValueFromParentModuleCall(varRef string) (interface{}, bool) {

	// if module call struct is nil, nothing to process
	if r.moduleCall == nil {
		return varRef, false
	}

	// get variable name from varRef
	varName := getVarName(varRef)

	// get initialized variables from module call
	moduleCallBody, ok := r.moduleCall.Config.(*hclsyntax.Body)
	if !ok {
		return varRef, false
	}

	// get varName from module call, if present
	varAttr, present := moduleCallBody.Attributes[varName]
	if !present {
		return varRef, false
	}

	// read source file
	fileBytes, err := ioutil.ReadFile(r.moduleCall.SourceAddrRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", r.moduleCall.SourceAddr, err)
		return varRef, false
	}

	// extract values from attribute expressions as golang interface{}
	c := converter{bytes: fileBytes}
	val, err := c.convertExpression(varAttr.Expr)
	if err != nil {
		return varRef, false
	}
	zap.S().Debugf("resolved variable reference from module call; var ref: '%v', value: '%v', type: '%v'", varRef, val, reflect.TypeOf(val))

	// replace the variable reference string with actual value
	if reflect.TypeOf(val).Kind() == reflect.String {

		valStr := val.(string)

		// if resolved value is a reference eg. "${module.somevalue}",
		// then return varRef
		if isRef(valStr) {
			return varRef, false
		}

		resolvedVal := varRefPattern.ReplaceAll([]byte(varRef), []byte(valStr))
		return string(resolvedVal), true
	}

	// return extracted value
	return val, true
}

// getVarValue returns the variable value as configured in IaC config in module
func (r *RefResolver) getVarValue(varRef string) (interface{}, bool) {

	// get variable name from varRef
	varName := getVarName(varRef)

	// check if variable name exists in the map of variables read from IaC
	hclVar, present := r.variables[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in config", varName, varRef)
		return varRef, false
	}

	// return varRef if default value is not present, or value is a NilVal,
	// or if default value is not known
	if hclVar.Default.IsNull() || hclVar.Default.RawEquals(cty.NilVal) || !hclVar.Default.IsKnown() {
		return varRef, false
	}

	// default value is of cty.Value type, convert it to native golang type
	// based on cty.Type, determine golang type
	for _, converter := range ctyConverterFuncs {
		if val, err := converter(hclVar.Default); err == nil {
			zap.S().Debugf("resolved variable reference; var ref: '%v', value: '%v', type: '%v'", varRef, val, reflect.TypeOf(val))

			// replace the variable reference string with actual value
			if reflect.TypeOf(val).Kind() == reflect.String {
				valStr := val.(string)
				resolvedVal := varRefPattern.ReplaceAll([]byte(varRef), []byte(valStr))
				return string(resolvedVal), true
			}
			return val, true
		}
	}
	zap.S().Debugf("failed to convert cty.Value '%v' to golang native value", hclVar.Default.GoString())

	// return original reference if cty conversion attempts fail
	return varRef, false
}

// isRef returns true if the given string is a variable reference
func isRef(attrVal string) bool {
	return refPattern.MatchString(attrVal)
}

// isVarRef returns true if the given string is a variable reference
func isVarRef(attrVal string) bool {
	return varRefPattern.MatchString(attrVal)
}

// getVarName returns the actual variable name as configured in IaC. It trims
// of "${var." prefix and "}" suffix and returns the variable name
func getVarName(varRef string) string {

	// 1. split at "${var.", remove everything before
	split := strings.Split(varRef, varPrefix)
	varName := split[1]

	// 2. split at "}", remove everything after
	split = strings.Split(varName, varSuffix)
	varName = split[0]

	zap.S().Debugf("extracted variable name %q from reference %q", varName, varRef)
	return varName
}
