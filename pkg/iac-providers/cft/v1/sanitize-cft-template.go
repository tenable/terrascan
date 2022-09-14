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

package cftv1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/awslabs/goformation/v6/cloudformation"
	"github.com/awslabs/goformation/v6/cloudformation/policies"
	"github.com/awslabs/goformation/v6/intrinsics"
	"go.uber.org/zap"
)

func (a *CFTV1) sanitizeCftTemplate(data []byte, isYAML bool) (map[string]interface{}, error) {
	var (
		intrinsified []byte
		err          error
	)

	if isYAML {
		// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
		intrinsified, err = intrinsics.ProcessYAML(data, nil)
		if err != nil {
			return nil, fmt.Errorf("error while resolving intrinsic functions, error %w", err)
		}
	} else {
		// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
		intrinsified, err = intrinsics.ProcessJSON(data, nil)
		if err != nil {
			return nil, fmt.Errorf("error while resolving intrinsic functions, error %w", err)
		}
	}

	templateFileMap := make(map[string]interface{})

	err = json.Unmarshal(intrinsified, &templateFileMap)
	if err != nil {
		return nil, err
	}

	// sanitize Parameters
	params, ok := templateFileMap["Parameters"]
	if ok {
		pMap, ok := params.(map[string]interface{})
		if ok {
			for pName := range pMap {
				zap.S().Debug(fmt.Sprintf("inspecting parameter '%s'", pName))
				inspectAndSanitizeParameters(pMap[pName])
			}
		}
	}

	// sanitize resources
	r, ok := templateFileMap["Resources"]
	if ok {
		rMap, ok := r.(map[string]interface{})
		if ok {
			for rName := range rMap {
				zap.S().Debug("inspecting resource", zap.String("Resource Name", rName))
				if shouldRemoveResource := inspectAndSanitizeResource(rMap[rName]); shouldRemoveResource {
					// we would remove any resource from the map for which goformation doesn't have a type defined
					delete(rMap, rName)
				}
			}
		}
	}

	return templateFileMap, nil
}

func inspectAndSanitizeParameters(p interface{}) {
	paramMap, ok := p.(map[string]interface{})
	if !ok {
		zap.S().Debug("invalid data for 'Parameters', should be of type map[string]interface{}")
		return
	}
	structFieldsMap := examineStruct(reflect.TypeOf(cloudformation.Parameter{}))
	if structFieldsMap != nil {
		for paramName := range paramMap {
			v, ok := structFieldsMap[paramName]
			if !ok {
				zap.S().Debug(fmt.Sprintf("attribute '%s', not present in 'Parameter' struct fields", paramName))
				continue
			}
			val := fixWithType(paramMap[paramName], v.Type)
			if val != nil {
				paramMap[paramName] = val
			}
		}
	}
}

func inspectAndSanitizeResource(r interface{}) (shouldRemoveResource bool) {
	resMap, ok := r.(map[string]interface{})
	if !ok {
		zap.S().Debug("invalid data for 'Resource', should be of type map[string]interface{}")
		return
	}

	// get the type of the resource
	t, ok := resMap["Type"]
	if !ok {
		zap.S().Debug("resource must have an attribute 'Type'")
		return
	}

	tVal, ok := t.(string)
	if !ok {
		zap.S().Debug("attribute 'Type' should be a string")
		return
	}

	goformationCftObj, ok := cloudformation.AllResources()[tVal]
	if !ok {
		shouldRemoveResource = true
		zap.S().Debug(fmt.Sprintf("not goformation resource present for '%s'", tVal))
		return
	}

	cftObjType := reflect.TypeOf(goformationCftObj)
	// if the object is of pointer type, get type of its concrete value
	if cftObjType.Kind() == reflect.Ptr {
		cftObjType = cftObjType.Elem()
	}
	structFieldsMap := examineStruct(cftObjType)
	if structFieldsMap != nil {
		// sanitize the properties of the resource
		prop, ok := resMap["Properties"]
		if !ok {
			zap.S().Debug("resource doesn't have 'Properties'")
			return
		}

		propMap, ok := prop.(map[string]interface{})
		if !ok {
			zap.S().Debug("'Properties' should be of type map[string]interface{}")
			return
		}

		for propName := range propMap {
			structField, ok := structFieldsMap[propName]
			if !ok {
				zap.S().Debug(fmt.Sprintf("attribute '%s', not present in '%s' struct fields", propName, tVal))
				continue
			}
			val := fixWithType(propMap[propName], structField.Type)
			if val != nil {
				propMap[propName] = val
			}
		}

		inspectAndSanitizeResourceAttributes(resMap)
	}
	return
}

