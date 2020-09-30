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

package utils

import (
	"testing"
)

func TestGenRandomString(t *testing.T) {

	table := []struct {
		name string
		want int
	}{
		{
			name: "gen random string 1",
			want: 3,
		},
		{
			name: "gen random string 2",
			want: 6,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := GenRandomString(tt.want)
			if len(got) != tt.want {
				t.Errorf("got: '%v', want: '%v'", len(got), tt.want)
			}
		})
	}
}
