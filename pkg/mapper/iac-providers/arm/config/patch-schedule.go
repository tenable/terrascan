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
	arm_dayOfWeek       = "dayOfWeek"
	arm_startHourUtc    = "startHourUtc"
	arm_scheduleEntries = "scheduleEntries"
)

const (
	tf_dayOfWeek    = "day_of_week"
	tf_startHourUTC = "start_hour_utc"
)

// RedisCacheConfig returns config for patch_schedule
func PatchScheduleConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	sch := convert.ToMap(r.Properties, arm_scheduleEntries)
	return map[string]interface{}{
		tf_dayOfWeek:    fn.LookUp(nil, params, convert.ToString(sch, arm_dayOfWeek)).(string),
		tf_startHourUTC: fn.LookUp(nil, params, convert.ToString(sch, arm_startHourUtc)).(float64),
		tf_tags:         r.Tags,
	}
}
