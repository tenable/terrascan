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

// Engine Policy Engine interface
type Engine interface {
	//Init method to initialize engine with policy path, and a pre load filter
	Init(string, PreLoadFilter) error
	Configure() error
	Evaluate(EngineInput, PreScanFilter) (EngineOutput, error)
	GetResults() EngineOutput
	Release() error
}

// FilterSpecification defines a function that
// RegoMetadata filter specifications should implement
type FilterSpecification interface {
	IsSatisfied(r *RegoMetadata) bool
}

// PreLoadFilter defines functions, that a pre load filter should implement
type PreLoadFilter interface {
	IsAllowed(r *RegoMetadata) bool
	IsFiltered(r *RegoMetadata) bool
}

// PreScanFilter defines function, that a pre scan filter should implement
type PreScanFilter interface {
	Filter(rmap map[string]*RegoData, input EngineInput) map[string]*RegoData
}
