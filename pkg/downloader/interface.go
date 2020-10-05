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

package downloader

// Downloader helps in downloading different kinds of modules from
// different types of sources
type Downloader interface {
	Download(url, destDir string) (finalDir string, err error)
	DownloadWithType(remoteType, url, dest string) (finalDir string, err error)
	GetURLSubDir(url, dest string) (urlWithType string, subDir string, err error)
	SubDirGlob(string, string) (string, error)
}

// NewDownloader returns a new downloader
func NewDownloader() Downloader {
	return NewGoGetter()
}
