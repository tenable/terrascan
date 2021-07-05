package k8sv1

import (
	"encoding/json"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	yamltojson "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	k8sappsv1 "k8s.io/api/apps/v1"
	k8sbatchv1 "k8s.io/api/batch/v1"
	k8sbatchv1beta1 "k8s.io/api/batch/v1beta1"
	k8scorev1 "k8s.io/api/core/v1"
)

func (k *K8sV1) extractContainerImages(kind string, doc *utils.IacDocument) ([]output.ContainerNameAndImage, []output.ContainerNameAndImage, error) {
	var containerImages = make([]output.ContainerNameAndImage, 0)
	var initContainerImages = make([]output.ContainerNameAndImage, 0)
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range pod.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range deployment.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range rc.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range job.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range ss.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range rs.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
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
			containerImages = append(containerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
		for _, container := range ds.Spec.Template.Spec.InitContainers {
			initContainerImages = append(initContainerImages, output.ContainerNameAndImage{Name: container.Name, Image: container.Image})
		}
	default:
		zap.S().Debugf("the container image extraction for kubernetes workload of kind %s is not supported.", kind)
	}
	return containerImages, initContainerImages, nil
}
