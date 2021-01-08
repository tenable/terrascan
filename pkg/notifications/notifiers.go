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
	"reflect"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

const (
	notificationsConfigKey = "notifications"
)

var (
	errNotifierNotSupported   = fmt.Errorf("notifier not supported")
	errNotifierTypeNotPresent = fmt.Errorf("notifier type not present in toml config")
	// ErrTomlKeyNotPresent will be returned when config file does not have notificationsConfigKey
	ErrTomlKeyNotPresent = fmt.Errorf("key not present in toml config")
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

	if configFile == "" {
		zap.S().Debug("no config file specified")
		return notifiers, nil
	}

	configReader, err := config.NewTerrascanConfigReader(configFile)
	if err != nil || configReader == nil {
		return notifiers, err
	}

	// get config for 'notifications'
	notifications := configReader.GetNotifications()
	if notifications == nil {
		zap.S().Debug("key '%s' not present in toml config", notificationsConfigKey)
		return notifiers, ErrTomlKeyNotPresent
	}

	// create notifiers
	var allErrs error
	for _, notifier := range notifications {
		if notifier.NotifierType == "" {
			zap.S().Error(errNotifierTypeNotPresent)
			allErrs = utils.WrapError(errNotifierTypeNotPresent, allErrs)
			continue
		}

		if !IsNotifierSupported(notifier.NotifierType) {
			zap.S().Errorf("notifier type '%s' not supported", notifier.NotifierType)
			allErrs = utils.WrapError(errNotifierNotSupported, allErrs)
			continue
		}

		// create a new notifier
		n, err := NewNotifier(notifier.NotifierType)
		if err != nil {
			allErrs = utils.WrapError(err, allErrs)
			continue
		}

		// populate data
		err = n.Init(notifier.NotifierConfig)
		if err != nil {
			allErrs = utils.WrapError(err, allErrs)
			continue
		}

		// add to the list of notifiers
		notifiers = append(notifiers, n)
	}

	// return list of notifiers
	return notifiers, allErrs
}

// IsNotifierSupported returns true/false depending on whether the notifier
// is supported in terrascan or not
func IsNotifierSupported(notifierType string) bool {
	if _, supported := supportedNotifiers[supportedNotifierType(notifierType)]; !supported {
		return false
	}
	return true
}
