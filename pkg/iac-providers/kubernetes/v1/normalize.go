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
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	k8sappsv1 "k8s.io/api/apps/v1"
	k8sbatchv1 "k8s.io/api/batch/v1"
	k8sbatchv1beta1 "k8s.io/api/batch/v1beta1"
	k8scorev1 "k8s.io/api/core/v1"
)

const (
	terrascanMaxSeverity = "runterrascan.io/maxseverity"
	terrascanMinSeverity = "runterrascan.io/minseverity"
)

var (
	errUnsupportedDoc = fmt.Errorf("unsupported document type")
	// ErrNoKind is returned when the "kind" key is not available (not a valid kubernetes resource)
	ErrNoKind = fmt.Errorf("kind does not exist")

	infileInstructionNotPresentLog = "%s not present for resource: %s"
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

func (k *K8sV1) extractContainerImages(kind string, doc *utils.IacDocument) ([]string, []string, error) {
	var containerImages = make([]string, 0)
	var initContainerImages = make([]string, 0)
	var data []byte
	var err error

	if doc.Data == nil {
		return containerImages, initContainerImages, errors.Errorf("document does not have any resource data for unmarshalling")
	}

	if doc.Type == utils.YAMLDoc {
		data, err = yamltojson.YAMLToJSON(doc.Data)
		if err != nil {
			return nil, nil, err
		}
	} else {
		data = doc.Data
	}

	switch kind {
	case "Pod":
		pod := k8scorev1.Pod{}
		err = json.Unmarshal(data, &pod)
		if err != nil {
			err := errors.Errorf("error unmarshalling pod: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range pod.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range pod.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "Deployment":
		deployment := k8sappsv1.Deployment{}
		err = json.Unmarshal(data, &deployment)
		if err != nil {
			err := errors.Errorf("error unmarshalling deployment: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range deployment.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range deployment.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "ReplicationController":
		rc := k8scorev1.ReplicationController{}
		err = json.Unmarshal(data, &rc)
		if err != nil {
			err := errors.Errorf("error unmarshalling replicationcontroller: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range rc.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range rc.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "Job":
		job := k8sbatchv1.Job{}
		err = json.Unmarshal(data, &job)
		if err != nil {
			err := errors.Errorf("error unmarshalling job: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range job.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range job.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "CronJob":
		cronjob := k8sbatchv1beta1.CronJob{}
		err = json.Unmarshal(data, &cronjob)
		if err != nil {
			err := errors.Errorf("error unmarshalling cronjob: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "StatefulSet":
		ss := k8sappsv1.StatefulSet{}
		err = json.Unmarshal(data, &ss)
		if err != nil {
			err := errors.Errorf("error unmarshalling statefulset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range ss.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range ss.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "ReplicaSet":
		rs := k8sappsv1.ReplicaSet{}
		err = json.Unmarshal(data, &rs)
		if err != nil {
			err := errors.Errorf("error unmarshalling replicaset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range rs.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range rs.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	case "DaemonSet":
		ds := k8sappsv1.DaemonSet{}
		err = json.Unmarshal(data, &ds)
		if err != nil {
			err := errors.Errorf("error unmarshalling daemonset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		for _, container := range ds.Spec.Template.Spec.Containers {
			containerImages = append(containerImages, container.Image)
		}
		for _, container := range ds.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, container.Image)
		}
	default:
		zap.S().Debugf("the container image extraction for kubernetes workload of kind %s is not supported.", kind)
	}
	return containerImages, initContainerImages, nil
}

// Normalize takes the input document and normalizes it
func (k *K8sV1) Normalize(doc *utils.IacDocument) (*output.ResourceConfig, error) {

	resource, jsonData, err := k.extractResource(doc)
	if err != nil {
		return nil, err
	}

	var resourceConfig output.ResourceConfig
	resourceConfig.K8sWorkloadContainerImages = make([]string, 0)
	resourceConfig.K8sWorkloadInitContainerImages = make([]string, 0)
	var containerImages, initContainerImages []string
	resourceConfig.Type = k.getNormalizedName(resource.Kind)

	switch resource.Kind {
	case "":
		// error case
		return nil, ErrNoKind
	// non-namespaced resources
	case "ClusterRole":
		fallthrough
	// pod and all kinds of workloads
	case "Pod", "Deployment", "ReplicaSet", "ReplicationController", "Job", "CronJob", "StatefulSet", "DaemonSet":
		containerImages, initContainerImages, err = k.extractContainerImages(resource.Kind, doc)
		if err != nil {
			return nil, err
		}
		fallthrough
	default:
		// namespaced-resources
		namespace := resource.Metadata.Namespace
		if namespace == "" {
			namespace = "default"
		}

		resourceConfig.ID = resourceConfig.Type + "." + resource.Metadata.NameOrGenerateName() + "-" + namespace
	}

	resourceConfig.K8sWorkloadContainerImages = append(resourceConfig.K8sWorkloadContainerImages, containerImages...)
	resourceConfig.K8sWorkloadInitContainerImages = append(resourceConfig.K8sWorkloadInitContainerImages, initContainerImages...)

	// read and update skip rules, if present
	skipRules := utils.ReadSkipRulesFromMap(resource.Metadata.Annotations, resourceConfig.ID)
	if skipRules != nil {
		resourceConfig.SkipRules = append(resourceConfig.SkipRules, skipRules...)
	}

	maxSeverity, minSeverity := readMinMaxSeverityFromAnnotations(resource.Metadata.Annotations, resourceConfig.ID)

	resourceConfig.MaxSeverity = maxSeverity
	resourceConfig.MinSeverity = minSeverity

	configData := make(map[string]interface{})
	if err = json.Unmarshal(*jsonData, &configData); err != nil {
		return nil, err
	}

	resourceConfig.Name = resource.Metadata.NameOrGenerateName()
	resourceConfig.Config = configData

	return &resourceConfig, nil
}

// readMinMaxSeverityFromAnnotations finds the min max severity values set in annotations for the resource
func readMinMaxSeverityFromAnnotations(annotations map[string]interface{}, resourceID string) (maxSeverity, minSeverity string) {
	var (
		minSeverityAnnotation interface{}
		maxSeverityAnnotation interface{}
		ok                    bool
	)
	if minSeverityAnnotation, ok = annotations[terrascanMinSeverity]; !ok {
		zap.S().Debugf(infileInstructionNotPresentLog, terrascanMinSeverity, resourceID)
	} else if minSeverity, ok = minSeverityAnnotation.(string); !ok {
		zap.S().Debugf("%s must be a string cantaining value as (High | Low| Medium)", terrascanMinSeverity)
	}
	if maxSeverityAnnotation, ok = annotations[terrascanMaxSeverity]; !ok {
		zap.S().Debugf(infileInstructionNotPresentLog, terrascanMaxSeverity, resourceID)
	} else if maxSeverity, ok = maxSeverityAnnotation.(string); !ok {
		zap.S().Debugf("%s must be a string cantaining value as (High | Low| Medium)", terrascanMaxSeverity)
	}
	return
}
