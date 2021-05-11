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
	"strings"

	"github.com/accurics/terrascan/pkg/mapper/convert"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	arm_emails              = "emails"
	arm_phone               = "phone"
	arm_alertNotifications  = "alertNotifications"
	arm_notificationsByRole = "notificationsByRole"
	arm_state               = "state"
)

const (
	tf_email              = "email"
	tf_phone              = "phone"
	tf_alertNotifications = "alert_notifications"
	tf_alertsToAdmins     = "alerts_to_admins"
)

// SecurityCenterContactConfig returns config for azurerm_security_center_contact
func SecurityCenterContactConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location: fn.LookUp(nil, params, r.Location).(string),
		tf_name:     fn.LookUp(nil, params, r.Name).(string),
		tf_tags:     r.Tags,
		tf_phone:    fn.LookUp(nil, params, convert.ToString(r.Properties, arm_phone)).(string),
		tf_email:    fn.LookUp(nil, params, convert.ToString(r.Properties, arm_emails)).(string),
	}

	notifications := convert.ToMap(r.Properties, arm_alertNotifications)
	state := convert.ToString(notifications, arm_state)
	cf[tf_alertNotifications] = strings.EqualFold(strings.ToUpper(state), "ON")

	notifications = convert.ToMap(r.Properties, arm_notificationsByRole)
	state = convert.ToString(notifications, state)
	cf[tf_alertsToAdmins] = strings.EqualFold(strings.ToUpper(state), "ON")

	return cf
}
