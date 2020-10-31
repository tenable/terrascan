package policy

import (
	"encoding/xml"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
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

func EngineOutputFromViolationStore(store *results.ViolationStore) EngineOutput {
	return EngineOutput{
		xml.Name{},
		store,
	}
}

func (me EngineOutput) AsViolationStore() results.ViolationStore {
	if me.ViolationStore == nil {
		return results.ViolationStore{}
	}
	return results.ViolationStore{
		Violations: me.Violations,
		Count:      me.Count,
	}
}
