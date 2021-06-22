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
	armLogin    = "login"
	armSid      = "sid"
	armTenantID = "tenantId"
)
const (
	tfLogin    = "login"
	tfObjectID = "object_id"
)

// SQLActiveDirectoryAdministratorConfig returns config for azurerm_sql_active_directory_administrator
func SQLActiveDirectoryAdministratorConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tfLocation: fn.LookUp(nil, params, r.Location).(string),
		tfName:     fn.LookUp(nil, params, r.Name).(string),
		tfTags:     r.Tags,
		tfLogin:    fn.LookUp(vars, params, convert.ToString(r.Properties, armLogin)).(string),
		tfObjectID: fn.LookUp(vars, params, convert.ToString(r.Properties, armSid)).(string),
		tfTenantID: fn.LookUp(vars, params, convert.ToString(r.Properties, armTenantID)).(string),
	}
}
