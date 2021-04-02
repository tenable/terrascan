package downloader

import (
	"fmt"
	"net/url"

	"git.fuzzbuzz.io/fuzz"
)

func FuzzGetURLSubDir(f *fuzz.F) {
	// func (g *goGetter) GetURLSubDir(remoteURL, destPath string) (string, string, error) {
	g := newGoGetter()

	remoteURL := f.String("remoteURL").Seeds("github.com/accurics/terrascan",
		"github.com/accurics/terrascan//some-subdir",
		"git@github.com:accurics/terrascan.git//some-subdir",
		"git::ssh://username@example.com/storage.git//some-subdir").Get()
	parsedURL, _, err := g.GetURLSubDir(remoteURL, f.String("destPath").Get())
	if err == nil {
		_, parseErr := url.Parse(parsedURL)
		if parseErr != nil {
			f.Fail(fmt.Sprintf("url from getURLSubDir invalid: %v", remoteURL))
		}
	}
}
