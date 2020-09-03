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
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

var (
	errBadResourceType      = fmt.Errorf("bad resource type")
	errKeyDoesNotExist      = fmt.Errorf("key does not exist")
	errMetadataDoesNotExist = fmt.Errorf("metadata does not exist")
	errMetadataNameField    = fmt.Errorf("unable to parse the metadata name field")
	errInvalidNamespaceType = fmt.Errorf("invalid namespace type")
)

func (k *K8sV1) normalize(doc *utils.IacDocument) (*output.ResourceConfig, error) {

	// if the document is yaml, convert it to json first
	var data *map[string]interface{}
	if doc.Type == utils.YAMLDoc {
		var err error
		data, err = utils.YAMLtoJSON(doc.Data)
		if err != nil {
			return nil, err
		}
	}

	// resource type
	_, ok := (*data)["kind"]
	if !ok {
		return nil, errBadResourceType
	}

	var resourceType string
	resourceType, ok = (*data)["kind"].(string)
	if !ok {
		return nil, errBadResourceType
	}

	metadataVal, ok := (*data)["metadata"]
	if !ok {
		return nil, errKeyDoesNotExist
	}

	var metadata map[string]interface{}
	metadata, ok = metadataVal.(map[string]interface{})
	if !ok {
		return nil, errMetadataDoesNotExist
	}

	namespace := "default"
	var resourceName string
	if resourceType == "Namespace" || resourceType == "ClusterRole" {
		resourceName, ok = metadata["name"].(string)
		if !ok {
			return nil, errMetadataNameField
		}
	} else {
		// sets the namespace
		// if no namespace is specified, the default namespace is used
		var namespaceVal interface{}
		if namespaceVal, ok = metadata["namespace"]; ok {
			// set the namespace if available, otherwise use the default
			namespace, _ = namespaceVal.(string)
		}

		// extract the resource name and set the resource id
		resourceName, ok = metadata["name"].(string)
		if !ok {
			return nil, errInvalidNamespaceType
		}
	}

	return &output.ResourceConfig{
		Type:   kubernetesTypeName + resourceType,
		ID:     kubernetesTypeName + resourceType + "." + resourceName + "." + namespace,
		Name:   resourceName,
		Config: data,
	}, nil
}
