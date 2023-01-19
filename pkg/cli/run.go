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

package cli

import (
	"errors"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/tenable/terrascan/pkg/downloader"
	"github.com/tenable/terrascan/pkg/runtime"
	"github.com/tenable/terrascan/pkg/utils"
	"github.com/tenable/terrascan/pkg/writer"
	"go.uber.org/zap"
)

const (
	yamlOutputFormat = "yaml"
	jsonOutputFormat = "json"
)

// ScanOptions represents scan command and its optional flags
type ScanOptions struct {
	// Policy path directory
	policyPath []string

	// Cloud type (aws, azure, gcp, github)
	policyType []string

	// IaC type (terraform)
	iacType string

	// IaC version (for terraform:v14)
	iacVersion string

	// Path to a single IaC file
	iacFilePath string

	// Path to a directory containing one or more IaC files
	iacDirPath string

	// remoteType indicates the type of remote backend. Supported backends are
	// git s3, gcs, http.
	remoteType string

	// remoteURL points to the remote Iac repository on git, s3, gcs, http
	remoteURL string

	// configOnly will output resource config (should only be used for debugging purposes)
	configOnly bool

	// configWithError will output resource config and encountered errors
	configWithError bool

	// config file path
	configFile string

	// the output format for writing the results
	outputType string

	// UseColors indicates whether to use color output
	UseColors bool
	useColors string // used for flag processing

	// ScanRules is the array of rules to scan
	scanRules []string

	// SkipRules is the array of rules to skip while scanning
	skipRules []string

	// categories is the array categories of policy violations that should be reported
	categories []string

	// severity is the level of severity of policy violations that should be reported
	severity string

	// verbose indicates whether to display all fields in default human readable output
	verbose bool

	// showPassedRules indicates whether to display passed rules or not
	showPassedRules bool

	// nonRecursive enables recursive scan for the terraform iac provider
	nonRecursive bool

	// useTerraformCache provides ability to use terraform init local cache for modules rather than downloading them.
	useTerraformCache bool

	// FindVulnerabilities gives option to scan container images for vulnerabilities
	findVulnerabilities bool

	// notificationWebhookURL is the URL where terrascan will send the scan report and normalized config json
	notificationWebhookURL string

	// notificationWebhookToken is the auth token to call the notification webhook URL
	notificationWebhookToken string

	// repoURL lets us specify URL of the repository being scanned
	repoURL string

	// repoRef lets us specify the branch of the repository being scanned
	repoRef string

	// logOutputDir lets us specify the directory to write scan result and log files
	logOutputDir string

	// valuesFiles is the array of helm values file
	valuesFiles []string
}

// NewScanOptions returns a new pointer to ScanOptions
func NewScanOptions() *ScanOptions {
	return new(ScanOptions)
}

// Scan executes the terrascan scan command
func (s *ScanOptions) Scan() error {
	if err := s.Init(); err != nil {
		zap.S().Error("scan init failed", zap.Error(err))
		return err
	}

	if err := s.Run(); err != nil {
		zap.S().Error("scan run failed", zap.Error(err))
		return err
	}
	return nil
}

// Init initalises and validates ScanOptions
func (s *ScanOptions) Init() error {
	s.initColor()
	if err := s.validate(); err != nil {
		zap.S().Error("failed to start scan", zap.Error(err))
		return err
	}
	return nil
}

// validate config only for human readable output
// rest command options are validated by the executor
func (s ScanOptions) validate() error {
	// human readable output doesn't support --config-only and --config-with-error flag
	// if --config-only/--config-with-error flag is set, then exit with an error
	// asking the user to use yaml or json output format
	if (s.configOnly || s.configWithError) && !(strings.EqualFold(s.outputType, yamlOutputFormat) || strings.EqualFold(s.outputType, jsonOutputFormat)) {
		return errors.New("please use yaml or json output format when using --config-only or --config-with-error flags")
	}
	return nil
}

// initialises use colors options
func (s *ScanOptions) initColor() {
	switch strings.ToLower(s.useColors) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			s.UseColors = true
		} else {
			s.UseColors = false
		}

	case "true":
		fallthrough
	case "t":
		fallthrough
	case "y":
		fallthrough
	case "1":
		fallthrough
	case "force":
		s.UseColors = true

	default:
		s.UseColors = false
	}
}

// Run executes terrascan in CLI mode
func (s *ScanOptions) Run() error {

	// temp dir to download the remote repo
	tempDir := utils.GenerateTempDir()
	defer os.RemoveAll(tempDir)

	// download remote repository
	err := s.downloadRemoteRepository(tempDir)
	if err != nil {
		return err
	}

	// create a new runtime executor for processing IaC
	executor, err := runtime.NewExecutor(s.iacType, s.iacVersion, s.policyType, s.iacFilePath, s.iacDirPath,
		s.policyPath, s.scanRules, s.skipRules, s.categories, s.severity, s.nonRecursive, s.useTerraformCache,
		s.findVulnerabilities, s.notificationWebhookURL, s.notificationWebhookToken, s.repoURL, s.repoRef, s.valuesFiles,
	)
	if err != nil {
		return err
	}

	// executor output
	results, err := executor.Execute(s.configOnly, s.configWithError)
	if err != nil {
		return err
	}

	// set the ResourcePath to remoteURL if remote directory is scanned.
	if !s.configOnly && s.remoteURL != "" {
		results.Violations.ViolationStore.Summary.ResourcePath = s.remoteURL
	}

	// write results to console
	err = s.writeResults(results)
	if err != nil {
		zap.S().Error("failed to write results", zap.Error(err))
		return err
	}

	if !s.configOnly && flag.Lookup("test.v") == nil {
		os.RemoveAll(tempDir)
		exitCode := getExitCode(results)
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}
	return nil
}

func (s *ScanOptions) downloadRemoteRepository(tempDir string) error {
	d := downloader.NewDownloader()
	path, err := d.DownloadWithType(s.remoteType, s.remoteURL, tempDir)
	if path != "" {
		s.iacDirPath = path
	}
	if err == downloader.ErrEmptyURLType {
		// url and type empty, proceed with regular scanning
		zap.S().Debugf("remote url and type not configured, proceeding with regular scanning")
	} else if err != nil {
		// some error while downloading remote repository
		return err
	}
	return nil
}

func (s ScanOptions) writeResults(results runtime.Output) error {

	var writers []io.Writer

	outputWriter := NewOutputWriter(s.UseColors)
	writers = append(writers, outputWriter)

	fileWriter, closeFile := NewFileWriter(s.logOutputDir, s.outputType)
	if fileWriter != nil {
		writers = append(writers, fileWriter)
		defer closeFile()
	}

	if s.configOnly {
		return writer.Write(s.outputType, results.ResourceConfig, writers)
	}

	if s.configWithError {
		return writer.Write(s.outputType, results, writers)
	}

	// add verbose flag to the scan summary
	results.Violations.ViolationStore.Summary.ShowViolationDetails = s.verbose

	if !s.showPassedRules {
		results.Violations.ViolationStore.PassedRules = nil
	}

	return writer.Write(s.outputType, results.Violations, writers)
}

// getExitCode returns appropriate exit code for terrascan based on scan output
func getExitCode(o runtime.Output) int {
	if len(o.Violations.ViolationStore.DirScanErrors) > 0 {
		if o.Violations.ViolationStore.Summary.ViolatedPolicies > 0 {
			return 5
		}
		return 4
	} else if o.Violations.ViolationStore.Summary.ViolatedPolicies > 0 {
		return 3
	}
	return 0
}
