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

	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/policies"
	"github.com/awslabs/goformation/v7/intrinsics"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// PARAMETERS is a constant to fetch Parameters from CFT
const PARAMETERS = "Parameters"

// RESOURCES is a constant to fetch Resources from CFT
const RESOURCES = "Resources"

func (a *CFTV1) sanitizeCftTemplate(fileName string, data []byte, isYAML bool) (map[string]interface{}, error) {
	var (
		intrinsified []byte
		err          error
	)
	if isYAML {
		fallbackToDefaultProcessing := false

		for i := 0; i < 1; i++ {
			// convert the yaml into json
			jsonData, err := a.ReadYAMLFileIntoJSON(fileName)
			if err == nil {
				jsonData, err := a.resolveResourceIDs(jsonData)
				if err != nil {
					zap.S().Debug(fmt.Sprintf("error while resolving Resource IDs, error %s", err.Error()))
					fallbackToDefaultProcessing = true
					break
				}
				intrinsified, err = intrinsics.ProcessJSON(jsonData, nil)
				if err != nil {
					zap.S().Debug(fmt.Sprintf("error while resolving Resource IDs, error %s", err.Error()))
					fallbackToDefaultProcessing = true
					break
				}
			}
		}
		if fallbackToDefaultProcessing || len(intrinsified) == 0 { // fallback to default behaviour of yaml processing
			data, err = removeRefAnchors(data)
			if err != nil {
				return nil, err
			}
			// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
			intrinsified, err = intrinsics.ProcessYAML(data, nil)
			if err != nil {
				return nil, fmt.Errorf("error while resolving intrinsic functions, error %w", err)
			}
		}
	} else {
		jsonData, err := a.resolveResourceIDs(data)
		if err == nil {
			intrinsified, err = intrinsics.ProcessJSON(jsonData, nil)
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
	}
	templateFileMap := make(map[string]interface{})

	err = json.Unmarshal(intrinsified, &templateFileMap)
	if err != nil {
		return nil, err
	}

	// sanitize Parameters
	params, ok := templateFileMap[PARAMETERS]
	var pMap map[string]interface{}
	pMapconverted := make(map[string]interface{})
	if ok {
		pMap, ok = params.(map[string]interface{})
		if ok {
			for pName := range pMap {
				zap.S().Debug(fmt.Sprintf("inspecting parameter '%s'", pName))
				inspectAndSanitizeParameters(pMap[pName])
				resultMap, found := convertFloat64ToString(pMap[pName])
				if found {
					pMapconverted[pName] = resultMap
				}
			}
		}
	}

	// sanitize resources
	r, ok := templateFileMap[RESOURCES]
	if ok {
		rMap, ok := r.(map[string]interface{})
		if ok {
			for rName := range rMap {
				zap.S().Debug("inspecting resource", zap.String("Resource Name", rName))
				if shouldRemoveResource := inspectAndSanitizeResource(rMap[rName], pMapconverted); shouldRemoveResource {
					// we would remove any resource from the map for which goformation doesn't have a type defined
					delete(rMap, rName)
				}
			}
		}
	}

	return templateFileMap, nil
}

func removeRefAnchors(data []byte) ([]byte, error) {
	const REF = "!ref"
	const DoubleColon = "::"
	strdata := string(data)
	words := strings.Split(strdata, " ")

	var cfnmap map[any]any
	err := yaml.Unmarshal(data, &cfnmap)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling yaml, error %w", err)
	}

	cfnJSONMap := anyMapToStringMap(cfnmap)
	paramsMap, paramsOk := cfnJSONMap[PARAMETERS].(map[string]any)

	for i := range words {
		current := strings.ToLower(words[i])
		if len(words) == i+1 {
			break
		}

		if strings.Contains(current, REF) {
			next := strings.TrimSpace(words[i+1])
			nextLower := strings.ToLower(words[i+1])

			if strings.Contains(nextLower, "aws::") {
				if i+1 < len(words) { // check edge case
					words[i+1] = strings.ReplaceAll(nextLower, DoubleColon, "-")
				}
				continue
			}

			if paramsOk {
				if _, ok := paramsMap[next]; ok {
					continue
				}
			}
		}

		if strings.Contains(current, REF) {
			words[i] = strings.Replace(current, REF, "", 1)
		}
	}

	strdata = strings.Join(words, " ")
	return []byte(strdata), nil
}

