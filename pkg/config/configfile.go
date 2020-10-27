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
	"os"

	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

var (
	errNotPresent           = fmt.Errorf("config file not present")
	errNotifierNotSupported = fmt.Errorf("notifier not supported")
	errTomlLoadConfig       = fmt.Errorf("failed to load toml config")
)

// NewNotifiers returns a list of notifiers configured in the config file
func LoadConfig(configFile string) (*toml.Tree, error) {

	// empty config file path
	if configFile == "" {
		zap.S().Debug("no config file specified")
		return nil, nil
	}

	// check if file exists
	_, err := os.Stat(configFile)
	if err != nil {
		zap.S().Errorf("config file '%s' not present", configFile)
		return nil, errNotPresent
	}

	// parse toml config file
	config, err := toml.LoadFile(configFile)
	if err != nil {
		zap.S().Errorf("failed to load toml config file '%s'. error: '%v'", err)
		return nil, errTomlLoadConfig
	}

	// return config Tree
	return config, nil
}
