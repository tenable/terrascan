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

package k8sv1

import (
	"encoding/json"
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	yamltojson "github.com/ghodss/yaml"
	"github.com/iancoleman/strcase"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const terrascanSkip = "terrascanSkip"

var (
	errUnsupportedDoc = fmt.Errorf("unsupported document type")
	// ErrNoKind is returned when the "kind" key is not available (not a valid kubernetes resource)
	ErrNoKind = fmt.Errorf("kind does not exist")
)

// k8sMetadata is used to pull the name, namespace types and annotations for a given resource
type k8sMetadata struct {
	Name         string                 `yaml:"name" json:"name"`
	GenerateName string                 `yaml:"generateName,omitempty" json:"generateName,omitempty"`
	Namespace    string                 `yaml:"namespace" json:"namespace"`
	Annotations  map[string]interface{} `yaml:"annotations" json:"annotations"`
}

// NameOrGenerateName gets the metadata's Name member, or if Name is not set then GenerateName (for CRDs, for example)
func (m k8sMetadata) NameOrGenerateName() string {
	if len(m.Name) > 0 {
		return m.Name
	}
	return m.GenerateName
}

// k8sResource is a generic struct to handle all k8s resource types
type k8sResource struct {
	APIVersion string      `yaml:"apiVersion" json:"apiVersion"`
	Kind       string      `yaml:"kind" json:"kind"`
	Metadata   k8sMetadata `yaml:"metadata" json:"metadata"`
}

// extractResource takes the incoming document and extracts the resource using a go struct
// returns the resource data and raw json byte output ready for normalization
func (k *K8sV1) extractResource(doc *utils.IacDocument) (*k8sResource, *[]byte, error) {
	var resource k8sResource
	switch doc.Type {
	case utils.YAMLDoc:
		data, err := yamltojson.YAMLToJSON(doc.Data)
		if err != nil {
			return nil, nil, err
		}
		err = yaml.Unmarshal(data, &resource)
		if err != nil {
			return nil, nil, err
		}
		return &resource, &data, nil
	case utils.JSONDoc:
		err := json.Unmarshal(doc.Data, &resource)
		if err != nil {
			return nil, nil, err
		}
		return &resource, &doc.Data, nil
	default:
		return nil, nil, errUnsupportedDoc
	}
}

// getNormalizedName returns the normalized name
// this matches the terraform-defined resource type when applicable
func (k *K8sV1) getNormalizedName(kind string) string {
	var name string
	switch kind {
	case "DaemonSet":
		name = kubernetesTypeName + "_daemonset"
	default:
		name = kubernetesTypeName + "_" + strcase.ToSnake(kind)
	}
	return name
}

// Normalize takes the input document and normalizes it
func (k *K8sV1) Normalize(doc *utils.IacDocument) (*output.ResourceConfig, error) {

	resource, jsonData, err := k.extractResource(doc)
	if err != nil {
		return nil, err
	}

	var resourceConfig output.ResourceConfig

	resourceConfig.Type = k.getNormalizedName(resource.Kind)

	switch resource.Kind {
	case "":
		// error case
		return nil, ErrNoKind
	// non-namespaced resources
	case "ClusterRole":
		fallthrough
	case "Namespace":
		resourceConfig.ID = resourceConfig.Type + "." + resource.Metadata.NameOrGenerateName()
	default:
		// namespaced-resources
		namespace := resource.Metadata.Namespace
		if namespace == "" {
			namespace = "default"
		}

		resourceConfig.ID = resourceConfig.Type + "." + resource.Metadata.NameOrGenerateName() + "." + namespace
	}

	// read and update skip rules, if present
	skipRules := readSkipRulesFromAnnotations(resource.Metadata.Annotations, resourceConfig.ID)
	if skipRules != nil {
		resourceConfig.SkipRules = append(resourceConfig.SkipRules, skipRules...)
	}

	configData := make(map[string]interface{})
	if err = json.Unmarshal(*jsonData, &configData); err != nil {
		return nil, err
	}

	resourceConfig.Name = resource.Metadata.NameOrGenerateName()
	resourceConfig.Config = configData

	return &resourceConfig, nil
}

func readSkipRulesFromAnnotations(annotations map[string]interface{}, resourceID string) []string {

	var skipRulesFromAnnotations interface{}
	var ok bool
	if skipRulesFromAnnotations, ok = annotations[terrascanSkip]; !ok {
		zap.S().Debugf("%s not present for resource: %s", terrascanSkip, resourceID)
		return nil
	}

	skipRules := make([]string, 0)
	if rules, ok := skipRulesFromAnnotations.([]interface{}); ok {
		for _, rule := range rules {
			if value, ok := rule.(string); ok {
				skipRules = append(skipRules, value)
			} else {
				zap.S().Debugf("each rule in %s must be of string type", terrascanSkip)
			}
		}
	} else {
		zap.S().Debugf("%s must be an array of rules to skip", terrascanSkip)
	}

	return skipRules
}
