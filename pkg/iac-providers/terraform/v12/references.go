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
	// reference patterns
	refPattern       = regexp.MustCompile(`\$\{.*\}`)
	moduleRefPattern = regexp.MustCompile(`\$\{module\.(.*)\.(.*)}`)
	varRefPattern    = regexp.MustCompile(`\$\{var\..*\}`)
)

// RefResolver tries to resolve all the references in the given terraform
// config
type RefResolver struct {
	variables        map[string]*hclConfigs.Variable
	parentModuleCall *hclConfigs.ModuleCall
	parentChildren   map[string]*hclConfigs.Config
	children         map[string]*hclConfigs.Config
}

// NewRefResolver returns a new RefResolver struct
func NewRefResolver(variables map[string]*hclConfigs.Variable,
	parentModuleCall *hclConfigs.ModuleCall,
	parentChildren map[string]*hclConfigs.Config,
	children map[string]*hclConfigs.Config) *RefResolver {
	return &RefResolver{
		variables:        variables,
		parentModuleCall: parentModuleCall,
		parentChildren:   parentChildren,
		children:         children,
	}
}

// ResolveRefs figures out all the variable references in the resource config
// and tries to replace the variable references, if possible,
// with actual value
func (r *RefResolver) ResolveRefs(config jsonObj) jsonObj {

	// iterate over every attribute in the config to resolve references
	for k, v := range config {

		var (
			valKind = reflect.TypeOf(v).Kind()
			valType = reflect.TypeOf(v).String()
		)

		switch {
		case valKind == reflect.String:

			// case 1: config value is a string; in resource config, refs
			// are of the type string
			config[k] = r.ResolveStrRef(v.(string))

		case valType == "tfv12.jsonObj" && valKind == reflect.Map:

			// case 2: config value is of type jsonObj
			config[k] = r.ResolveRefs(v.(jsonObj))

		case valType == "[]tfv12.jsonObj" && valKind == reflect.Slice:

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

		default:
			zap.S().Debugf("not processing attribute name: '%v', value: '%v' for resolving references", k, v)
		}
	}

	return config
}

// ResolveStrRef tries to resolve a string reference. Reference can be a
// variable "${var.foo}", cross module variable "${module.foo.bar}", local
// value "${local.foo}"
func (r *RefResolver) ResolveStrRef(ref string) interface{} {

	switch {
	case isVarRef(ref):

		/*
			Variable values initialized in the parent module call take
			precedence over values initialized in the same module,
			hence we execute the following steps:
			- resolve variables initialized in parent module call
			- resolve variables initialized in same module
		*/

		// 1. resolve variables initialized in parent module call
		val := r.ResolveVarRefFromParentModuleCall(ref)

		if reflect.TypeOf(val).Kind() == reflect.String {

			valStr := val.(string)
			/*
			 Now, if the output of ResolveVarRefFromParentModuleCall is a string
			 then there are following possibilties, output can be a:
			 - variable reference (resolve the variable reference)
			 - cross module variable reference in parent module call
			 	(resolve module reference)
			 - some other reference
			 - an absolute value (return from here)
			*/

			switch {
			case isVarRef(valStr):

				// resolve variable reference
				return r.ResolveVarRef(valStr)

			case isModuleRef(valStr):

				// resolve cross module reference in parent module
				return r.ResolveModuleRef(valStr, r.parentChildren)

			case isRef(valStr):

				// some other reference
				return r.ResolveStrRef(valStr)
			}
		}

		// hopefully, the variable has been resolved here
		return val

	case isModuleRef(ref):

		// resolve cross module references
		return r.ResolveModuleRef(ref, r.children)

	default:
		return ref
	}
}

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
		zap.S().Debugf("resolved str variable ref: '%v', value: '%v'", varRef, resolvedVal)
		return string(resolvedVal)
	}

	// return extracted value
	zap.S().Debugf("resolved variable ref: '%v', value: '%v'", varRef, val)
	return val
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
