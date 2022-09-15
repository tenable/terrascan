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

package rego

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/utils"
)

var metaDataReferenceIDPattern = regexp.MustCompile(fmt.Sprintf("(%s|%s)", utils.MetaDataReferenceIDRegex, utils.MetaDataIDRegex))
var metaDataIDPattern = regexp.MustCompile(utils.MetaDataIDRegex)

func TestPolicyValidation(t *testing.T) {

	// The root directory for reading policies
	path := "."

	dirList, err := utils.FindAllDirectories(path)
	if err != nil {
		t.Errorf("Error while walking the policy directories : %s", err)
		return
	}

	for _, dir := range dirList {
		t.Run("Policy Validation", func(t *testing.T) {
			Validate(dir, t)
		})
	}

}

func Validate(dir string, t *testing.T) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err)
		return
	}
	metadataFiles := utils.FilterFileInfoBySuffix(&dirEntries, []string{".json"})

	for j := range metadataFiles {
		filePath := filepath.Join(dir, *metadataFiles[j])
		var regoMetadata *policy.RegoMetadata
		regoMetadata, err = LoadRegoMetadata(filePath)
		if err != nil {
			t.Errorf("Error while reading %s : %s", filePath, err)
			continue
		}

		validateRequiredFields(regoMetadata, filePath, t)

		if !metaDataReferenceIDPattern.MatchString(regoMetadata.ReferenceID) {
			t.Errorf("%s invalid reference_id pattern", filePath)
		}

		if !metaDataIDPattern.MatchString(regoMetadata.ID) {
			t.Errorf("%s invalid id pattern", filePath)
		}

		if _, err := os.Stat(filepath.Join(dir, regoMetadata.File)); errors.Is(err, os.ErrNotExist) {
			t.Errorf("%s doesn't exist", filepath.Join(dir, regoMetadata.File))
		}

	}

}

func validateRequiredFields(regoMetadata *policy.RegoMetadata, filepath string, t *testing.T) {

	if regoMetadata.Name == "" {
		validationErrorLogger("name", filepath, t)
	}
	if regoMetadata.File == "" {
		validationErrorLogger("file", filepath, t)
	}
	if regoMetadata.PolicyType == "" {
		validationErrorLogger("policy_type", filepath, t)
	}
	if regoMetadata.ResourceType == "" {
		validationErrorLogger("resource_type", filepath, t)
	}
	if regoMetadata.Severity == "" {
		validationErrorLogger("severity", filepath, t)
	}
	if regoMetadata.Description == "" {
		validationErrorLogger("description", filepath, t)
	}
	if regoMetadata.Category == "" {
		validationErrorLogger("category", filepath, t)
	}
	if regoMetadata.Version == 0 {
		validationErrorLogger("version", filepath, t)
	}
	if regoMetadata.ID == "" {
		validationErrorLogger("id", filepath, t)
	}

}

func validationErrorLogger(field string, filepath string, t *testing.T) {
	t.Errorf("Required Field missing in %s : \"%s\"", filepath, field)
}
