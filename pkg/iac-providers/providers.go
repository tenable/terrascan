package iacprovider

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

var (
	errIacNotSupported = fmt.Errorf("iac not supported")
)

// NewIacProvider returns a new IacProvider
func NewIacProvider(iacType, iacVersion string) (iacProvider IacProvider, err error) {

	// get IacProvider from supportedIacProviders
	iacProviderObject, supported := supportedIacProviders[supportedIacType(iacType)][supportedIacVersion(iacVersion)]
	if !supported {
		zap.S().Errorf("IaC type:'%s', version: '%s' not supported", iacType, iacVersion)
		return iacProvider, errIacNotSupported
	}

	return reflect.New(iacProviderObject).Interface().(IacProvider), nil
}

// IsIacSupported returns true/false depending on whether the IaC
// provider is supported in terrascan or not
func IsIacSupported(iacType, iacVersion string) bool {
	if _, supported := supportedIacProviders[supportedIacType(iacType)]; !supported {
		return false
	}
	if _, supported := supportedIacProviders[supportedIacType(iacType)][supportedIacVersion(iacVersion)]; !supported {
		return false
	}
	return true
}
