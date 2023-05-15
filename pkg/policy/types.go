package policy

import (
	"encoding/xml"

	"github.com/open-policy-agent/opa/rego"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/results"
)

const (
	version12 = "v12"
	version1  = "v1"
	version3  = "v3"
	version4  = "v4"
)

// EngineInput Contains data used as input to the engine
type EngineInput struct {
	InputData *output.AllResourceConfigs
}

// EngineOutput Contains data output from the engine
type EngineOutput struct {
	XMLName                 xml.Name `json:"-" yaml:"-" xml:"results"`
	*results.ViolationStore `json:"results" yaml:"results" xml:"results"`
}

// EngineOutputFromViolationStore returns an EngineOutput initialized from ViolationStore
func EngineOutputFromViolationStore(store *results.ViolationStore) EngineOutput {
	return EngineOutput{
		xml.Name{},
		store,
	}
}

// AsViolationStore returns EngineOutput as a ViolationStore
func (me EngineOutput) AsViolationStore() results.ViolationStore {
	if me.ViolationStore == nil {
		return results.ViolationStore{}
	}
	return results.ViolationStore{
		Violations:        me.Violations,
		Vulnerabilities:   me.Vulnerabilities,
		SkippedViolations: me.SkippedViolations,
		PassedRules:       me.PassedRules,
		Summary:           me.Summary,
	}
}

// RegoMetadata The rego metadata struct which is read and saved from disk
type RegoMetadata struct {
	Name         string                 `json:"name"`
	File         string                 `json:"file"`
	PolicyType   string                 `json:"policy_type"`
	ResourceType string                 `json:"resource_type"`
	TemplateArgs map[string]interface{} `json:"template_args"`
	Severity     string                 `json:"severity"`
	Description  string                 `json:"description"`
	ReferenceID  string                 `json:"reference_id"`
	Category     string                 `json:"category"`
	Version      int                    `json:"version"`
	ID           string                 `json:"id"`
}

// RegoData Stores all information needed to evaluate and report on a rego rule
type RegoData struct {
	Metadata      RegoMetadata
	RawRego       []byte
	PreparedQuery *rego.PreparedEvalQuery
}
