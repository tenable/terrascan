package iacProvider

import (
	"fmt"
	"log"
	"reflect"
)

// NewIacProvider returns a new IacProvider
func NewIacProvider(iacType, iacVersion string) (iacProvider IacProvider, err error) {

	// get IacProvider from supportedIacProviders
	iacProviderObject, supported := supportedIacProviders[supportedIacType(iacType)][supportedIacVersion(iacVersion)]
	if !supported {
		errMsg := fmt.Sprintf("IaC type:'%s', version: '%s' not supported", iacType, iacVersion)
		log.Printf(errMsg)
		return iacProvider, fmt.Errorf("errMsg")
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
