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

// Global initalizes GlobalConfig struct
var global *TerrascanConfig = &TerrascanConfig{}

// TerrascanConfig struct defines global variables/configurations across terrascan
type TerrascanConfig struct {
	Policy        `toml:"policy,omitempty"`
	Notifications map[string]Notifier `toml:"notifications,omitempty"`
	Rules         `toml:"rules,omitempty"`
	Category      `toml:"category,omitempty"`
	Severity      `toml:"severity,omitempty"`
	K8sDenyRules  `toml:"k8s-deny-rules,omitempty"`
}

// Category defines the categories of violations that you want to be reported
type Category struct {
	List []string `toml:"list"`
}

// Severity defines the minimum level of severity of violations that you want to be reported
type Severity struct {
	Level string `toml:"level"`
}

// Policy struct defines policy specific configurations
type Policy struct {
	// policy repo local path
	BasePath string `toml:"path,omitempty"`
	// local filepath where repository containing policies is cached at
	RepoPath string `toml:"rego_subdir,omitempty"`

	// policy git url and branch
	RepoURL string `toml:"repo_url,omitempty"`
	Branch  string `toml:"branch,omitempty"`
}

// Notifier represent a single notification in the terrascan config file
type Notifier struct {
	NotifierType   string      `toml:"type"`
	NotifierConfig interface{} `toml:"config"`
}

// Rules represents scan and skip rules in the terrascan config file
type Rules struct {
	ScanRules []string `toml:"scan-rules,omitempty"`
	SkipRules []string `toml:"skip-rules,omitempty"`
}

// K8s deny rules in the terrascan config file
type K8sDenyRules struct {
	DeniedSeverity string   `toml:"denied-severity,omitempty"`
	Categories     []string `toml:"denied-categories,omitempty"`
}
