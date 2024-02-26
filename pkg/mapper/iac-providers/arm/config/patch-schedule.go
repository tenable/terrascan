/*
    Copyright (C) 2022 Tenable, Inc.

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
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
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
		tfDayOfWeek:    fn.LookUpString(nil, params, convert.ToString(sch, armDayOfWeek)),
		tfStartHourUTC: fn.LookUpFloat64(nil, params, convert.ToString(sch, armStartHourUtc)),
		tfTags:         functions.PatchAWSTags(r.Tags),
	}
}
