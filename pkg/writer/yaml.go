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

package writer

import (
	"io"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	yamlFormat supportedFormat = "yaml"
)

func init() {
	RegisterWriter(yamlFormat, YAMLWriter)
}

// YAMLWriter prints data in YAML format
func YAMLWriter(data interface{}, writer io.Writer) error {
	j, _ := yaml.Marshal(data)
	if _, err := writer.Write(j); err != nil {
		zap.S().Debugf("failed to write output error: '%v'", err)
	}

	if _, err := writer.Write([]byte{'\n'}); err != nil {
		zap.S().Debugf("failed to write newline error: '%v'", err)
	}
	return nil
}
