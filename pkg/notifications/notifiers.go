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

package notifications

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

const (
	notificationsConfigKey = "notifications"
)

var (
	errNotPresent           = fmt.Errorf("config file not present")
	errNotifierNotSupported = fmt.Errorf("notifier not supported")
	errTomlLoadConfig       = fmt.Errorf("failed to load toml config")
	errTomlKeyNotPresent    = fmt.Errorf("key not present in toml config")
)

// NewNotifier returns a new notifier
func NewNotifier(notifierType string) (notifier Notifier, err error) {

	// get notifier from supportedNotifierss
	notifierObject, supported := supportedNotifiers[supportedNotifierType(notifierType)]
	if !supported {
		zap.S().Errorf("notifier type '%s' not supported", notifierType)
		return notifier, errNotifierNotSupported
	}

	// successful
	return reflect.New(notifierObject).Interface().(Notifier), nil
}

// NewNotifiers returns a list of notifiers configured in the config file
func NewNotifiers(configFile string) ([]Notifier, error) {

	var notifiers []Notifier

	// empty config file path
	if configFile == "" {
		zap.S().Infof("no config file specified")
		return notifiers, nil
	}

	// check if file exists
	_, err := os.Stat(configFile)
	if err != nil {
		zap.S().Errorf("config file '%s' not present", configFile)
		return notifiers, errNotPresent
	}

	// parse toml config file
	config, err := toml.LoadFile(configFile)
	if err != nil {
		zap.S().Errorf("failed to load toml config file '%s'. error: '%v'", err)
		return notifiers, errTomlLoadConfig
	}

	// get config for 'notifications'
	keyConfig := config.Get(notificationsConfigKey)
	if keyConfig == nil {
		zap.S().Infof("key '%s' not present in toml config", notificationsConfigKey)
		return notifiers, errTomlKeyNotPresent
	}

	// get all the notifier types configured in TOML config
	keyTomlConfig := keyConfig.(*toml.Tree)
	notifierTypes := keyTomlConfig.Keys()

	// create notifiers
	for _, nType := range notifierTypes {

		// check if toml config present for notifier type
		nTypeConfig := keyTomlConfig.Get(nType)
		if nTypeConfig == nil {
			zap.S().Errorf("notifier '%v' config not present", nType)
			return notifiers, errTomlKeyNotPresent
		}

		// create a new notifier
		n, err := NewNotifier(nType)
		if err != nil {
			continue
		}

		// populate data
		err = n.Init(nTypeConfig)
		if err != nil {
			continue
		}

		// add to the list of notifiers
		notifiers = append(notifiers, n)
	}

	// return list of notifiers
	return notifiers, nil
}

// IsNotifierSupported returns true/false depending on whether the notifier
// is supported in terrascan or not
func IsNotifierSupported(notifierType string) bool {
	if _, supported := supportedNotifiers[supportedNotifierType(notifierType)]; !supported {
		return false
	}
	return true
}
