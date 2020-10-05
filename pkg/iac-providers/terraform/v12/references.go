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
	"reflect"
	"regexp"
	"strings"

	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

var (
	// variable reference pattern
	varRefPattern = regexp.MustCompile(`\$\{var\..*\}`)
	varPrefix     = "${var."
	varSuffix     = "}"
)

// ResolveRefs figures out all the variable references in the resource config
// and tries to replace the variable references, if possible,
// with actual value
func ResolveRefs(config jsonObj, variables map[string]*hclConfigs.Variable) jsonObj {

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
				config[k] = getVarValue(v.(string), variables)
			}
		} else if valType == "tfv12.jsonObj" && valKind == reflect.Map {
			// case 2: config value is of type jsonObj
			config[k] = ResolveRefs(v.(jsonObj), variables)
		} else if valType == "[]tfv12.jsonObj" && valKind == reflect.Slice {
			// case 3: config value is of type []jsonObj

			// type assert interface{} -> []jsonObj
			sConfig, ok := v.([]jsonObj)
			if !ok {
				continue
			}

			// iterate over the []jsonObj to resolve refs
			for i, c := range sConfig {
				sConfig[i] = ResolveRefs(c, variables)
			}
			config[k] = sConfig

		} else {
			zap.S().Debugf("cannot resolve refs for var: '%v', type '%v', kind: '%v',\nvalue: '%v'\n\n", k, reflect.TypeOf(v), valKind, v)
		}
	}

	return config
}

// getVarValue returns the variable value as configured in IaC config
func getVarValue(varRef string, variables map[string]*hclConfigs.Variable) interface{} {

	// get variable name from varRef
	varName := getVarName(varRef)

	// check if variable name exists in the map of variables read from IaC
	hclVar, present := variables[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in config", varName, varRef)
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
