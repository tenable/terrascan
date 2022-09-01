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
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/utils"
)

// CreateTerrascanConfigFile creates a config file with test policy path
func CreateTerrascanConfigFile(configFileName, policyRootRelPath string, terrascanConfig *config.TerrascanConfig) error {
	policyAbsPath, err := filepath.Abs(policyRootRelPath)
	if err != nil {
		return err
	}

	if utils.IsWindowsPlatform() {
		policyAbsPath = strings.ReplaceAll(policyAbsPath, "\\", "\\\\")
	}

	if terrascanConfig == nil {
		terrascanConfig = &config.TerrascanConfig{}
	}

	terrascanConfig.BasePath = policyAbsPath
	terrascanConfig.RepoPath = policyAbsPath

	// create config file in work directory
	file, err := os.Create(configFileName)
	if err != nil {
		return fmt.Errorf("config file creation failed, err: %v", err)
	}

	contentBytes, err := toml.Marshal(terrascanConfig)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(contentBytes))
	if err != nil {
		return fmt.Errorf("error while writing to config file, err: %v", err)
	}
	return nil
}

// CreateCertificate creates certificates required to run server in the folder specified
func CreateCertificate(certsFolder, certFileName, privKeyFileName string) (string, string, error) {
	// create certs folder to keep certificates
	os.Mkdir(certsFolder, 0755)
	certFileAbsPath, err := filepath.Abs(filepath.Join(certsFolder, "server.crt"))
	if err != nil {
		return "", "", err
	}
	privKeyFileAbsPath, err := filepath.Abs(filepath.Join(certsFolder, "priv.key"))
	if err != nil {
		return "", "", err
	}
	err = GenerateCertificates(certFileAbsPath, privKeyFileAbsPath)
	if err != nil {
		return "", "", err
	}

	return certFileAbsPath, privKeyFileAbsPath, nil
}

// DeleteDefaultKindCluster deletes the default kind cluster
func DeleteDefaultKindCluster() error {
	cmd := exec.Command("kind", "delete", "cluster")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// CreateDefaultKindCluster creates the default kind cluster
func CreateDefaultKindCluster() error {
	cmd := exec.Command("kind", "create", "cluster")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// GetIP finds preferred outbound ip of the machine
func GetIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4(), nil
			}
		}
	}
	return nil, fmt.Errorf("could not find ip address of the machine")
}
