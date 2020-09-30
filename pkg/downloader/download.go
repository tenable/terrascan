package downloader

import (
	"fmt"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"go.uber.org/zap"
)

// NewDownloader returns a new downloader
func NewDownloader() *Downloader {
	return &Downloader{}
}

// list of errors
var (
	ErrEmptyURLType     = fmt.Errorf("empty remote url and type")
	ErrEmptyURLDest     = fmt.Errorf("remote url or destination dir path cannot be empty")
	ErrEmptyURLTypeDest = fmt.Errorf("empty remote url or type or desitnation dir path")
)

// GetURLWithType returns the download URL with it's respective type prefix
// along with subDir path, if present.
func GetURLWithType(remoteURL, destPath string) (string, string, error) {

	// get subDir, if present
	repoURL, subDir := SplitAddrSubdir(remoteURL)
	zap.S().Debugf("downloading %q to %q", repoURL, destPath)

	// check if a detector is present for the given url with type
	URLWithType, err := getter.Detect(repoURL, destPath, goGetterDetectors)
	if err != nil {
		zap.S().Errorf("failed to detect url with type for %q. error: '%v'", remoteURL, err)
		return "", "", err
	}
	zap.S().Debugf("remote URL: %q; url with type: %q", remoteURL, URLWithType)

	// get actual subDir path
	URLWithType, realSubDir := SplitAddrSubdir(URLWithType)
	if realSubDir != "" {
		subDir = filepath.Join(realSubDir, subDir)
	}

	if URLWithType != repoURL {
		zap.S().Debugf("detector rewrote %q to %q", repoURL, URLWithType)
	}

	// successful
	return URLWithType, subDir, nil
}

// Download retrieves the remote repository referenced in the given remoteURL
// into the destination path and then returns the full path to any subdir
// indicated in the URL
func (d Downloader) Download(remoteURL, destPath string) (string, error) {

	zap.S().Debugf("download with remote url: %q, destination dir: %q",
		remoteURL, destPath)

	// validations: remote url or destination dir path cannot be empty
	if remoteURL == "" || destPath == "" {
		zap.S().Error(ErrEmptyURLDest)
		return "", ErrEmptyURLDest
	}

	// get repository url, subdir from given remote url
	URLWithType, subDir, err := GetURLWithType(remoteURL, destPath)
	if err != nil {
		return "", err
	}

	// downloading from remote addr
	client := getter.Client{
		Src:           URLWithType,
		Dst:           destPath,
		Pwd:           destPath,
		Mode:          getter.ClientModeDir,
		Detectors:     goGetterNoDetectors,
		Decompressors: goGetterDecompressors,
		Getters:       goGetterGetters,
	}
	err = client.Get()
	if err != nil {
		zap.S().Errorf("failed to download %q. error: '%v'", URLWithType, err)
		return "", err
	}

	// Our subDir string can contain wildcards until this point, so that
	// e.g. a subDir of * can expand to one top-level directory in a .tar.gz
	// archive. Now that we've expanded the archive successfully we must
	// resolve that into a concrete path.
	finalDir := destPath
	if subDir != "" {
		finalDir, err = getter.SubdirGlob(destPath, subDir)
		if err != nil {
			zap.S().Errorf("failed to expand %q to %q", subDir, finalDir)
			return "", err
		}
		zap.S().Debugf("expanded %q to %q", subDir, finalDir)
	}

	// If we got this far then we have apparently succeeded in downloading
	// the requested object!
	return filepath.Clean(finalDir), nil
}

// DownloadWithType retrieves the remote repository referenced in the
// given remoteURL into the installation path and then returns the full path
// to any subdir indicated in the remoteURL
//
// DownloadWithType enforces download type on go-getter to get rid of any
// ambiguities in remoteURL
func (d Downloader) DownloadWithType(remoteType, remoteURL, destPath string) (string, error) {

	zap.S().Debugf("download with remote type: %q, remote URL: %q, destination dir: %q",
		remoteType, remoteURL, destPath)

	// validations
	// remoteURL and repoType cannot be empty
	if remoteURL == "" && remoteType == "" {
		zap.S().Error(ErrEmptyURLType)
		return "", ErrEmptyURLType
	}

	// remoteURL, remoteType, destination path cannot be empty
	if remoteURL == "" || remoteType == "" || destPath == "" {
		zap.S().Error(ErrEmptyURLTypeDest)
		return "", ErrEmptyURLDest
	}
	URLWithType := fmt.Sprintf("%s::%s", remoteType, remoteURL)

	// Download
	return d.Download(URLWithType, destPath)
}

// SubDirGlob returns the actual subdir with globbing processed
func SubDirGlob(destPath, subDir string) (string, error) {
	return getter.SubdirGlob(destPath, subDir)
}

// SplitAddrSubdir splits the given address into a package portion
// and a sub-directory portion.
//
// The package portion defines the URL what should be downloaded and then the
// sub-directory portion, if present, specifies a sub-directory within
// the downloaded object .
//
// The subDir portion will be returned as empty if no subdir separator
// ("//") is present in the address.
func SplitAddrSubdir(addr string) (repoURL, subDir string) {
	return getter.SourceDirSubdir(addr)
}
