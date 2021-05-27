package helm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type HelmOptions struct {
	KubeCli   kubernetes.Interface
	Name      string
	Namespace string
	Local     bool
	HelmBin   string
	SetFlags  map[string]string // use this instead of values
	DryRun    bool              // use `--dry-run` to get output of resources
	Debug     bool
	ChartPath string
}

// HelmInstallInfo is used to return helm install information
type HelmInstallInfo struct {
	DryRunManifests string
}

func (h *HelmOptions) Install(ctx context.Context) (*HelmInstallInfo, error) {
	if h.KubeCli == nil {
		if err := h.setKubeCli(); err != nil {
			return nil, err
		}
	}

	if err := h.preFlightSetting(); err != nil {
		return nil, fmt.Errorf("failed to set preflight settings, error: %w", err)
	}

	cmd := h.helmInstallCommand()
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute helm command: %s, error: %w", stdErr.String(), err)
	}

	return &HelmInstallInfo{
		DryRunManifests: string(output),
	}, nil
}

func (h *HelmOptions) setKubeCli() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to get kubernetes clientset, error: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes clientset w/ config, error: %w", err)
	}

	h.KubeCli = clientset

	return nil
}

// LoadConfig returns a kubernetes client config based on global settings.
func loadConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	homeConfig := filepath.Join(os.Getenv("HOME"), ".kube/config")
	return clientcmd.BuildConfigFromFlags("", homeConfig)

}

func (h *HelmOptions) preFlightSetting() error {
	if _, err := h.KubeCli.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: h.Namespace}}, metav1.CreateOptions{}); err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create namespace, error: %w", err)
		}
	}

	return nil
}

func (h *HelmOptions) helmInstallCommand() *exec.Cmd {
	cmd := exec.Command(h.HelmBin, "install", h.Name, h.ChartPath, "--namespace", h.Namespace, "--wait")
	if len(h.SetFlags) > 0 {
		var setFlags bytes.Buffer
		for key, value := range h.SetFlags {
			setFlags.WriteString(key + "=" + value + ",")
		}
		cmd.Args = append(cmd.Args, "--set", setFlags.String())
	}

	if h.DryRun {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	if h.Debug {
		cmd.Args = append(cmd.Args, "--debug")
	}
	return cmd
}
