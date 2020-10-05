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

package tfv12

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/downloader"
	"go.uber.org/zap"
)

var localSourcePrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

func isLocalSourceAddr(addr string) bool {
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}

// RemoteModuleInstaller helps in downloading remote modules, it also maintains a
// cache of all the installed modules and their respective resolved addresses
// (URL)
type RemoteModuleInstaller struct {
	cache      InstalledCache
	downloader downloader.Downloader
}

// NewRemoteModuleInstaller returns a RemoteModuleInstaller initialized with a
// new cache and downloader
func NewRemoteModuleInstaller() *RemoteModuleInstaller {
	return &RemoteModuleInstaller{
		cache:      make(map[string]string),
		downloader: downloader.NewDownloader(),
	}
}

// InstalledCache remembers the final resolved addresses of all the sources
// already downloaded.
//
// The keys in InstalledCache are resolved and trimmed source addresses
// (with a scheme always present, and without any "subdir" component),
// and the values are the paths where each source was previously installed.
type InstalledCache map[string]string

// DownloadModule retrieves the package referenced in the given address
// into the installation path and then returns the full path to any subdir
// indicated in the address.
func (r *RemoteModuleInstaller) DownloadModule(addr, destPath string) (string, error) {

	// split url and subdir
	URLWithType, subDir, err := r.downloader.GetURLSubDir(addr, destPath)
	if err != nil {
		return "", err
	}

	// check if the module has already been downloaded
	if prevDir, exists := r.cache[URLWithType]; exists {
		zap.S().Debugf("module %q already installed at %q", URLWithType, prevDir)
		destPath = prevDir
	} else {
		destPath, err := r.downloader.Download(URLWithType, destPath)
		if err != nil {
			zap.S().Debugf("failed to download remote module. error: '%v'", err)
			return "", err
		}
		// Remember where we installed this so we might reuse this directory
		// on subsequent calls to avoid re-downloading.
		r.cache[URLWithType] = destPath
	}

	// Our subDir string can contain wildcards until this point, so that
	// e.g. a subDir of * can expand to one top-level directory in a .tar.gz
	// archive. Now that we've expanded the archive successfully we must
	// resolve that into a concrete path.
	var finalDir string
	if subDir != "" {
		finalDir, err = r.downloader.SubDirGlob(destPath, subDir)
		if err != nil {
			return "", err
		}
		zap.S().Debugf("expanded %q to %q", subDir, finalDir)
	} else {
		finalDir = destPath
	}

	// If we got this far then we have apparently succeeded in downloading
	// the requested object!
	return filepath.Clean(finalDir), nil
}

// CleanUp cleans up all the locally downloaded modules
func (r *RemoteModuleInstaller) CleanUp() {
	for url, path := range r.cache {
		zap.S().Debugf("deleting %q installed at %q", url, path)
		os.RemoveAll(path)
	}
}
