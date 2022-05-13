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

package policy

import (
	"reflect"
	"sort"
	"testing"
)

func TestSupportedPolicyTypes(t *testing.T) {
	t.Run("supported policy types", func(t *testing.T) {
		var want []string
		for k := range supportedCloudProvider {
			want = append(want, string(k))
		}
		sort.Strings(want)
		got := SupportedPolicyTypes(true)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}

func TestSupportedNotIndirectPolicyTypes(t *testing.T) {
	t.Run("supported policy types", func(t *testing.T) {
		var want []string
		for k, v := range supportedCloudProvider {
			if !v.isIndirect {
				want = append(want, string(k))
			}
		}
		sort.Strings(want)
		got := SupportedPolicyTypes(false)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}
