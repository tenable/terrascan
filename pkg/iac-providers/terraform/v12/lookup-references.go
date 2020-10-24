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

	"go.uber.org/zap"
)

var (
	// reference patterns
	lookupRefPattern = regexp.MustCompile(`(\$\{)?lookup\((?P<table>\S+)\,\s*?(?P<key>\S+)\)(\})?`)
)

// isLookupRef returns true if the given string is a lookup value reference
func isLookupRef(attrVal string) bool {
	return lookupRefPattern.MatchString(attrVal)
}

// getLookupName returns the actual lookup value name as configured in IaC. It
// trims of "${lookup(." prefix and ")}" suffix and returns the lookup value name
func getLookupName(lookupRef string) (string, string, string) {

	// 1. extract the exact lookup value reference from the string
	lookupExpr := lookupRefPattern.FindString(lookupRef)

	// 2. extract lookup value name from lookup value reference
	match := lookupRefPattern.FindStringSubmatch(lookupRef)
	result := make(map[string]string)
	for i, name := range lookupRefPattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	table := result["table"]
	key := result["key"]

	zap.S().Debugf("extracted lookup table %q key %q from reference %q", table, key, lookupRef)
	return table, key, lookupExpr
}

// ResolveLookupRef returns the lookup value as configured in IaC config in module
func (r *RefResolver) ResolveLookupRef(lookupRef string) interface{} {

	// get lookup name from lookupRef
	table, key, _ := getLookupName(lookupRef)

	// resolve key, if it is a reference
	resolvedKey := r.ResolveStrRef(key)

	// check if key is still an unresolved reference
	if reflect.TypeOf(resolvedKey).Kind() == reflect.String && isRef(resolvedKey.(string)) {
		zap.S().Debugf("failed to resolve key ref: '%v'", key)
		return lookupRef
	}

	// resolve table, if it is a ref
	lookup := r.ResolveStrRef(table)

	// check if lookup is a map
	if reflect.TypeOf(lookup).String() != "map[string]interface {}" {
		zap.S().Debugf("failed to resolve lookup ref %q, table name %q into a map, received %v", lookupRef, table, reflect.TypeOf(lookup).String())
		return lookupRef
	}

	// check if key is present in lookup table
	resolved, ok := lookup.(map[string]interface{})[resolvedKey.(string)]
	if !ok {
		zap.S().Debugf("key %q not present in lookup table %v", key, lookup)
		return lookupRef
	}

	zap.S().Debugf("resolved lookup ref %q to value %v", lookupRef, resolved)
	return resolved
}
