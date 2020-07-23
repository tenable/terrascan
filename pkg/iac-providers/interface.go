package iacprovider

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// IacProvider defines the interface which every IaC provider needs to implement
// to claim support in terrascan
type IacProvider interface {
	LoadIacFile(string) (output.AllResourceConfigs, error)
	LoadIacDir(string) (output.AllResourceConfigs, error)
}
