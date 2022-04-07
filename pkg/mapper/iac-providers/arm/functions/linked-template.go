package functions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/utils"
	getter "github.com/hashicorp/go-getter"
	"go.uber.org/zap"
)

// ResolveLinkedTemplate downloads temlate for the given uri and returns its path
func ResolveLinkedTemplate(uri string) ([]byte, error) {
	tempDir := utils.GenerateTempDir()
	defer os.RemoveAll(tempDir)
	path, err := downloadTemplate(uri, tempDir)
	if err != nil {
		return nil, err
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileData, nil

}

func downloadTemplate(uri string, dst string) (string, error) {
	parts := strings.Split(uri, "/")
	path := filepath.Join(dst, parts[len(parts)-1])
	client := getter.Client{
		Src:  uri,
		Dst:  path,
		Mode: getter.ClientModeFile,
	}
	err := client.Get()
	if err != nil {
		zap.S().Debug("unable to parse linked termplate", zap.Error(err), zap.String("file", path))
		return "", err
	}
	return path, nil
}
