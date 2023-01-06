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

package helmv3

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	k8sv1 "github.com/tenable/terrascan/pkg/iac-providers/kubernetes/v1"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

var (
	errSkipTestDir     = fmt.Errorf("skipping test directory")
	errBadChartName    = fmt.Errorf("invalid chart name in Chart.yaml")
	errBadChartVersion = fmt.Errorf("invalid chart version in Chart.yaml")
)

const valuesFiles = "valuesFiles"

// LoadIacDir loads all helm charts under the specified directory
func (h *HelmV3) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	// find all Chart.yaml files within the specified directory structure
	fileMap, err := utils.FindFilesBySuffix(absRootDir, h.getHelmChartFilenames())
	if err != nil {
		zap.S().Debug("error while searching for helm charts", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(h.errIacLoadDirs, results.DirScanErr{IacType: "helm", Directory: absRootDir, ErrMessage: err.Error()})
	}

	if len(fileMap) == 0 {
		errMsg := fmt.Sprintf("no helm charts found in directory %s", absRootDir)
		return allResourcesConfig, multierror.Append(h.errIacLoadDirs, results.DirScanErr{IacType: "helm", Directory: absRootDir, ErrMessage: errMsg})
	}

	// fileDir now contains the chart path
	iacDocumentMap := make(map[string][]*utils.IacDocument)
	for fileDir, chartFilename := range fileMap {
		chartPath := filepath.Join(fileDir, *chartFilename[0])
		logger := zap.S().With("chart path", chartPath)
		logger.Debug("processing chart", zap.Error(err))

		// load helm charts into a map of IaC documents
		var iacDocuments []*utils.IacDocument
		var chartMap helmChartData
		iacDocuments, chartMap, err = h.loadChart(chartPath, options)
		if err != nil && err != errSkipTestDir {
			errMsg := fmt.Sprintf("error occurred while loading chart. err: %v", err)
			logger.Debug("error occurred while loading chart", zap.Error(err))
			h.errIacLoadDirs = multierror.Append(h.errIacLoadDirs, results.DirScanErr{IacType: "helm", Directory: fileDir, ErrMessage: errMsg})
			continue
		}

		iacDocumentMap[chartPath] = iacDocuments

		// for each chart, add a normalized helm_chart resource
		var config *output.ResourceConfig
		config, err = h.createHelmChartResource(chartPath, chartMap)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create helm chart resource. err: %v", err)
			logger.Error("failed to create helm chart resource", zap.Any("config", config), zap.Error(err))
			h.errIacLoadDirs = multierror.Append(h.errIacLoadDirs, results.DirScanErr{IacType: "helm", Directory: fileDir, ErrMessage: errMsg})
			continue
		}

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}

	// normalize all rendered IaC documents using the kubernetes code
	for _, iacDocuments := range iacDocumentMap {
		for _, doc := range iacDocuments {
			// helmv3 supports the kubernetes v1 api
			var k k8sv1.K8sV1
			var config *output.ResourceConfig
			config, err = k.Normalize(doc)
			if err != nil {
				// ignore logging errors when the "kind" field is not available because helm chart rendering can create an empty file
				// in that case, we should not output an error as it was the user's intention to prevent rendering the resource
				if err != k8sv1.ErrNoKind {
					zap.S().Error("unable to normalize data", zap.Error(err), zap.String("file", doc.FilePath))
					h.errIacLoadDirs = multierror.Append(h.errIacLoadDirs, err)
				}
				continue
			}

			config.Line = 1
			config.Source = doc.FilePath

			allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
		}
	}

	return allResourcesConfig, h.errIacLoadDirs
}

// createHelmChartResource returns normalized Helm Chart resource data
func (h *HelmV3) createHelmChartResource(chartPath string, chartData map[string]interface{}) (*output.ResourceConfig, error) {
	var config output.ResourceConfig

	logger := zap.S().With("helm chart path", chartPath)

	jsonData, err := json.Marshal(chartData)
	if err != nil {
		logger.Error("unable to marshal chart to json", zap.Error(err))
		return nil, err
	}

	configData := make(map[string]interface{})
	if err = json.Unmarshal(jsonData, &configData); err != nil {
		logger.Error("unable to unmarshal normalized config data", zap.Error(err))
		logger.Debug("failed config data", zap.Any("config", configData))
		return nil, err
	}

	chartName, ok := chartData["name"].(string)
	if !ok {
		logger.Error("unable to determine chart name", zap.Error(err))
		return nil, err
	}

	config.Type = "helm_chart"
	config.Name = chartName
	config.Line = 1
	config.Source = chartPath
	config.ID = config.Type + "." + config.Name
	config.Config = configData

	return &config, nil
}

