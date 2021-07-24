package k8sv1

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

func (*K8sV1) getFileType(file string) string {
	if strings.HasSuffix(file, YAMLExtension) {
		return YAMLExtension
	} else if strings.HasSuffix(file, YAMLExtension2) {
		return YAMLExtension2
	} else if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}

// LoadIacDir loads all k8s files in the current directory
func (k *K8sV1) LoadIacDir(absRootDir string, nonRecursive, _ bool) (output.AllResourceConfigs, error) {
	// set the root directory being scanned
	k.absRootDir = absRootDir

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, K8sFileExtensions())
	if err != nil {
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(k.errIacLoadDirs, results.DirScanErr{IacType: "k8s", Directory: absRootDir, ErrMessage: err.Error()})
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])

			var configData output.AllResourceConfigs
			if configData, err = k.LoadIacFile(file); err != nil {
				errMsg := fmt.Sprintf("error while loading iac file '%s'. err: %v", file, err)
				zap.S().Debug("error while loading iac files", zap.String("IAC file", file), zap.Error(err))
				k.errIacLoadDirs = multierror.Append(k.errIacLoadDirs, results.DirScanErr{IacType: "k8s", Directory: fileDir, ErrMessage: errMsg})
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, k.errIacLoadDirs
}
