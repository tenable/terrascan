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
	"encoding/xml"
	"io"

	"go.uber.org/zap"
)

const (
	xmlFormat supportedFormat = "xml"
)

func init() {
	RegisterWriter(xmlFormat, XMLWriter)
}

// XMLWriter prints data in XML format
func XMLWriter(data interface{}, writer io.Writer) error {
	j, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		zap.S().Errorf("failed to write XML output. error: '%v'", err)
		return err
	}
	writer.Write(j)
	writer.Write([]byte{'\n'})
	return nil
}