// renderChart renders a helm chart with the given template files and values
// returns and IaC document for each rendered file
func (h *HelmV3) renderChart(chartPath string, chartMap helmChartData, templateDir string, templateFiles []*string, valueMap map[string]interface{}) ([]*utils.IacDocument, error) {
	const defaultNamespaceName = "default"

	iacDocuments := make([]*utils.IacDocument, 0)
	logger := zap.S().With("helm chart path", chartPath)

	if filepath.Base(templateDir) == helmTestDir {
		logger.Debug("skipping test dir", zap.String("dir", templateDir))
		return iacDocuments, errSkipTestDir
	}

	// create a list containing raw template file data
	chartFiles := make([]*chart.File, 0)
	for _, templateFile := range templateFiles {
		var fileData []byte
		fileData, err := os.ReadFile(filepath.Join(templateDir, *templateFile))
		if err != nil {
			logger.Debug("error while reading template file", zap.String("file", *templateFile), zap.Error(err))
			return iacDocuments, err
		}

		chartFiles = append(chartFiles, &chart.File{
			Name: filepath.Join(helmTemplateDir, *templateFile),
			Data: fileData,
		})
	}

	// chart name and version are required parameters
	chartName, ok := chartMap["name"].(string)
	if !ok {
		logger.Debug("chart name is invalid")
		return iacDocuments, errBadChartName
	}

	var chartVersion string
	chartVersion, ok = chartMap["version"].(string)
	if !ok {
		logger.Debug("chart version is invalid")
		return iacDocuments, errBadChartVersion
	}

	// build the minimum helm chart data input
	c := &chart.Chart{
		Metadata:  &chart.Metadata{Name: chartName, Version: chartVersion},
		Templates: chartFiles,
	}

	// create rendered values
	options := chartutil.ReleaseOptions{
		Name:      chartName,
		Namespace: defaultNamespaceName,
	}

	v, err := chartutil.ToRenderValues(c, valueMap, options, nil)
	if err != nil {
		logger.Debug("value rendering failed", zap.Any("values", v), zap.Error(err))
		return iacDocuments, err
	}

	// render all files within the chart
	var renderData map[string]string
	var e engine.Engine

	// lint mode tells the rendering engine to be less strict when it comes to required variables
	e.LintMode = true
	renderData, err = e.Render(c, v)
	if err != nil {
		logger.Debug("error encountered while rendering chart", zap.String("template dir", templateDir), zap.Error(err))
		return iacDocuments, err
	}

	for renderFile := range renderData {
		iacDocuments = append(iacDocuments, &utils.IacDocument{
			Data:      []byte(renderData[renderFile]),
			Type:      utils.YAMLDoc,
			StartLine: 1,
			EndLine:   1,
			FilePath:  renderFile,
		})
	}

	return iacDocuments, nil
}

