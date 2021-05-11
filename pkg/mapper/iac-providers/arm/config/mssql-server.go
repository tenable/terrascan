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
	arm_administratorLogin         = "administratorLogin"
	arm_administratorLoginPassword = "administratorLoginPassword"
	arm_minimumTLSVersion          = "minimalTlsVersion"
)

const (
	tf_administratorLogin         = "administrator_login"
	tf_administratorLoginPassword = "administrator_login_password"
	tf_minimumTLSVersion          = "minimum_tls_version"
)

// MSSQLServerConfig returns config for azurerm_mssql_server
func MSSQLServerConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tf_location:                   fn.LookUp(nil, params, r.Location).(string),
		tf_name:                       fn.LookUp(nil, params, r.Name).(string),
		tf_tags:                       r.Tags,
		tf_administratorLogin:         fn.LookUp(vars, params, convert.ToString(r.Properties, arm_administratorLogin)).(string),
		tf_administratorLoginPassword: fn.LookUp(vars, params, convert.ToString(r.Properties, arm_administratorLoginPassword)).(string),
		tf_minimumTLSVersion:          fn.LookUp(vars, params, convert.ToString(r.Properties, arm_minimumTLSVersion)).(string),
	}
}
