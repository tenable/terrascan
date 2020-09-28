package downloader

import (
	"path/filepath"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"
	"go.uber.org/zap"
)

// list of supported detectors
var goGetterDetectors = []getter.Detector{
	new(getter.GitHubDetector),
	new(getter.GitDetector),
	new(getter.BitBucketDetector),
	new(getter.GCSDetector),
	new(getter.S3Detector),
	new(getter.FileDetector),
}

// empty list of detectors
var goGetterNoDetectors = []getter.Detector{}

// list of supported decompressors
var goGetterDecompressors = map[string]getter.Decompressor{
	"bz2": new(getter.Bzip2Decompressor),
	"gz":  new(getter.GzipDecompressor),
	"xz":  new(getter.XzDecompressor),
	"zip": new(getter.ZipDecompressor),

	"tar.bz2":  new(getter.TarBzip2Decompressor),
	"tar.tbz2": new(getter.TarBzip2Decompressor),

	"tar.gz": new(getter.TarGzipDecompressor),
	"tgz":    new(getter.TarGzipDecompressor),

	"tar.xz": new(getter.TarXzDecompressor),
	"txz":    new(getter.TarXzDecompressor),
}

// list of supported getters
var goGetterGetters = map[string]getter.Getter{
	"file":  new(getter.FileGetter),
	"gcs":   new(getter.GCSGetter),
	"git":   new(getter.GitGetter),
	"hg":    new(getter.HgGetter),
	"s3":    new(getter.S3Getter),
	"http":  getterHTTPGetter,
	"https": getterHTTPGetter,
}

var getterHTTPClient = cleanhttp.DefaultClient()

var getterHTTPGetter = &getter.HttpGetter{
	Client: getterHTTPClient,
	Netrc:  true,
}

// NewDownloader returns a new downloader
func NewDownloader() *Downloader {
	return &Downloader{}
}

// Download retrieves the package referenced in the given address
// into the installation path and then returns the full path to any subdir
// indicated in the address.
func (d Downloader) Download(addr, instPath string) (string, error) {

	// get subDir, if present
	packageAddr, subDir := SplitAddrSubdir(addr)
	zap.S().Debugf("downloading %q to %q", packageAddr, instPath)

	// check if a detector is present for the given address
	realAddr, err := getter.Detect(packageAddr, instPath, goGetterDetectors)
	if err != nil {
		zap.S().Errorf("failed to detect resolved address for %q. error: '%v'", addr, err)
		return "", err
	}
	zap.S().Debugf("resolved address: %q; provider address: %q", realAddr, addr)

	// get actual subDir path
	var realSubDir string
	realAddr, realSubDir = SplitAddrSubdir(realAddr)
	if realSubDir != "" {
		subDir = filepath.Join(realSubDir, subDir)
	}

	if realAddr != packageAddr {
		zap.S().Debugf("detector rewrote %q to %q", packageAddr, realAddr)
	}

	// downloading from remote addr
	client := getter.Client{
		Src:           realAddr,
		Dst:           instPath,
		Pwd:           instPath,
		Mode:          getter.ClientModeDir,
		Detectors:     goGetterNoDetectors, // we already did detection above
		Decompressors: goGetterDecompressors,
		Getters:       goGetterGetters,
	}
	err = client.Get()
	if err != nil {
		zap.S().Errorf("failed to download %q. error: '%v'", realAddr, err)
		return "", err
	}

	// Our subDir string can contain wildcards until this point, so that
	// e.g. a subDir of * can expand to one top-level directory in a .tar.gz
	// archive. Now that we've expanded the archive successfully we must
	// resolve that into a concrete path.
	finalDir := instPath
	if subDir != "" {
		finalDir, err = getter.SubdirGlob(instPath, subDir)
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

// SplitAddrSubdir splits the given address (which is assumed to be a
// registry address or go-getter-style address) into a package portion
// and a sub-directory portion.
//
// The package portion defines what should be downloaded and then the
// sub-directory portion, if present, specifies a sub-directory within
// the downloaded object (an archive, VCS repository, etc) that contains
// the module's configuration files.
//
// The subDir portion will be returned as empty if no subdir separator
// ("//") is present in the address.
func SplitAddrSubdir(addr string) (packageAddr, subDir string) {
	return getter.SourceDirSubdir(addr)
}

var localSourcePrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

// IsLocalPath returns true if the given "addr" is a local file path
func IsLocalPath(addr string) bool {
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}
