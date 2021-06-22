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
	armEmails              = "emails"
	armPhone               = "phone"
	armAlertNotifications  = "alertNotifications"
	armNotificationsByRole = "notificationsByRole"
	armState               = "state"
)

const (
	tfEmail              = "email"
	tfPhone              = "phone"
	tfAlertNotifications = "alert_notifications"
	tfAlertsToAdmins     = "alerts_to_admins"
)

// SecurityCenterContactConfig returns config for azurerm_security_center_contact
func SecurityCenterContactConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUp(nil, params, r.Location).(string),
		tfName:     fn.LookUp(nil, params, r.Name).(string),
		tfTags:     r.Tags,
		tfPhone:    fn.LookUp(nil, params, convert.ToString(r.Properties, armPhone)).(string),
		tfEmail:    fn.LookUp(nil, params, convert.ToString(r.Properties, armEmails)).(string),
	}

	notifications := convert.ToMap(r.Properties, armAlertNotifications)
	state := convert.ToString(notifications, armState)
	cf[tfAlertNotifications] = strings.EqualFold(strings.ToUpper(state), "ON")

	notifications = convert.ToMap(r.Properties, armNotificationsByRole)
	state = convert.ToString(notifications, state)
	cf[tfAlertsToAdmins] = strings.EqualFold(strings.ToUpper(state), "ON")

	return cf
}