func anyMapToStringMap(cfnmap map[any]any) map[string]any {
	res := map[string]any{}
	for k, v := range cfnmap {
		switch v2 := v.(type) {
		case map[any]any:
			res[fmt.Sprint(k)] = anyMapToStringMap(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
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

func inspectAndSanitizeResource(r interface{}, pMap map[string]interface{}) (shouldRemoveResource bool) {
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
			findKeyAndReplace(propMap[propName], pMap)
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

func convertFloat64ToString(paramMap interface{}) (map[string]interface{}, bool) {
	valMapNew := make(map[string]interface{})
	foundfloat := false
	if valMap, ok := paramMap.(map[string]interface{}); ok {
		for paramName := range valMap {
			var newBound []string
			valToCheck := valMap[paramName]
			switch val := valToCheck.(type) {
			case int, float64, int32, float32, int8, int16, int64:
				valToCheck = fmt.Sprintf("%v", val)
				foundfloat = true
				valMapNew[paramName] = valToCheck
			}
			if arrayValue, ok := valToCheck.([]interface{}); ok {
				newBound = make([]string, len(arrayValue))
				for i := range arrayValue {
					switch val := arrayValue[i].(type) {
					case int, float64, int32, float32, int8, int16, int64:
						newBound[i] = fmt.Sprintf("%v", val)
						foundfloat = true
					}
				}
				if foundfloat {
					valMapNew[paramName] = newBound
				}
			}
		}
		return valMapNew, foundfloat
	}
	return valMapNew, foundfloat
}

// findKeyAndReplace key in interface (recursively) and return value as interface
func findKeyAndReplace(obj interface{}, propValues map[string]interface{}) (interface{}, bool) {
	// if the argument is not a map, ignore it
	mobj, ok := obj.(map[string]interface{})
	if !ok {
		return nil, false
	}
	for k, v := range mobj {
		// key match, return value
		if val, ok := propValues[k]; ok {
			if valMap, ok := val.(map[string]interface{}); ok {
				if val2, ok := valMap["Default"]; ok {
					mobj[k] = val2
				}
				return v, true
			}
		}
		// if the value is a map, search recursively
		if m, ok := v.(map[string]interface{}); ok {
			if res, ok := findKeyAndReplace(m, propValues); ok {
				return res, true
			}
		}
		// if the value is an array, search recursively
		if va, ok := v.([]interface{}); ok {
			for _, a := range va {
				if res, ok := findKeyAndReplace(a, propValues); ok {
					return res, true
				}
			}
		}
	}
	return nil, false
}

// ReadYAMLFileIntoJSON converts the given file into JSON string
func (a *CFTV1) ReadYAMLFileIntoJSON(fileName string) ([]byte, error) {
	templateSample, err := a.File(fileName)
	if err != nil {
		return nil, err
	}
	gostruct, err := templateSample.Map()
	if err != nil {
		zap.S().Errorf("failed to map yaml to json. error : %s in file %s", err.Error(), fileName)
		return nil, err
	}
	jsonData, err := json.Marshal(gostruct)
	if err != nil {
		zap.S().Errorf("failed to convert yaml to json. error : %s", err.Error())
		return nil, err
	}
	return jsonData, nil
}
func (a *CFTV1) getMapOfResourceIds(allData interface{}) map[string]string {
	mapOfresourceIds := make(map[string]string)
	mapOfParameters := make(map[string]interface{})
	if templateFileMap, ok := allData.(map[string]interface{}); ok {
		r, ok := templateFileMap[PARAMETERS]
		if ok {
			rMap, ok := r.(map[string]interface{})
			if ok {
				for rName, val := range rMap {
					zap.S().Debug("inspecting resource", zap.String("Parameters Name", rName))
					if val1, ok := val.(map[string]interface{}); ok {
						mapOfParameters[rName] = val1["Default"]
					}
				}
			}
		}
	}

	if templateFileMap, ok := allData.(map[string]interface{}); ok {
		r, ok := templateFileMap[RESOURCES]
		if ok {
			rMap, ok := r.(map[string]interface{})
			if ok {
				for rName := range rMap {
					zap.S().Debug("inspecting resource", zap.String("Resource Name", rName))
					if _, ok := mapOfParameters[rName]; !ok {
						mapOfresourceIds[rName] = rName
					}
				}
			}
		}
	}
	return mapOfresourceIds
}

// resolveResourceIDs resolves the indirect resource to resource references
func (a *CFTV1) resolveResourceIDs(jsonData []byte) ([]byte, error) {
	var unmarshalled interface{}
	if err := json.Unmarshal(jsonData, &unmarshalled); err != nil {
		return nil, fmt.Errorf("invalid JSON: %s", err)
	}
	mapOfparentReferences := a.getMapOfResourceIds(unmarshalled)
	unmarshalledResult := a.resolveIndirectReferences(nil, unmarshalled, "", mapOfparentReferences)
	resultBytes, err := json.Marshal(unmarshalledResult)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %s", err)
	}
	jsonData = resultBytes
	return jsonData, nil
}

// resolveIndirectReferences finds if the references are of the Name of the resource and replace it with actual Name
func (a *CFTV1) resolveIndirectReferences(parent interface{}, input interface{}, parentKey string, mapOfReferences map[string]string) interface{} {

	switch value := input.(type) {

	case map[string]interface{}:
		processed := map[string]interface{}{}
		for key, val := range value {
			if key == "Ref" && parentKey != "" {
				valuStr := fmt.Sprintf("%v", val)
				if _, ok := mapOfReferences[valuStr]; ok {
					if parentVal, ok := parent.(map[string]interface{}); ok {
						parentVal[parentKey] = valuStr
						val = valuStr
						return val
					}
				}
			}
			processed[key] = a.resolveIndirectReferences(value, val, key, mapOfReferences)
		}
		return processed

	case []interface{}:

		// We found an array in the JSON - recurse through it's elements looking for intrinsic functions
		processed := []interface{}{}
		for _, val := range value {
			processed = append(processed, a.resolveIndirectReferences(parent, val, parentKey, mapOfReferences))
		}
		return processed

	case nil:
		return value
	case bool:
		return value
	case float64:
		return value
	case string:
		return value
	default:
		return nil

	}
}
