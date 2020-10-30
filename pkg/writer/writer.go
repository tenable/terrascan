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
	"fmt"
	"io"

	"go.uber.org/zap"
)

var (
	errNotSupported = fmt.Errorf("output format not supported")
)

// Write method writes in the given format using the respective writer func
func Write(format string, data interface{}, writer io.Writer) error {

	writerFunc, present := writerMap[supportedFormat(format)]
	if !present {
		zap.S().Errorf("output format '%s' not supported", format)
		return errNotSupported
	}

	return writerFunc(data, writer)
}
