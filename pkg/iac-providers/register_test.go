package iacprovider

import (
	"reflect"
	"testing"
)

type MockIacProvider struct{}

func TestRegisterIacProvider(t *testing.T) {

	/*
		table := []struct {
			name       string
			iacType    supportedIacType
			iacVersion supportedIacVersion
			want       reflect.Type
		}{
			{
				name:       "mock iac type and version",
				iacType:    supportedIacType("mockIacType"),
				iacVersion: supportedIacVersion("mockIacVersion"),
				want:       reflect.TypeOf(MockIacProvider{}),
			},
			{
				name:       "mock iac type default version",
				iacType:    supportedIacType("mockIacType"),
				iacVersion: supportedIacVersion(""),
				want:       reflect.TypeOf(MockIacProvider{}),
			},
		}
	*/

	t.Run("mock iac provider", func(t *testing.T) {

		var (
			iacType           = supportedIacType("mockIacType")
			iacVersion        = supportedIacVersion("mockIacVersion")
			defaultIacVersion = iacVersion
			want              = reflect.TypeOf(MockIacProvider{})
		)

		RegisterIacProvider(iacType, iacVersion, defaultIacVersion, want)

		if _, present := supportedIacProviders[iacType]; !present {
			t.Errorf("mockIacType not registered")
		}
		got, present := supportedIacProviders[iacType][iacVersion]
		if !present {
			t.Errorf("mockIacVersion not registered")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})

	t.Run("mock iac default version", func(t *testing.T) {

		var (
			iacType           = supportedIacType("mockIacType")
			iacVersion        = supportedIacVersion("mockIacVersion")
			defaultIacVersion = iacVersion
			want              = reflect.TypeOf(MockIacProvider{})
		)

		RegisterIacProvider(iacType, iacVersion, defaultIacVersion, want)

		if _, present := supportedIacProviders[iacType]; !present {
			t.Errorf("mockIacType not registered")
		}
		got, present := supportedIacProviders[iacType][defaultIacVersion]
		if !present {
			t.Errorf("defaultIacVersion not registered")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}
