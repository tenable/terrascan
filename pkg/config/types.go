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

package config

// Global initializes GlobalConfig struct
var global *TerrascanConfig

// TerrascanConfig struct defines global variables/configurations across terrascan
type TerrascanConfig struct {
	Policy              `toml:"policy,omitempty" yaml:"policy,omitempty"`
	Notifications       map[string]Notifier `toml:"notifications,omitempty" yaml:"notifications,omitempty"`
	Rules               `toml:"rules,omitempty" yaml:"rules,omitempty"`
	Category            `toml:"category,omitempty" yaml:"category,omitempty"`
	Severity            `toml:"severity,omitempty" yaml:"severity,omitempty"`
	K8sAdmissionControl `toml:"k8s-admission-control,omitempty" yaml:"k8s-admission-control,omitempty"`
}

// Category defines the categories of violations that you want to be reported
type Category struct {
	List []string `toml:"list" yaml:"list"`
}

// Severity defines the minimum level of severity of violations that you want to be reported
type Severity struct {
	Level string `toml:"level" yaml:"level"`
}

// Policy struct defines policy specific configurations
type Policy struct {
	// policy repo local path
	BasePath string `toml:"path,omitempty" yaml:"path,omitempty"`
	// local filepath where repository containing policies is cached at
	RepoPath string `toml:"rego_subdir,omitempty" yaml:"rego_subdir,omitempty"`

	// policy git url and branch
	RepoURL string `toml:"repo_url,omitempty" yaml:"repo_url,omitempty"`
	Branch  string `toml:"branch,omitempty" yaml:"branch,omitempty"`

	// policy environment and access token
	Environment string `toml:"environment,omitempty" yaml:"environment,omitempty"`
	AccessToken string `toml:"access_token,omitempty" yaml:"access_token,omitempty"`
}

// Notifier represent a single notification in the terrascan config file
type Notifier struct {
	NotifierType   string      `toml:"type" yaml:"type"`
	NotifierConfig interface{} `toml:"config" yaml:"config"`
}

// Rules represents scan and skip rules in the terrascan config file
type Rules struct {
	ScanRules []string `toml:"scan-rules,omitempty" yaml:"scan-rules,omitempty"`
	SkipRules []string `toml:"skip-rules,omitempty" yaml:"skip-rules,omitempty"`
}

// K8sAdmissionControl deny rules in the terrascan config file
type K8sAdmissionControl struct {
	Dashboard      bool     `toml:"dashboard,omitempty" yaml:"dashboard,omitempty"`
	DeniedSeverity string   `toml:"denied-severity,omitempty" yaml:"denied-severity,omitempty"`
	Categories     []string `toml:"denied-categories,omitempty" yaml:"denied-categories,omitempty"`
	SaveRequests   bool     `toml:"save-requests,omitempty" yaml:"save-requests,omitempty"`
}
