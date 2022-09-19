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

package k8sv1

import (
	"encoding/json"

	yamltojson "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
	k8sappsv1 "k8s.io/api/apps/v1"
	k8sbatchv1 "k8s.io/api/batch/v1"
	k8sbatchv1beta1 "k8s.io/api/batch/v1beta1"
	k8scorev1 "k8s.io/api/core/v1"
)

func (k *K8sV1) extractContainerImages(kind string, doc *utils.IacDocument) ([]output.ContainerDetails, []output.ContainerDetails, error) {
	var containerImages = make([]output.ContainerDetails, 0)
	var initContainerImages = make([]output.ContainerDetails, 0)
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
		containerImages = append(containerImages, readContainers(pod.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(pod.Spec.InitContainers)...)

	case "Deployment":
		deployment := k8sappsv1.Deployment{}
		err = json.Unmarshal(data, &deployment)
		if err != nil {
			err := errors.Errorf("error unmarshalling deployment: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(deployment.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(deployment.Spec.Template.Spec.InitContainers)...)

	case "ReplicationController":
		rc := k8scorev1.ReplicationController{}
		err = json.Unmarshal(data, &rc)
		if err != nil {
			err := errors.Errorf("error unmarshalling replicationcontroller: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		if rc.Spec.Template != nil {
			containerImages = append(containerImages, readContainers(rc.Spec.Template.Spec.Containers)...)
			initContainerImages = append(initContainerImages, readContainers(rc.Spec.Template.Spec.InitContainers)...)
		}

	case "Job":
		job := k8sbatchv1.Job{}
		err = json.Unmarshal(data, &job)
		if err != nil {
			err := errors.Errorf("error unmarshalling job: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(job.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(job.Spec.Template.Spec.InitContainers)...)

	case "CronJob":
		cronjob := k8sbatchv1beta1.CronJob{}
		err = json.Unmarshal(data, &cronjob)
		if err != nil {
			err := errors.Errorf("error unmarshalling cronjob: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(cronjob.Spec.JobTemplate.Spec.Template.Spec.InitContainers)...)
	case "StatefulSet":
		ss := k8sappsv1.StatefulSet{}
		err = json.Unmarshal(data, &ss)
		if err != nil {
			err := errors.Errorf("error unmarshalling statefulset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(ss.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(ss.Spec.Template.Spec.InitContainers)...)
	case "ReplicaSet":
		rs := k8sappsv1.ReplicaSet{}
		err = json.Unmarshal(data, &rs)
		if err != nil {
			err := errors.Errorf("error unmarshalling replicaset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(rs.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(rs.Spec.Template.Spec.InitContainers)...)
	case "DaemonSet":
		ds := k8sappsv1.DaemonSet{}
		err = json.Unmarshal(data, &ds)
		if err != nil {
			err := errors.Errorf("error unmarshalling daemonset: %v", err)
			zap.S().Errorf(err.Error())
			return nil, nil, err
		}
		containerImages = append(containerImages, readContainers(ds.Spec.Template.Spec.Containers)...)
		initContainerImages = append(initContainerImages, readContainers(ds.Spec.Template.Spec.InitContainers)...)
	default:
		zap.S().Debugf("the container image extraction for kubernetes workload of kind %s is not supported.", kind)
	}
	return containerImages, initContainerImages, nil
}

// readContainers prepares list of containers and init containers from k8scorev1.Container object
func readContainers(containers []k8scorev1.Container) (containerImages []output.ContainerDetails) {
	for _, container := range containers {
		containerImages = append(containerImages, output.ContainerDetails{Name: container.Name, Image: container.Image})
	}
	return
}
