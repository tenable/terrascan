package policy

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
)

// EngineInput Contains data used as input to the engine
type EngineInput struct {
	InputData *output.AllResourceConfigs
}

// EngineOutput Contains data output from the engine
type EngineOutput struct {
	*results.ViolationStore `json:"results" yaml:"results" xml:"results,attr"`
}
