package iacProvider

import (
	"reflect"

	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
)

// SupportedIacType data type for supported IaC provider
type supportedIacType string

// supported IaC providers
const (
	terraform supportedIacType = "terraform"
)

// supportedIacVersion data type for supported Iac provider
type supportedIacVersion string

// supported Iac versions
const (
	defaultVersion supportedIacVersion = "default"
	terraformV12   supportedIacVersion = "v12"
)

// map of supported IaC providers
var supportedIacProviders map[supportedIacType](map[supportedIacVersion]reflect.Type)

// initializes a map of supported IaC providers
func init() {
	supportedIacProviders = make(map[supportedIacType](map[supportedIacVersion]reflect.Type))

	// terraform support
	supportedTerraformVersions := make(map[supportedIacVersion]reflect.Type)
	supportedTerraformVersions[terraformV12] = reflect.TypeOf(tfv12.TfV12{})
	supportedIacProviders[terraform] = supportedTerraformVersions
}
