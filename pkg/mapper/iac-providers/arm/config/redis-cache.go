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
	"github.com/accurics/terrascan/pkg/mapper/convert"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	arm_sku              = "sku"
	arm_family           = "family"
	arm_capacity         = "capacity"
	arm_enableNonSSLPort = "enableNonSslPort"
)

const (
	tf_enableNonSSLPort = "enable_non_ssl_port"
	tf_capacity         = "capacity"
	tf_family           = "family"
)

// RedisCacheConfig returns config for azurerm_redis_cache
func RedisCacheConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:         fn.LookUp(nil, params, r.Location).(string),
		tf_name:             fn.LookUp(nil, params, r.Name).(string),
		tf_tags:             r.Tags,
		tf_enableNonSSLPort: fn.LookUp(nil, params, convert.ToString(r.Properties, arm_enableNonSSLPort)).(bool),
	}

	s := convert.ToMap(r.Properties, arm_sku)
	cf[tf_skuName] = fn.LookUp(nil, params, convert.ToString(s, tf_name)).(string)
	cf[tf_family] = fn.LookUp(nil, params, convert.ToString(s, arm_family)).(string)
	cf[tf_capacity] = fn.LookUp(nil, params, convert.ToString(s, arm_capacity)).(float64)

	return cf
}
