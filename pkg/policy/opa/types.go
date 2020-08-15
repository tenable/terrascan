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

package opa

import (
	"context"
	"time"

	"github.com/accurics/terrascan/pkg/policy"

	"github.com/open-policy-agent/opa/rego"
)

// RegoMetadata The rego metadata struct which is read and saved from disk
type RegoMetadata struct {
	Name         string                 `json:"name"`
	File         string                 `json:"file"`
	TemplateArgs map[string]interface{} `json:"template_args"`
	Severity     string                 `json:"severity"`
	Description  string                 `json:"description"`
	ReferenceID  string                 `json:"reference_id"`
	Category     string                 `json:"category"`
	Version      int                    `json:"version"`
}

// RegoData Stores all information needed to evaluate and report on a rego rule
type RegoData struct {
	Metadata      RegoMetadata
	RawRego       []byte
	PreparedQuery *rego.PreparedEvalQuery
}

// EngineStats Contains misc stats
type EngineStats struct {
	ruleCount         int
	regoFileCount     int
	metadataFileCount int
	metadataCount     int
	runTime           time.Duration
}

// Engine Implements the policy engine interface
type Engine struct {
	results     policy.EngineOutput
	context     context.Context
	regoFileMap map[string][]byte
	regoDataMap map[string]*RegoData
	stats       EngineStats
}
