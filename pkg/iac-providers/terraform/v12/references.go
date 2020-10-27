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

	hclConfigs "github.com/hashicorp/terraform/configs"
)

var (
	// reference patterns
	refPattern = regexp.MustCompile(`\$\{.*\}`)
)

// RefResolver tries to resolve all the references in the given terraform
// config
type RefResolver struct {
	Config           *hclConfigs.Config
	ParentModuleCall *hclConfigs.ModuleCall
}

// NewRefResolver returns a new RefResolver struct
func NewRefResolver(config *hclConfigs.Config,
	parentModuleCall *hclConfigs.ModuleCall) *RefResolver {
	return &RefResolver{
		Config:           config,
		ParentModuleCall: parentModuleCall,
	}
}

// isRef returns true if the given string is a variable reference
func isRef(val string) bool {
	return refPattern.MatchString(val) || isVarRef(val) || isModuleRef(val) || isLocalRef(val) || isLookupRef(val)
}

// ResolveRefs figures out all the variable references in the resource config
// and tries to replace the variable references, if possible,
// with actual value
func (r *RefResolver) ResolveRefs(config jsonObj) jsonObj {

	// iterate over every attribute in the config to resolve references
	for k, v := range config {

		var (
			vKind = reflect.TypeOf(v).Kind()
			vType = reflect.TypeOf(v).String()
		)

		switch {
		case vKind == reflect.String:

			// case 1: config value is a string; in resource config, refs
			// are of the type string
			config[k] = r.ResolveStrRef(v.(string))

		case vType == "tfv12.jsonObj" && vKind == reflect.Map:

			// case 2: config value is of type jsonObj
			config[k] = r.ResolveRefs(v.(jsonObj))

		case vType == "[]interface {}" && vKind == reflect.Slice:

			// case 3: config value is a []interface{}
			sConfig, ok := v.([]interface{})
			if !ok {
				continue
			}

			// if the golang native type is string, then try and resolve
			// references
			if len(sConfig) > 0 && reflect.TypeOf(sConfig[0]).Kind() == reflect.String {
				for i, c := range sConfig {
					sConfig[i] = r.ResolveStrRef(c.(string))
				}
				config[k] = sConfig
			}

		case vType == "[]tfv12.jsonObj" && vKind == reflect.Slice:

			// case 4: config value is of type []jsonObj

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

// ResolveStrRef tries to resolve a string reference. Reference can be a
// variable "${var.foo}", cross module variable "${module.foo.bar}", local
// value "${local.foo}"
func (r *RefResolver) ResolveStrRef(ref string) interface{} {

	switch {
	case isModuleRef(ref):

		// resolve cross module references
		return r.ResolveModuleRef(ref, r.Config.Children)

	case isLocalRef(ref):

		// resolve local value references
		return r.ResolveLocalRef(ref)

	case isLookupRef(ref):

		// resolve lookup references
		return r.ResolveLookupRef(ref)

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
				return r.ResolveModuleRef(valStr, r.Config.Parent.Children)

			case isRef(valStr):

				// some other reference
				return r.ResolveStrRef(valStr)
			}
		}

		// hopefully, the variable has been resolved here
		return val

	default:
		return ref
	}
}
