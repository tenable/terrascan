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

package policy

import (
	"reflect"
	"sort"
	"testing"
)

func TestPolicyTypeAllExpandedCorrectly(t *testing.T) {
	t.Run("policy type all gets right policy names", func(t *testing.T) {

		want := SupportedPolicyTypes(false)
		got := supportedCloudProvider["all"].policyNames()

		sort.Strings(want)
		sort.Strings(got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})

	t.Run("policy type all gets right policy paths", func(t *testing.T) {

		want := GetDefaultPolicyPaths(SupportedPolicyTypes(false))
		got := GetDefaultPolicyPaths([]string{"all"})

		sort.Strings(want)
		sort.Strings(got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}
