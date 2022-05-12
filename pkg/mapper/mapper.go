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

package mapper

import (
	"github.com/tenable/terrascan/pkg/mapper/core"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft"
)

// NewMapper returns a mapper based on IaC provider.
func NewMapper(iacType string) core.Mapper {
	switch iacType {
	case "cft":
		return cft.Mapper()
	case "arm":
		return arm.Mapper()
	}
	return nil
}
