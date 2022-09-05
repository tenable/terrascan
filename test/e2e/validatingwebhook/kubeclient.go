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

package validatingwebhook

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "default"
)

// KubernetesClient will connect to local k8s cluster,
// and help perform resource operation
type KubernetesClient struct {
	client *kubernetes.Clientset
}

// NewKubernetesClient creates a new Kubernetes client
func NewKubernetesClient() (*KubernetesClient, error) {
	kubernetesClient := new(KubernetesClient)
	var err error
	kubernetesClient.client, err = kubernetesClient.getK8sClient()
	if err != nil {
		return nil, err
	}
	return kubernetesClient, nil
}

// getK8sClient creates a kubernetes clientset with default config path
func (k *KubernetesClient) getK8sClient() (*kubernetes.Clientset, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("home directory not found, error: %s", err.Error())
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config, error: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create K8s clientset, error: %s", err.Error())
	}

	return clientset, nil
}

// CreateValidatingWebhookConfiguration creates a ValidatingWebhookConfiguration
func (k *KubernetesClient) CreateValidatingWebhookConfiguration(webhookFile, certFile, apiKey, port string) (*admissionv1.ValidatingWebhookConfiguration, error) {
	webhooks := admissionv1.ValidatingWebhookConfiguration{}
	data, err := os.ReadFile(webhookFile)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 1024)
	err = decoder.Decode(&webhooks)
	if err != nil {
		return nil, err
	}

	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	webhook := &webhooks.Webhooks[0]
	webhook.ClientConfig.CABundle = certData
	ip, err := GetIP()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s:%s/v1/k8s/webhooks/%s/scan/validate", ip.String(), port, apiKey)
	webhook.ClientConfig.URL = &url

	admr := k.client.AdmissionregistrationV1()

	createdWebhookConfig, err := admr.ValidatingWebhookConfigurations().Create(context.TODO(), &webhooks, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createdWebhookConfig, nil
}

// DeleteValidatingWebhookConfiguration will delete the specified webhook name
func (k *KubernetesClient) DeleteValidatingWebhookConfiguration(webhookConfigName string) error {
	return k.client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(context.TODO(), webhookConfigName, metav1.DeleteOptions{})
}

// CreatePod will create a pod by parsing a resource file
func (k *KubernetesClient) CreatePod(resourceFile string) (*v1.Pod, error) {
	pod := v1.Pod{}
	data, err := os.ReadFile(resourceFile)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 1024)
	err = decoder.Decode(&pod)
	if err != nil {
		return nil, err
	}

	createdPod, err := k.client.CoreV1().Pods(namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createdPod, nil
}

// DeletePod will delete the specified pod name
func (k *KubernetesClient) DeletePod(podName string) error {
	return k.client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

// CreateService will a service by parsing a resource file
func (k *KubernetesClient) CreateService(resourceFile string) (*v1.Service, error) {
	service := v1.Service{}
	data, err := os.ReadFile(resourceFile)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 1024)
	err = decoder.Decode(&service)
	if err != nil {
		return nil, err
	}

	createdService, err := k.client.CoreV1().Services(namespace).Create(context.TODO(), &service, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createdService, nil
}

// DeleteService will delete the specified service name
func (k *KubernetesClient) DeleteService(serviceName string) error {
	return k.client.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
}

// WaitForServiceAccount will wait for the default serviceaccount to get created
func (k *KubernetesClient) WaitForServiceAccount(ctx context.Context) error {
	for {
		svcAcc, err := k.client.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
		if err != nil || ctx.Err() != nil {
			return err
		}
		if len(svcAcc.Items) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	return nil
}
