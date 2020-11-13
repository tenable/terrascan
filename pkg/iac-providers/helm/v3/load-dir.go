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

package helmv3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	k8sv1 "github.com/accurics/terrascan/pkg/iac-providers/kubernetes/v1"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

var (
	errSkipTestDir       = fmt.Errorf("skipping test directory")
	errNoHelmChartsFound = fmt.Errorf("no helm charts found")
	errBadChartName      = fmt.Errorf("bad chart name in Chart.yaml")
	errBadChartVersion   = fmt.Errorf("bad chart version in Chart.yaml")
)

// LoadIacDir loads all helm charts under the specified directory
func (h *HelmV3) LoadIacDir(absRootDir string) (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	// find all Chart.yaml files within the specified directory structure
	fileMap, err := utils.FindFilesBySuffix(absRootDir, h.getHelmChartFilenames())
	if err != nil {
		zap.S().Error("error while searching for helm charts", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(fileMap) == 0 {
		err = errNoHelmChartsFound
		zap.S().Error("", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
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
		iacDocuments, chartMap, err = h.loadChart(chartPath)
		if err != nil && err != errSkipTestDir {
			logger.Error("error occurred while loading chart", zap.Error(err))
			continue
		}

		iacDocumentMap[chartPath] = iacDocuments

		// for each chart, add a normalized helm_chart resource
		var config *output.ResourceConfig
		config, err = h.createHelmChartResource(chartPath, chartMap)
		if err != nil {
			logger.Error("failed to create helm chart resource", zap.Any("config", config), zap.Error(err))
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
				}
				continue
			}

			config.Line = 1
			config.Source = doc.FilePath

			allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
		}
	}

	return allResourcesConfig, nil
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
		fileData, err := ioutil.ReadFile(filepath.Join(templateDir, *templateFile))
		if err != nil {
			logger.Error("unable to read template file", zap.String("file", *templateFile), zap.Error(err))
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
		logger.Error("chart name was invalid")
		return iacDocuments, errBadChartName
	}

	var chartVersion string
	chartVersion, ok = chartMap["version"].(string)
	if !ok {
		logger.Error("chart version was invalid")
		return iacDocuments, errBadChartVersion
	}

	// build the minimum helm chart data input
	c := &chart.Chart{
		Metadata:  &chart.Metadata{Name: chartName, Version: chartVersion},
		Templates: chartFiles,
	}

	// create rendered values
	options := chartutil.ReleaseOptions{
		Name:      defaultChartName,
		Namespace: chartName + "-namespace",
	}

	v, err := chartutil.ToRenderValues(c, valueMap, options, nil)
	if err != nil {
		logger.Error("value rendering failed", zap.Any("values", v), zap.Error(err))
		return iacDocuments, err
	}

	// render all files within the chart
	var renderData map[string]string
	var e engine.Engine

	// lint mode tells the rendering engine to be less strict when it comes to required variables
	e.LintMode = true
	renderData, err = e.Render(c, v)
	if err != nil {
		logger.Error("error encountered while rendering chart", zap.String("template dir", templateDir), zap.Error(err))
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
func (h *HelmV3) loadChart(chartPath string) ([]*utils.IacDocument, helmChartData, error) {
	iacDocuments := make([]*utils.IacDocument, 0)
	chartMap := make(helmChartData)
	logger := zap.S().With("chart path", chartPath)

	// load the chart file and values file from the specified chart path
	chartFileBytes, err := ioutil.ReadFile(chartPath)
	if err != nil {
		logger.Error("unable to read", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	if err = yaml.Unmarshal(chartFileBytes, &chartMap); err != nil {
		logger.Error("unable to unmarshal values", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	var fileInfo os.FileInfo
	chartDir := filepath.Dir(chartPath)
	valuesFile := filepath.Join(chartDir, helmValuesFilename)
	fileInfo, err = os.Stat(valuesFile)
	if err != nil {
		logger.Error("unable to stat values.yaml", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	logger.With("file name", fileInfo.Name())
	var valueFileBytes []byte
	valueFileBytes, err = ioutil.ReadFile(valuesFile)
	if err != nil {
		logger.Error("unable to read values.yaml", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	var valueMap map[string]interface{}
	if err = yaml.Unmarshal(valueFileBytes, &valueMap); err != nil {
		logger.Error("unable to unmarshal values.yaml", zap.Error(err))
		return iacDocuments, chartMap, err
	}

	// for each template file found, render and save an iacDocument
	var templateFileMap map[string][]*string
	templateFileMap, err = utils.FindFilesBySuffix(filepath.Join(chartDir, helmTemplateDir), h.getHelmTemplateExtensions())
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
