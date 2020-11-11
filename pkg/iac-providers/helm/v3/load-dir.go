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
	fileMap, err := utils.FindFilesBySuffix(absRootDir, []string{helmChartFilename})
	if err != nil {
		zap.S().Error("error while searching for helm charts", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(fileMap) == 0 {
		zap.S().Error("", zap.String("root dir", absRootDir), zap.Error(err))
		err = errNoHelmChartsFound
		return allResourcesConfig, err
	}

	// fileDir now contains the chart path
	iacDocumentMap := make(map[string][]*utils.IacDocument)
	for fileDir, chartFilename := range fileMap {
		chartPath := filepath.Join(fileDir, *chartFilename[0])
		zap.S().Debug("processing chart", zap.String("chart path", chartPath), zap.Error(err))

		var iacDocuments []*utils.IacDocument
		var chartMap map[string]interface{}
		iacDocuments, chartMap, err = h.loadChart(chartPath)
		if err != nil && err != errSkipTestDir {
			zap.S().Warn("error occurred while loading chart", zap.String("chart path", chartPath), zap.Error(err))
			continue
		}

		iacDocumentMap[chartPath] = iacDocuments

		var config *output.ResourceConfig
		config, err = h.createHelmChartResource(chartPath, chartMap)
		if err != nil {
			zap.S().Debug("failed to create helm chart resource", zap.Any("config", config))
			continue
		}

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}

	for _, iacDocuments := range iacDocumentMap {
		for _, doc := range iacDocuments {
			// @TODO add k8s version check
			var k k8sv1.K8sV1
			var config *output.ResourceConfig
			config, err = k.Normalize(doc)
			if err != nil {
				zap.S().Warn("unable to normalize data", zap.Error(err), zap.String("file", doc.FilePath))
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

	jsonData, err := json.Marshal(chartData)
	if err != nil {
		zap.S().Warn("unable to marshal chart to json", zap.String("chart path", chartPath))
		return nil, err
	}

	configData := make(map[string]interface{})
	if err = json.Unmarshal(jsonData, &configData); err != nil {
		zap.S().Warn("unable to unmarshal normalized config data", zap.String("chart path", chartPath))
		zap.S().Debug("failed config data", zap.Any("config", configData))
		return nil, err
	}

	chartName, ok := chartData["name"].(string)
	if !ok {
		zap.S().Warn("unable to determine chart name", zap.String("chart path", chartPath))
		return nil, err
	}

	config.Type = "helm_chart"
	config.Name = chartName
	config.Line = 0
	config.Source = chartPath
	config.ID = config.Type + "." + config.Name
	config.Config = configData

	return &config, nil
}

func (h *HelmV3) loadChart(chartPath string) ([]*utils.IacDocument, map[string]interface{}, error) {
	iacDocuments := make([]*utils.IacDocument, 0)
	chartMap := make(map[string]interface{})

	// load the chart file and values file from the specified chart path
	chartFileBytes, err := ioutil.ReadFile(chartPath)
	if err != nil {
		zap.S().Warn("unable to read", zap.String("file", chartPath))
		return iacDocuments, chartMap, err
	}

	if err = yaml.Unmarshal(chartFileBytes, &chartMap); err != nil {
		zap.S().Warn("unable to unmarshal values", zap.String("file", chartPath))
		return iacDocuments, chartMap, err
	}

	var fileInfo os.FileInfo
	chartDir := filepath.Dir(chartPath)
	valuesFile := filepath.Join(chartDir, helmValuesFilename)
	fileInfo, err = os.Stat(valuesFile)
	if err != nil {
		zap.S().Warn("unable to stat values.yaml", zap.String("chart path", chartPath))
		return iacDocuments, chartMap, err
	}

	var valueFileBytes []byte
	valueFileBytes, err = ioutil.ReadFile(valuesFile)
	if err != nil {
		zap.S().Warn("unable to read values.yaml", zap.String("file", fileInfo.Name()))
		return iacDocuments, chartMap, err
	}

	var valueMap map[string]interface{}
	if err = yaml.Unmarshal(valueFileBytes, &valueMap); err != nil {
		zap.S().Warn("unable to unmarshal values.yaml", zap.String("file", fileInfo.Name()))
		return iacDocuments, chartMap, err
	}

	// for each template file found, render and save an iacDocument
	var templateFileMap map[string][]*string
	templateFileMap, err = utils.FindFilesBySuffix(filepath.Join(chartDir, helmTemplateDir), h.getHelmTemplateExtensions())
	if err != nil {
		zap.S().Warn("error while calling FindFilesBySuffix", zap.String("filepath", fileInfo.Name()))
		return iacDocuments, chartMap, err
	}
	for templateDir, templateFiles := range templateFileMap {
		if filepath.Base(templateDir) == helmTestDir {
			zap.S().Debug("skipping test dir", zap.String("dir", templateDir))
			return iacDocuments, chartMap, errSkipTestDir
		}
		chartFiles := make([]*chart.File, 0)
		for _, templateFile := range templateFiles {
			var fileData []byte
			fileData, err = ioutil.ReadFile(filepath.Join(templateDir, *templateFile))
			if err != nil {
				zap.S().Warn("unable to read template file", zap.String("file", *templateFile))
				return iacDocuments, chartMap, err
			}

			chartFiles = append(chartFiles, &chart.File{
				Name: filepath.Join(helmTemplateDir, *templateFile),
				Data: fileData,
			})
		}

		chartName, ok := chartMap["name"].(string)
		if !ok {
			return iacDocuments, chartMap, errBadChartName
		}

		var chartVersion string
		chartVersion, ok = chartMap["version"].(string)
		if !ok {
			return iacDocuments, chartMap, errBadChartVersion
		}

		c := &chart.Chart{
			Metadata:  &chart.Metadata{Name: chartName, Version: chartVersion},
			Templates: chartFiles,
		}

		var v chartutil.Values
		v, err = chartutil.CoalesceValues(c, chartutil.Values{
			"Values": valueMap,
			"Release": chartutil.Values{
				"Name": defaultChartName,
			},
		})
		if err != nil {
			zap.S().Warn("error encountered in CoalesceValues", zap.String("chart path", chartPath))
			return iacDocuments, chartMap, err
		}

		var renderData map[string]string
		renderData, err = engine.Render(c, v)
		if err != nil {
			zap.S().Warn("error encountered while rendering chart", zap.String("chart path", chartPath),
				zap.String("template dir", templateDir))
			return iacDocuments, chartMap, err
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
	}

	return iacDocuments, chartMap, nil
}

func (h *HelmV3) getHelmTemplateExtensions() []string {
	return []string{"yaml", "tpl"}
}
