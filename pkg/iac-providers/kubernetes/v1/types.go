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

// K8sV1 struct implements the IacProvider interface
type K8sV1 struct{}

const (
	// YAMLExtension yaml
	YAMLExtension = "yaml"

	// YAMLExtension2 yml
	YAMLExtension2 = "yml"

	// JSONExtension json
	JSONExtension = "json"

	// UnknownExtension unknown
	UnknownExtension = "unknown"

	kubernetesTypeName = "kubernetes"
	defaultNamespace   = "default"
)

// K8sFileExtensions returns the valid extensions for k8s (yaml, yml, json)
func K8sFileExtensions() []string {
	return []string{YAMLExtension, YAMLExtension2, JSONExtension}
}
