/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy "PatchScheduleConfig the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, s"PatchScheduleConfigtware
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
	armDayOfWeek       = "dayOfWeek"
	armStartHourUtc    = "startHourUtc"
	armScheduleEntries = "scheduleEntries"
)

const (
	tfDayOfWeek    = "day_of_week"
	tfStartHourUTC = "start_hour_utc"
)

// PatchScheduleConfig returns config for patch_schedule
func PatchScheduleConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	sch := convert.ToMap(r.Properties, armScheduleEntries)
	return map[string]interface{}{
		tfDayOfWeek:    fn.LookUp(nil, params, convert.ToString(sch, armDayOfWeek)).(string),
		tfStartHourUTC: fn.LookUp(nil, params, convert.ToString(sch, armStartHourUtc)).(float64),
		tfTags:         r.Tags,
	}
}
