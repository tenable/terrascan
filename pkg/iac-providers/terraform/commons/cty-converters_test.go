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
	"reflect"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

const ctyValCreationError = "error while creating ctyval"

func TestCtyToFloat(t *testing.T) {

	testFlVal := 10.23
	testFloatCtyVal, err := gocty.ToCtyValue(testFlVal, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testStrVal := "test"
	testStringCtyVal, err := gocty.ToCtyValue(testStrVal, cty.String)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "input is float type",
			args: args{
				ctyVal: testFloatCtyVal,
			},
			want: testFlVal,
		},
		{
			name: "input is not float type",
			args: args{
				ctyVal: testStringCtyVal,
			},
			want:    0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToFloat(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToFloat() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtyToStr(t *testing.T) {
	testStrVal := "abc"
	testStringCtyVal, err := gocty.ToCtyValue(testStrVal, cty.String)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testIntVal := 10
	testIntCtyVal, err := gocty.ToCtyValue(testIntVal, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "input is string type",
			args: args{
				ctyVal: testStringCtyVal,
			},
			want: testStrVal,
		},
		{
			name: "input is not string type",
			args: args{
				ctyVal: testIntCtyVal,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToStr(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToStr() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtyToInt(t *testing.T) {
	testStrVal := "xyz"
	testStringCtyVal, err := gocty.ToCtyValue(testStrVal, cty.String)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testIntVal := 2
	testIntCtyVal, err := gocty.ToCtyValue(testIntVal, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "input is int type",
			args: args{
				ctyVal: testIntCtyVal,
			},
			want: testIntVal,
		},
		{
			name: "input is not int type",
			args: args{
				ctyVal: testStringCtyVal,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToInt(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToInt() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtyToBool(t *testing.T) {
	testBoolVal := true
	testBoolCtyVal, err := gocty.ToCtyValue(testBoolVal, cty.Bool)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testIntVal := 20
	testIntCtyVal, err := gocty.ToCtyValue(testIntVal, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "input is bool type",
			args: args{
				ctyVal: testBoolCtyVal,
			},
			want: testBoolVal,
		},
		{
			name: "input is not bool type",
			args: args{
				ctyVal: testIntCtyVal,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToBool(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToBool() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtyToSlice(t *testing.T) {
	var testNilSlice []interface{}

	testNumSlice := []int{1, 2, 3}
	testNumSliceCtyVal, err := gocty.ToCtyValue(testNumSlice, cty.List(cty.Number))
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testStrSlice := []string{"one", "two", "three"}
	testStrSliceCtyVal, err := gocty.ToCtyValue(testStrSlice, cty.List(cty.String))
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testFloatSlice := []float64{0.11, 2.22, 0.03}
	testFloatSliceCtyVal, err := gocty.ToCtyValue(testFloatSlice, cty.List(cty.Number))
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testBoolSlice := []bool{true, false, true}
	testBoolSliceCtyVal, err := gocty.ToCtyValue(testBoolSlice, cty.List(cty.Bool))
	if err != nil {
		t.Error(ctyValCreationError)
	}

	tupleElem1 := 10
	tupleELem2 := 0.02
	tupleElem3 := "test"
	tupleElem4 := true

	tupleElem1CtyVal1, err := gocty.ToCtyValue(tupleElem1, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}
	tupleElem1CtyVal2, err := gocty.ToCtyValue(tupleELem2, cty.Number)
	if err != nil {
		t.Error(ctyValCreationError)
	}
	tupleElem1CtyVal3, err := gocty.ToCtyValue(tupleElem3, cty.String)
	if err != nil {
		t.Error(ctyValCreationError)
	}
	tupleElem1CtyVal4, err := gocty.ToCtyValue(tupleElem4, cty.Bool)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testTupleCtyVal := cty.TupleVal([]cty.Value{tupleElem1CtyVal1, tupleElem1CtyVal2, tupleElem1CtyVal3, tupleElem1CtyVal4})

	testBoolVar := false
	testBoolCtyVal, err := gocty.ToCtyValue(testBoolVar, cty.Bool)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testEmptySliceCtyVal := cty.ListValEmpty(cty.String)

	testStringSetCtyVal := cty.SetVal([]cty.Value{
		cty.StringVal("test1"),
		cty.StringVal("test2"),
		cty.StringVal("test3"),
		cty.StringVal("test4"),
	})

	testIntSetCtyVal := cty.SetVal([]cty.Value{
		cty.NumberIntVal(1),
		cty.NumberIntVal(2),
	})

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "input is number slice",
			args: args{
				ctyVal: testNumSliceCtyVal,
			},
			want: []interface{}{1, 2, 3},
		},
		{
			name: "input is string slice",
			args: args{
				ctyVal: testStrSliceCtyVal,
			},
			want: []interface{}{"one", "two", "three"},
		},
		{
			name: "input is float slice",
			args: args{
				ctyVal: testFloatSliceCtyVal,
			},
			want: []interface{}{0.11, 2.22, 0.03},
		},
		{
			name: "input is bool slice",
			args: args{
				ctyVal: testBoolSliceCtyVal,
			},
			want: []interface{}{true, false, true},
		},
		{
			name: "input is tuple with elems of all types",
			args: args{
				ctyVal: testTupleCtyVal,
			},
			want: []interface{}{10, 0.02, "test", true},
		},
		{
			name: "input is set with elems of string type",
			args: args{
				ctyVal: testStringSetCtyVal,
			},
			want: []interface{}{"test1", "test2", "test3", "test4"},
		},
		{
			name: "input is set with elems of int type",
			args: args{
				ctyVal: testIntSetCtyVal,
			},
			want: []interface{}{1, 2},
		},
		{
			name: "input is not of list or tuple type",
			args: args{
				ctyVal: testBoolCtyVal,
			},
			want:    testNilSlice,
			wantErr: true,
		},
		{
			name: "input is empty list",
			args: args{
				ctyVal: testEmptySliceCtyVal,
			},
			want: testNilSlice,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToSlice(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToSlice() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtyToMap(t *testing.T) {
	emptyMapCtyVal := cty.MapValEmpty(cty.String)
	emptyObjectCtyVal, err := gocty.ToCtyValue(map[string]interface{}{}, cty.EmptyObject)
	if err != nil {
		t.Error(ctyValCreationError)
	}

	testMap1 := map[string]cty.Value{
		"key1": cty.StringVal("value1"),
		"key2": cty.StringVal("value2"),
		"key3": cty.StringVal("value3"),
		"key4": cty.StringVal("value4"),
	}

	testMap2 := map[string]cty.Value{
		"one":   cty.NumberIntVal(1),
		"two":   cty.NumberIntVal(2),
		"three": cty.NumberIntVal(3),
	}

	testMapCtyVal1 := cty.MapVal(testMap1)
	testMapCtyVal2 := cty.MapVal(testMap2)

	testObjectCtyVal := cty.ObjectVal(map[string]cty.Value{
		"stringVal": cty.StringVal("someValue"),
		"intVar":    cty.NumberIntVal(100),
		"boolVar":   cty.BoolVal(false),
		"floatVar":  cty.NumberFloatVal(10.10),
	})

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "empty input map",
			args: args{
				ctyVal: emptyMapCtyVal,
			},
			wantErr: true,
		},
		{
			name: "empty input object",
			args: args{
				ctyVal: emptyObjectCtyVal,
			},
			wantErr: true,
		},
		{
			name: "valid input map, [string]string",
			args: args{
				ctyVal: testMapCtyVal1,
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			name: "valid input map, [string]int",
			args: args{
				ctyVal: testMapCtyVal2,
			},
			want: map[string]interface{}{
				"one":   1,
				"two":   2,
				"three": 3,
			},
		},
		{
			name: "valid input object var",
			args: args{
				ctyVal: testObjectCtyVal,
			},
			want: map[string]interface{}{
				"stringVal": "someValue",
				"intVar":    100,
				"boolVar":   false,
				"floatVar":  10.10,
			},
		},
		{
			name: "not map type",
			args: args{
				ctyVal: cty.BoolVal(true),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctyToMap(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ctyToMap() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ctyToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertPrimitiveType(t *testing.T) {
	intCtyVal := cty.NumberIntVal(-10)
	floatCtyVal := cty.NumberFloatVal(3.141)
	unsignedIntCtyVal := cty.NumberUIntVal(44)
	stringCtyVal := cty.StringVal("Peter Parker")
	boolCtyVal := cty.BoolVal(true)
	nonPrimitiveCtyVal := cty.ListVal([]cty.Value{cty.StringVal("test")})

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "signed int input",
			args: args{
				ctyVal: intCtyVal,
			},
			want: -10,
		},
		{
			name: "float input",
			args: args{
				ctyVal: floatCtyVal,
			},
			want: 3.141,
		},
		{
			name: "unsigned int input",
			args: args{
				ctyVal: unsignedIntCtyVal,
			},
			want: 44,
		},
		{
			name: "string input",
			args: args{
				ctyVal: stringCtyVal,
			},
			want: "Peter Parker",
		},
		{
			name: "boolean input",
			args: args{
				ctyVal: boolCtyVal,
			},
			want: true,
		},
		{
			name: "non primitive input",
			args: args{
				ctyVal: nonPrimitiveCtyVal,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertPrimitiveType(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertPrimitiveType() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPrimitiveType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertComplexType(t *testing.T) {

	nonComplexCtyVal := cty.NumberIntVal(5)
	complexListCtyVal := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"application": cty.StringVal("app1"),
			"port":        cty.NumberIntVal(9010),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"application": cty.StringVal("app2"),
			"port":        cty.NumberIntVal(8000),
		}),
	})
	complexMapCtyVal := cty.MapVal(map[string]cty.Value{
		"first": cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("Thor"),
			"ID":   cty.NumberIntVal(1),
		}),
		"second": cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("Spiderman"),
			"ID":   cty.NumberIntVal(2),
		}),
	})

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "non complex input",
			args: args{
				ctyVal: nonComplexCtyVal,
			},
			wantErr: true,
		},
		{
			name: "list of objects input",
			args: args{
				ctyVal: complexListCtyVal,
			},
			want: []interface{}{
				map[string]interface{}{
					"application": "app1",
					"port":        9010,
				},
				map[string]interface{}{
					"application": "app2",
					"port":        8000,
				},
			},
		},
		{
			name: "map of objects input",
			args: args{
				ctyVal: complexMapCtyVal,
			},
			want: map[string]interface{}{
				"first": map[string]interface{}{
					"name": "Thor",
					"ID":   1,
				},
				"second": map[string]interface{}{
					"name": "Spiderman",
					"ID":   2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertComplexType(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertComplexType() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertComplexType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertCtyToGoNative(t *testing.T) {
	primitiveCtyVal := cty.StringVal("primitive")
	complexCtyVal := cty.ListVal([]cty.Value{cty.BoolVal(false)})

	type args struct {
		ctyVal cty.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "primitive cty input",
			args: args{
				ctyVal: primitiveCtyVal,
			},
			want: "primitive",
		},
		{
			name: "complex cty input",
			args: args{
				ctyVal: complexCtyVal,
			},
			want: []interface{}{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertCtyToGoNative(tt.args.ctyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertCtyToGoNative() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCtyToGoNative() want = %v, want %v", got, tt.want)
			}
		})
	}
}
