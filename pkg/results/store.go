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

package results

// NewViolationStore returns a new violation store
func NewViolationStore() *ViolationStore {
	return &ViolationStore{
		Violations:        []*Violation{},
		SkippedViolations: []*Violation{},
	}
}

// AddResult Adds individual violations into the violation store
// when skip is true, violations are added to skipped violations
func (s *ViolationStore) AddResult(violation *Violation, isSkipped bool) {
	if isSkipped {
		s.SkippedViolations = append(s.SkippedViolations, violation)
	} else {
		s.Violations = append(s.Violations, violation)
	}
}

// GetResults Retrieves all violations from the violation store
// when skip is true, it returns only the skipped violations
func (s *ViolationStore) GetResults(isSkipped bool) []*Violation {
	if isSkipped {
		return s.SkippedViolations
	}
	return s.Violations
}
