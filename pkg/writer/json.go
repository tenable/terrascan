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
	"encoding/json"
	"io"
)

const (
	jsonFormat supportedFormat = "json"
)

func init() {
	RegisterWriter(jsonFormat, JSONWriter)
}

// JSONWriter prints data in JSON format
func JSONWriter(data interface{}, writer io.Writer) error {
	j, _ := json.MarshalIndent(data, "", "  ")
	writer.Write(j)
	writer.Write([]byte{'\n'})
	return nil
}
