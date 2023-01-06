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

package utils

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMergeMaps(t *testing.T) {

	RegisterFailHandler(Fail)

	type Address struct {
		city    string
		country string
	}
	type Phone struct {
		mobile string
		fixed  string
	}

	table := []struct {
		name    string
		a       map[interface{}]interface{}
		b       map[interface{}]interface{}
		want    map[interface{}]interface{}
		wantErr error
	}{
		{
			name: "two maps",
			a: map[interface{}]interface{}{
				"name": "food",
				"items": struct {
					source string
					price  float64
				}{"chicken", 1.75},
			},
			b: map[interface{}]interface{}{
				"name": "grocery",
				"items": struct {
					source string
					price  float64
				}{"tea", 100.75},
				"required": true,
			},
			want: map[interface{}]interface{}{
				"name": "grocery",
				"items": struct {
					source string
					price  float64
				}{"tea", 100.75},
				"required": true,
			},
		},
		{
			name: "second map overrides first ",
			a: map[interface{}]interface{}{
				"name": "food",
				"items": struct {
					source string
					price  float64
					phone  Phone
				}{
					"chicken", 1.75, Phone{
						mobile: "99999",
						fixed:  "44444",
					},
				},
			},
			b: map[interface{}]interface{}{
				"name": "grocery",
				"items": struct {
					source  string
					price   float64
					address Address
				}{"tea", 100.75, Address{
					city:    "New York",
					country: "USA",
				},
				},
				"required": true,
			},

			want: map[interface{}]interface{}{
				"name": "grocery",
				"items": struct {
					source  string
					price   float64
					address Address
				}{"tea", 100.75,
					Address{
						city:    "New York",
						country: "USA",
					},
				},
				"required": true,
			},
		},
		{
			name: "b map overrides a with data type change",
			a: map[interface{}]interface{}{
				"name": "grocery",
				"items": struct {
					source  string
					price   float64
					address Address
				}{"tea", 100.75, Address{
					city:    "New York",
					country: "USA",
				},
				},
				"required": true,
			},
			b: map[interface{}]interface{}{
				"name": "food",
				"items": struct {
					source string
					price  float64
					phone  Phone
				}{
					"chicken", 1.75, Phone{
						mobile: "99999",
						fixed:  "44444",
					},
				},
			},

			want: map[interface{}]interface{}{
				"name": "food",
				"items": struct {
					source string
					price  float64
					phone  Phone
				}{"chicken", 1.75,
					Phone{
						mobile: "99999",
						fixed:  "44444",
					},
				},
				"required": true,
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMaps(tt.a, tt.b)
			isEqual := reflect.DeepEqual(result, tt.want)
			Expect(isEqual).To(Equal(true))
		})
	}
}
