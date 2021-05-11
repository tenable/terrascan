package validatingwebhook

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

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

type KubernetesClient struct {
	client *kubernetes.Clientset
}

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
		fmt.Println("home directory not found", err)
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
	data, err := ioutil.ReadFile(webhookFile)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 1024)
	err = decoder.Decode(&webhooks)
	if err != nil {
		return nil, err
	}

	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	webhook := &webhooks.Webhooks[0]
	webhook.ClientConfig.CABundle = certData
	ip, err := getIP()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s:%s/v1/k8s/webhooks/%s/scan/validate", ip.String(), port, apiKey)
	webhook.ClientConfig.URL = &url

	admr := k.client.AdmissionregistrationV1()

	createdWebhookConfig, err := admr.ValidatingWebhookConfigurations().Create(context.TODO(), &webhooks, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return createdWebhookConfig, nil
}

func (k *KubernetesClient) DeleteValidatingWebhookConfiguration(webhookConfigName string) error {
	return k.client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(context.TODO(), webhookConfigName, metav1.DeleteOptions{})
}

func (k *KubernetesClient) CreatePod(resourceFile string) (*v1.Pod, error) {
	pod := v1.Pod{}
	data, err := ioutil.ReadFile(resourceFile)
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
		fmt.Println(err)
		return nil, err
	}
	return createdPod, err
}

func (k *KubernetesClient) DeletePod(podName string) error {
	return k.client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

func (k *KubernetesClient) CreateService(resourceFile string) (*v1.Service, error) {
	service := v1.Service{}
	data, err := ioutil.ReadFile(resourceFile)
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
		fmt.Println(err)
		return nil, err
	}
	return createdService, err
}

func (k *KubernetesClient) DeleteService(serviceName string) error {
	return k.client.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
}