func inspectAndSanitizeResourceAttributes(resource map[string]interface{}) {
	// every cft resource has 6 attributes as specified at https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-product-attribute-reference.html

	// sanitize CreationPolicy if present (CreationPolicy is an object)
	cp, ok := resource["CreationPolicy"]
	if ok {
		cpMap, ok := cp.(map[string]interface{})
		if ok {
			structFieldsMap := examineStruct(reflect.TypeOf(policies.CreationPolicy{}))
			for k := range cpMap {
				v, ok := structFieldsMap[k]
				if !ok {
					zap.S().Debug(fmt.Sprintf("attribute '%s' not present 'CreationPolicy' struct", k))
					continue
				}
				val := fixWithType(cpMap[k], v.Type)
				if val != nil {
					cpMap[k] = val
				}
			}
		}
	}

	// sanitize UpdatePolicy if present (UpdatePolicy is an object)
	up, ok := resource["UpdatePolicy"]
	if ok {
		upMap, ok := up.(map[string]interface{})
		if ok {
			structFieldsMap := examineStruct(reflect.TypeOf(policies.UpdatePolicy{}))
			for k := range upMap {
				v, ok := structFieldsMap[k]
				if !ok {
					zap.S().Debug(fmt.Sprintf("attribute '%s' not present 'UpdatePolicy' struct", k))
					continue
				}
				val := fixWithType(upMap[k], v.Type)
				if val != nil {
					upMap[k] = val
				}
			}
		}
	}

	// sanitize DependsOn if present (DependsOn is a slice)
	d, ok := resource["DependsOn"]
	if ok {
		// check if DependsOn is a slice
		_, ok = d.([]interface{})
		if !ok {
			newVal := make([]interface{}, 0)
			newVal = append(newVal, d)
			resource["DependsOn"] = newVal
		}
	}

	// Metadata is of type map[string]interface{}, we do not need to sanitize
	// DeletionPolicy is of type string, we do not need to sanitize
	// UpdateReplacePolicy is of type string, we do not need to sanitize
}

// fixWithType... tries to fix the orignal value based on type specified
// it doesn't try to fix, if type of original data is the type specified
func fixWithType(data interface{}, r reflect.Type) interface{} {
	switch t := data.(type) {
	case int, int8, int16, int32, int64:
		val := t.(int)
		switch r.Kind() {
		case reflect.Float32, reflect.Float64:
			return float64(val)
		case reflect.String:
			return strconv.Itoa(val)
		case reflect.Ptr:
			return fixWithType(data, r.Elem())
		}
	case string:
		switch r.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(t, 10, 64)
			if err == nil {
				return v
			}
		case reflect.Float32:
			v, err := strconv.ParseFloat(t, 32)
			if err == nil {
				return v
			}
		case reflect.Float64:
			v, err := strconv.ParseFloat(t, 64)
			if err == nil {
				return v
			}
		case reflect.Bool:
			v, err := strconv.ParseBool(t)
			if err == nil {
				return v
			}
		case reflect.Ptr:
			return fixWithType(data, r.Elem())
		}
	case bool:
		switch r.Kind() {
		case reflect.String:
			return strconv.FormatBool(t)
		case reflect.Ptr:
			return fixWithType(data, r.Elem())
		}

	case float32, float64:
		val := t.(float64)
		switch r.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(val)
		case reflect.String:
			return strconv.FormatFloat(val, 'f', -1, 64)
		case reflect.Ptr:
			return fixWithType(data, r.Elem())
		}

	case []interface{}:
		switch r.Kind() {
		case reflect.Array, reflect.Slice, reflect.Ptr:
			arr := []interface{}{}
			for x := range t {
				v := fixWithType(t[x], r.Elem())
				if v != nil {
					arr = append(arr, v)
				} else {
					arr = append(arr, t[x])
				}
			}
			return arr
		}
	case map[string]interface{}:
		switch r.Kind() {
		case reflect.Struct:
			sType := reflect.New(r).Type().Elem()
			mMap := examineStruct(sType)
			for k := range t {
				v, ok := mMap[k]
				if !ok {
					zap.S().Debug(fmt.Sprintf("attribute '%s' not present in struct '%s'", k, sType.String()))
					continue
				}
				val := fixWithType(t[k], v.Type)
				if val != nil {
					t[k] = val
				}
			}
			return t
		case reflect.Ptr:
			sType := reflect.New(r).Type().Elem().Elem()
			mMap := examineStruct(sType)
			for k := range t {
				v, ok := mMap[k]
				if !ok {
					zap.S().Debug(fmt.Sprintf("attribute '%s' not present in struct '%s'", k, sType.String()))
					continue
				}
				val := fixWithType(t[k], v.Type)
				if val != nil {
					t[k] = val
				}
			}
			return t
		}
	}
	return nil
}

func examineStruct(t reflect.Type) map[string]reflect.StructField {
	if t.Kind() != reflect.Struct {
		return nil
	}
	m := make(map[string]reflect.StructField)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		key := f.Name
		// we want to get the tag name in the struct definition
		// struct field name may be different than tag name
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			if i := strings.Index(tag, ","); i != -1 {
				tag = tag[:strings.Index(tag, ",")]
			}
			key = tag
		}
		m[key] = f
	}
	return m
}
