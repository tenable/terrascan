/*
    Copyright (C) 2022 Tenable, Inc.

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

package commons

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// list of available primitive cty to golang type converters
var (
	ctyNativeConverterFuncs = []func(cty.Value) (interface{}, error){ctyToStr, ctyToInt, ctyToFloat, ctyToBool}
)

// ctyToStr tries to convert the given cty.Value into golang string type
func ctyToStr(ctyVal cty.Value) (interface{}, error) {
	var val string
	err := gocty.FromCtyValue(ctyVal, &val)
	return val, err
}

// ctyToInt tries to convert the given cty.Value into golang int type
func ctyToInt(ctyVal cty.Value) (interface{}, error) {
	var val int
	err := gocty.FromCtyValue(ctyVal, &val)
	return val, err
}

// ctyToFloat tries to convert the given cty.Value into golang float type
func ctyToFloat(ctyVal cty.Value) (interface{}, error) {
	var val float64
	err := gocty.FromCtyValue(ctyVal, &val)
	return val, err
}

// ctyToBool tries to convert the given cty.Value into golang bool type
func ctyToBool(ctyVal cty.Value) (interface{}, error) {
	var val bool
	err := gocty.FromCtyValue(ctyVal, &val)
	return val, err
}

// ctyToSlice tries to convert the given cty.Value into golang slice of
// interface{}
func ctyToSlice(ctyVal cty.Value) (interface{}, error) {
	var val []interface{}
	var allErrs error

	if ctyVal.Type().IsListType() || ctyVal.Type().IsTupleType() || ctyVal.Type().IsSetType() {
		for _, v := range ctyVal.AsValueSlice() {
			nativeVal, err := convertCtyToGoNative(v)
			if err != nil {
				allErrs = errors.Wrap(allErrs, err.Error())
				continue
			}
			val = append(val, nativeVal)
		}
		if allErrs != nil {
			return nil, allErrs
		}
		return val, nil
	}
	return val, fmt.Errorf("incorrect type")
}

// ctyToMap tries to converts the incoming cty.Value into map[string]cty.Value
// then for every key value of this map, tries to convert the cty.Value into
// native golang value and create a new map[string]interface{}
func ctyToMap(ctyVal cty.Value) (interface{}, error) {

	if !(ctyVal.Type().IsMapType() || ctyVal.Type().IsObjectType()) {
		return nil, fmt.Errorf("not map type")
	}

	var (
		ctyValMap = ctyVal.AsValueMap() // map[string]cty.Value
		val       = make(map[string]interface{})
		allErrs   error
	)

	// cannot process an empty ctValMap
	if len(ctyValMap) < 1 {
		return nil, fmt.Errorf("empty ctyValMap")
	}

	// iterate over every key cty.Value pair, try to convert cty.Value into
	// golang value
	for k, v := range ctyValMap {
		// convert cty.Value to native golang type based on cty.Type
		nativeVal, err := convertCtyToGoNative(v)
		if err != nil {
			allErrs = errors.Wrap(allErrs, err.Error())
			continue
		}
		val[k] = nativeVal
	}
	if allErrs != nil {
		return nil, allErrs
	}

	// hopefully successful!
	return val, nil
}

// convertCtyToGoNative converts a cty.Value to its go native type
func convertCtyToGoNative(ctyVal cty.Value) (interface{}, error) {
	// no need to convert variable to any type in case value is null
	// added check here since this function is been called in recursive manner
	if ctyVal.IsNull() {
		return nil, nil
	}
	if ctyVal.Type().IsPrimitiveType() {
		return convertPrimitiveType(ctyVal)
	}
	return convertComplexType(ctyVal)
}

// convertPrimitiveType converts a primitive cty.Value to its go native type
func convertPrimitiveType(ctyVal cty.Value) (interface{}, error) {
	for _, converter := range ctyNativeConverterFuncs {
		if val, err := converter(ctyVal); err == nil {
			return val, err
		}
	}
	return nil, fmt.Errorf("ctyVal could not be resolved to native go type")
}

// convertPrimitiveType converts a complex cty.Value (list, tuple, set, map, object) to its go native type
func convertComplexType(ctyVal cty.Value) (interface{}, error) {
	if ctyVal.Type().IsListType() || ctyVal.Type().IsTupleType() || ctyVal.Type().IsSetType() {
		return ctyToSlice(ctyVal)
	}
	return ctyToMap(ctyVal)
}