// loadChart renders and loads all templates within a chart path
func (h *HelmV3) loadChart(chartPath string, options map[string]interface{}) ([]*utils.IacDocument, helmChartData, error) {
	iacDocuments := make([]*utils.IacDocument, 0)
	chartMap := make(helmChartData)
	logger := zap.S().With("chart path", chartPath)

	chartDir := filepath.Dir(chartPath)
	// load the chart file and values file from the specified chart path
	chartFileBytes, err := os.ReadFile(chartPath)
	if err != nil {
		logger.Debug("unable to read", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	if err = yaml.Unmarshal(chartFileBytes, &chartMap); err != nil {
		logger.Debug("unable to unmarshal values", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	valuesFilePaths := []string{}
	// check if custom values files are given as options
	if valuesFiles, ok := options[valuesFiles]; ok {
		valuesFilePaths = valuesFiles.([]string)
		logger.Debug("found user defined values.yaml files list", zap.Any("values", valuesFilePaths))
	}
	if len(valuesFilePaths) == 0 { // if no values-files list provided then use default values.yaml file
		valuesFilePaths = []string{helmValuesFilename}
		logger.Debug("defaulting to values.yaml file present in the current directory", zap.Any("values", valuesFilePaths))
	}

	allValuesFiles := make([]map[interface{}]interface{}, 0)

	for _, valuesFile := range valuesFilePaths {
		valuesFilePath := filepath.Join(chartDir, valuesFile)
		valuesMap, err := h.readFileIntoInterface(valuesFilePath)
		if err != nil {
			return iacDocuments, chartMap, err
		}
		allValuesFiles = append(allValuesFiles, valuesMap)
	}

	if len(allValuesFiles) > 0 {
		resultValueMap := allValuesFiles[0]
		for i := 1; i < len(valuesFilePaths); i++ {
			resultValueMap = utils.MergeMaps(resultValueMap, allValuesFiles[i])
		}

		outValuesBytes, err := yaml.Marshal(resultValueMap)
		if err != nil {
			logger.Debug("unable to marshal merged values.yaml", zap.Error(err))
			return iacDocuments, chartMap, err
		}
		var valueMap map[string]interface{}
		// UnMarshal back to map[string]interface{}
		if err = yaml.Unmarshal(outValuesBytes, &valueMap); err != nil {
			logger.Debug("unable to unmarshal values.yaml", zap.Error(err))
			return iacDocuments, chartMap, err
		}
		iacDocuments, chartMap, err = h.getIACDocumentsWithValues(chartPath, chartMap, valueMap)
		if err != nil {
			logger.Warn("error rendering chart with merged values file", zap.Error(err))
			return iacDocuments, chartMap, err
		}
	}
	return iacDocuments, chartMap, nil
}

// getIACDocumentsWithValues returns iacDocument given chart path and values map
func (h *HelmV3) getIACDocumentsWithValues(chartPath string, chartMap helmChartData, valueMap map[string]interface{}) ([]*utils.IacDocument, helmChartData, error) {
	iacDocuments := make([]*utils.IacDocument, 0)
	logger := zap.S().With("chart path", chartPath)

	chartDir := filepath.Dir(chartPath)
	// for each template file found, render and save an iacDocument
	var templateFileMap map[string][]*string
	templateFileMap, err := utils.FindFilesBySuffix(filepath.Join(chartDir, helmTemplateDir), h.getHelmTemplateExtensions())
	if err != nil {
		logger.Warn("error while calling FindFilesBySuffix", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	var renderedCharts []*utils.IacDocument
	for templateDir, templateFiles := range templateFileMap {
		renderedCharts, err = h.renderChart(chartPath, chartMap, templateDir, templateFiles, valueMap)
		if err != nil {
			continue
		}
		iacDocuments = append(iacDocuments, renderedCharts...)
	}
	return iacDocuments, chartMap, err
}

// getHelmTemplateExtensions returns valid helm template extensions
func (h *HelmV3) getHelmTemplateExtensions() []string {
	return []string{"yaml", "yml", "tpl"}
}

// getHelmChartFilenames returns valid chart filenames
func (h *HelmV3) getHelmChartFilenames() []string {
	return []string{"Chart.yaml"}
}

// Name returns name of the provider
func (h *HelmV3) Name() string {
	return "helm"
}

// readFileIntoInterface reads and converts file into interface
func (h *HelmV3) readFileIntoInterface(valuesFile string) (map[interface{}]interface{}, error) {
	logger := zap.S().With("readFileIntoInterface", valuesFile)
	fileInfo, err := os.Stat(valuesFile)
	if err != nil {
		logger.Debug("unable to stat values.yaml", zap.Error(err))
		return nil, err
	}

	logger.With("file name", fileInfo.Name())
	var valueFileBytes []byte
	valueFileBytes, err = os.ReadFile(valuesFile)
	if err != nil {
		logger.Debug("unable to read values.yaml", zap.Error(err))
		return nil, err
	}

	var valueMap map[interface{}]interface{}
	if err = yaml.Unmarshal(valueFileBytes, &valueMap); err != nil {
		logger.Debug("unable to unmarshal values.yaml", zap.Error(err))
		return nil, err
	}

	return valueMap, nil
}
