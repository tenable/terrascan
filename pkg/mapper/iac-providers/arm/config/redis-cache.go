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

package config

import (
	"github.com/tenable/terrascan/pkg/mapper/convert"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	armSku              = "sku"
	armFamily           = "family"
	armCapacity         = "capacity"
	armEnableNonSSLPort = "enableNonSslPort"
)

const (
	tfEnableNonSSLPort = "enable_non_ssl_port"
	tfCapacity         = "capacity"
	tfFamily           = "family"
)

// RedisCacheConfig returns config for azurerm_redis_cache
func RedisCacheConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(nil, params, r.Location),
		tfName:     fn.LookUpString(nil, params, r.Name),
		tfTags:     r.Tags,
	}

	if enabledNonSSLPort, ok := fn.LookUp(nil, params, convert.ToString(r.Properties, armEnableNonSSLPort)).(bool); ok {
		cf[tfEnableNonSSLPort] = enabledNonSSLPort
	}

	s := convert.ToMap(r.Properties, armSku)
	cf[tfSkuName] = fn.LookUpString(nil, params, convert.ToString(s, tfName))
	cf[tfFamily] = fn.LookUpString(nil, params, convert.ToString(s, armFamily))
	cf[tfCapacity] = fn.LookUpFloat64(nil, params, convert.ToString(s, armCapacity))

	return cf
}
