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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

var (
	// ErrTomlLoadConfig indicates error: Failed to load toml config
	errTomlLoadConfig = fmt.Errorf("failed to load toml config")
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
		zap.S().Debugf("config file: %s, doesn't exist", fileName)
		return configReader, ErrNotPresent
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		zap.S().Debugf("error loading config file", zap.Error(err))
		return configReader, errTomlLoadConfig
	}

	if err = toml.Unmarshal(data, &configReader.config); err != nil {
		return configReader, err
	}
	return configReader, nil
}

// GetPolicyConfig will return the policy config from the terrascan config file
func (r TerrascanConfigReader) GetPolicyConfig() Policy {
	return r.config.Policy
}

// GetNotifications will return the notifiers specified in the terrascan config file
func (r TerrascanConfigReader) GetNotifications() map[string]Notifier {
	return r.config.Notifications
}

// GetRules will return the rules specified in the terrascan config file
func (r TerrascanConfigReader) GetRules() Rules {
	return r.config.Rules
}
