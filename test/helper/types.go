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

package helper

import (
	"github.com/accurics/terrascan/pkg/results"
)

type violations []*results.Violation

func (v violations) Len() int {
	return len(v)
}

func (v violations) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v violations) Less(i, j int) bool {
	if v[i].File < v[j].File {
		return true
	}
	if v[i].File > v[j].File {
		return false
	}

	if v[i].ResourceType < v[j].ResourceType {
		return true
	}

	if v[i].ResourceType > v[j].ResourceType {
		return false
	}

	if v[i].RuleName < v[j].RuleName {
		return true
	}

	if v[i].RuleName > v[j].RuleName {
		return false
	}

	if v[i].ResourceName < v[j].ResourceName {
		return true
	}

	if v[i].ResourceName > v[j].ResourceName {
		return false
	}

	if v[i].LineNumber < v[j].LineNumber {
		return true
	}

	return v[i].LineNumber > v[j].LineNumber
}
