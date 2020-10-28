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
	"testing"
)

func testConfigEnv(configPath string, t *testing.T) {
	// Load the test data to compare against
	config, err := loadConfigFile(configPath)
	if err != nil {
		t.Error(err)
	}

	LoadGlobalConfig(configPath)

	if Global.Policy.BasePath != config.Policy.BasePath {
		t.Errorf("BasePath not overridden!  %v != %v", Global.Policy.BasePath, config.Policy.BasePath)
	}

	if Global.Policy.RepoPath != config.Policy.RepoPath {
		t.Errorf("RepoPath not overridden!  %v != %v", Global.Policy.RepoPath, config.Policy.RepoPath)
	}

	if Global.Policy.RepoURL != config.Policy.RepoURL {
		t.Errorf("RepoURL not overridden!  %v != %v", Global.Policy.RepoURL, config.Policy.RepoURL)
	}

	if Global.Policy.Branch != config.Policy.Branch {
		t.Errorf("Branch not overridden!  %v != %v", Global.Policy.Branch, config.Policy.Branch)
	}
}

func TestConfigFileLoadsCorrectly(t *testing.T) {
	testConfigEnv("./testdata/terrascan-config.toml", t)
}
