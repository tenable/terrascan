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

package utils

import (
	"github.com/pkg/errors"
)

// WrapError wraps given err with allErrs and returns a unified error
func WrapError(err, allErrs error) error {
	// if allErrs is empty, return err
	if allErrs == nil {
		return err
	}

	// if err empty return allErrs
	if err == nil {
		return allErrs
	}

	// wrap err with allErrs
	allErrs = errors.Wrap(err, allErrs.Error())
	return allErrs
}
