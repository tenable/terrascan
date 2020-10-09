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
	// reference patterns
	moduleRefPattern = regexp.MustCompile(`\$\{module\.(.*)\.(.*)}`)
)

// isModuleRef return true if the given string has a cross module reference
func isModuleRef(ref string) bool {
	return moduleRefPattern.MatchString(ref)
}

// getModuleVarName extracts and returns the module and variable name from the
// module reference string
func getModuleVarName(moduleRef string) (string, string) {

	var (
		modulePrefix = "${module."
		moduleSuffix = "}"
	)

	// ex of moduleRef: ${module.name.variable}
	// 1. split at "${var.", remove everything before
	split := strings.Split(moduleRef, modulePrefix)
	mod := split[1]

	// 2. split at "}", remove everything after
	split = strings.Split(mod, moduleSuffix)
	mod = split[0]

	// 3. split at "."; eg: "name.variable"
	split = strings.Split(mod, ".")
	if len(split) < 2 {
		return "", ""
	}

	// return module name and variable name
	return split[0], split[1]
}

// ResolveModuleRef tries to resolve cross module references
func (r *RefResolver) ResolveModuleRef(moduleRef string,
	children map[string]*hclConfigs.Config) interface{} {

	// get module and variable name
	moduleName, varName := getModuleVarName(moduleRef)

	// module and variable names cannot be empty
	if moduleName == "" || varName == "" {
		return moduleRef
	}

	// get module config
	module, ok := children[moduleName]
	if !ok {
		zap.S().Debugf("module: '%v' not present in children", moduleName)
		return moduleRef
	}

	// check if variable name exists in the map of variables read from the
	// referenced module
	hclVar, present := module.Module.Variables[varName]
	if !present {
		zap.S().Debugf("variable name: %q, ref: %q not present in config", varName, moduleRef)
		return moduleRef
	}

	// return moduleRef if default value is not present, or value is a NilVal,
	// or if default value is not known
	if hclVar.Default.IsNull() || hclVar.Default.RawEquals(cty.NilVal) || !hclVar.Default.IsKnown() {
		return moduleRef
	}

	// default value is of cty.Value type, convert it to native golang type
	// based on cty.Type, determine golang type
	for _, converter := range ctyConverterFuncs {
		if val, err := converter(hclVar.Default); err == nil {
			zap.S().Debugf("resolved module variable ref: '%v', value: '%v'", moduleRef, val)

			// replace the variable reference string with actual value
			if reflect.TypeOf(val).Kind() == reflect.String {
				valStr := val.(string)
				resolvedVal := varRefPattern.ReplaceAll([]byte(moduleRef), []byte(valStr))
				return string(resolvedVal)
			}
			return val
		}
	}
	zap.S().Debugf("failed to convert cty.Value '%v' to golang native value", hclVar.Default.GoString())

	return moduleRef
}
