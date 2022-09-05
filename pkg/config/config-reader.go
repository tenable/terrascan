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
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	tomlExtension  = ".toml"
	yamlExtension1 = ".yaml"
	yamlExtension2 = ".yml"
)

var (
	// ErrTomlLoadConfig indicates error: Failed to load toml config
	ErrTomlLoadConfig = fmt.Errorf("failed to load toml config")
	// ErrNotPresent indicates error: Config file not present
	ErrNotPresent = fmt.Errorf("config file not present")
)

// TerrascanConfigReader holds the terrascan config file name
type TerrascanConfigReader struct {
	config TerrascanConfig
}

// NewTerrascanConfigReader initialises and returns a config reader
func NewTerrascanConfigReader(fileName string) (*TerrascanConfigReader, error) {
	config := TerrascanConfig{}
	configReader := new(TerrascanConfigReader)
	configReader.config = config

	// empty file name check should be done by the caller, this is a safe check
	if fileName == "" {
		zap.S().Debug("no config file specified")
		return configReader, nil
	}

	// return error if file doesn't exist
	_, err := os.Stat(fileName)
	if err != nil {
		zap.S().Errorf("config file: %s, doesn't exist", fileName)
		return configReader, ErrNotPresent
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		zap.S().Error("error loading config file", zap.Error(err))
		return configReader, ErrTomlLoadConfig
	}

	// check the extension of the file and decode using the file contents
	// using the relevant package
	switch filepath.Ext(fileName) {
	case tomlExtension:
		err = toml.Unmarshal(data, &configReader.config)
	case yamlExtension1, yamlExtension2:
		err = yaml.Unmarshal(data, &configReader.config)
	default:
		err = fmt.Errorf("file format %q not support for terrascan config file",
			filepath.Ext(fileName))
	}
	if err != nil {
		return configReader, err
	}

	return configReader, nil
}

// GetPolicyConfig will return the policy config from the terrascan config file
func (r TerrascanConfigReader) getPolicyConfig() Policy {
	return r.config.Policy
}

// GetNotifications will return the notifiers specified in the terrascan config file
func (r TerrascanConfigReader) getNotifications() map[string]Notifier {
	return r.config.Notifications
}

// GetRules will return the rules specified in the terrascan config file
func (r TerrascanConfigReader) getRules() Rules {
	return r.config.Rules
}

// GetCategory will return the category specified in the terrascan config file
func (r TerrascanConfigReader) getCategory() Category {
	return r.config.Category
}

// GetSeverity will return the level of severity specified in the terrascan config file
func (r TerrascanConfigReader) getSeverity() Severity {
	return r.config.Severity
}

// GetK8sAdmissionControl will return the k8s deny rules specified in the terrascan config file
func (r TerrascanConfigReader) GetK8sAdmissionControl() K8sAdmissionControl {
	return r.config.K8sAdmissionControl
}
