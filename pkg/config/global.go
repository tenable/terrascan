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

import (
	"path/filepath"
	"strings"

	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

const (
	defaultPolicyRepoURL     = "https://github.com/tenable/terrascan.git"
	defaultPolicyBranch      = "master"
	defaultPolicyEnvironment = "https://cloud.tenable.com"
)

// ConfigEnvvarName env variable
const ConfigEnvvarName = "TERRASCAN_CONFIG"

var (
	defaultPolicyRepoPath = filepath.Join("pkg", "policies", "opa", "rego")
	defaultBasePolicyPath = filepath.Join(utils.GetHomeDir(), ".terrascan")
)

// LoadGlobalConfig loads policy configuration from specified configFile
// into var Global.Policy.  Members of Global.Policy that are not specified
// in configFile will get default values
func LoadGlobalConfig(configFile string) error {
	// Start with the defaults
	global = &TerrascanConfig{}

	global.Policy = Policy{
		BasePath: defaultBasePolicyPath,
		RepoPath: defaultPolicyRepoPath,
		RepoURL:  defaultPolicyRepoURL,
		Branch:   defaultPolicyBranch,
	}

	var configReader *TerrascanConfigReader
	var err error

	if configReader, err = NewTerrascanConfigReader(configFile); err != nil {
		return err
	}

	if configFile != "" {
		zap.S().Debugf("loading global config from: %s", configFile)
	}

	if len(configReader.getPolicyConfig().BasePath) > 0 && len(configReader.getPolicyConfig().RepoPath) == 0 {
		zap.S().Warnf("policy base path specified in configfile '%s', but rego_subdir path not specified. applying default rego_subdir value '%s'", configFile, GetPolicyRepoPath())
	}

	if len(configReader.getPolicyConfig().RepoPath) > 0 && len(configReader.getPolicyConfig().BasePath) == 0 {
		zap.S().Warnf("policy rego_subdir specified in configfile '%s', but base path not specified. applying default base path value '%s'", configFile, GetPolicyBasePath())
	}

	if len(configReader.getPolicyConfig().BasePath) > 0 {
		global.BasePath = configReader.getPolicyConfig().BasePath
	}

	if len(configReader.getPolicyConfig().RepoPath) > 0 {
		global.RepoPath = configReader.getPolicyConfig().RepoPath
	}

	absolutePolicyBasePath, absolutePolicyRepoPath, err := utils.GetAbsPolicyConfigPaths(GetPolicyBasePath(), GetPolicyRepoPath())
	if err != nil {
		zap.S().Error("error processing provided policy paths", zap.Error(err))
		return err
	}

	global.Policy.BasePath = absolutePolicyBasePath
	global.Policy.RepoPath = absolutePolicyRepoPath

	if len(configReader.getPolicyConfig().RepoURL) > 0 {
		global.Policy.RepoURL = configReader.getPolicyConfig().RepoURL
	}
	if len(configReader.getPolicyConfig().Branch) > 0 {
		global.Policy.Branch = configReader.getPolicyConfig().Branch
	}
	if len(configReader.getPolicyConfig().Environment) > 0 {
		global.Policy.Environment = configReader.getPolicyConfig().Environment
	}
	if len(configReader.getPolicyConfig().AccessToken) > 0 {
		global.Policy.AccessToken = configReader.getPolicyConfig().AccessToken
	}

	if len(configReader.getRules().ScanRules) > 0 {
		global.Rules.ScanRules = configReader.getRules().ScanRules
	}

	if len(configReader.getRules().SkipRules) > 0 {
		global.Rules.SkipRules = configReader.getRules().SkipRules
	}

	if len(configReader.getSeverity().Level) > 0 {
		global.Severity.Level = configReader.getSeverity().Level
	}

	if len(configReader.getNotifications()) > 0 {
		global.Notifications = configReader.getNotifications()
	}

	if len(configReader.getCategory().List) > 0 {
		global.Category.List = configReader.getCategory().List
	}

	global.K8sAdmissionControl = configReader.GetK8sAdmissionControl()

	zap.S().Debugf("global config loaded")

	return nil
}

// GetPolicyBasePath returns the configured policy base path
func GetPolicyBasePath() string {
	if global == nil {
		return defaultBasePolicyPath
	}
	return global.Policy.BasePath
}

// GetPolicyRepoPath return the configured path to the policies repo locally downloaded
func GetPolicyRepoPath() string {
	if global == nil {
		return defaultPolicyRepoPath
	}
	return global.Policy.RepoPath
}

// GetPolicyRepoURL returns the configured policy repo url
func GetPolicyRepoURL() string {
	if global == nil {
		return defaultPolicyRepoURL
	}
	return global.Policy.RepoURL
}

// GetPolicyBranch returns the configured policy repo url
func GetPolicyBranch() string {
	if global == nil {
		return defaultPolicyBranch
	}
	return global.Policy.Branch
}

// GetPolicyEnvironment returns the configured policy environment url
func GetPolicyEnvironment() string {
	if global == nil {
		return defaultPolicyEnvironment
	}
	return strings.TrimRight(global.Policy.Environment, "/")
}

// GetPolicyAccessToken returns the configured policy access token
func GetPolicyAccessToken() string {
	if global == nil {
		return ""
	}
	return global.Policy.AccessToken
}

// GetScanRules returns the configured scan rules
func GetScanRules() []string {
	if global == nil {
		return nil
	}
	return global.Rules.ScanRules
}

// GetSkipRules returns the configured skips rules
func GetSkipRules() []string {
	if global == nil {
		return nil
	}
	return global.Rules.SkipRules
}

// GetSeverityLevel returns the configured severity level
func GetSeverityLevel() string {
	if global == nil {
		return ""
	}
	return global.Severity.Level
}

// GetCategoryList returns the configured list of category of violations
func GetCategoryList() []string {
	if global == nil {
		return nil
	}
	return global.Category.List
}

// GetNotifications returns the configured notifier map
func GetNotifications() map[string]Notifier {
	if global == nil {
		return nil
	}
	return global.Notifications
}

// GetK8sAdmissionControl returns kubernetes admission control configuration
func GetK8sAdmissionControl() K8sAdmissionControl {
	if global == nil {
		return K8sAdmissionControl{}
	}
	return global.K8sAdmissionControl
}
