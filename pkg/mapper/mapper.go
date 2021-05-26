package mapper

import (
	"github.com/accurics/terrascan/pkg/mapper/core"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/cft"
)

// NewMapper returns a mapper based on IaC provider.
func NewMapper(iacType string) core.Mapper {
	switch iacType {
	case "cft":
		return cft.Mapper()
	}
	return nil
}
