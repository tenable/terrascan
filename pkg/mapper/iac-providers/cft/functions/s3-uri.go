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

package functions

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// borrowed from https://gist.github.com/kwilczynski/f6e626990d6d2395b42a12721b165b86
// modified per need

// S3URI holds the metadata for s3 url
type S3URI struct {
	uri       *url.URL
	VersionID *string
	Scheme    *string
	Bucket    *string
	Key       *string
}

// String returns pointer to the string object
func String(s string) *string {
	return &s
}

var (
	errBucketNotFound    = errors.New("bucket name could not be found")
	errHostnameNotFound  = errors.New("hostname could not be found")
	errInvalidS3Endpoint = errors.New("an invalid S3 endpoint URL")

	// Pattern used to parse multiple path and host style S3 endpoint URLs.
	s3URLPattern = regexp.MustCompile(`^(.+\.)?s3[.-](?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)\.`)
)

// ParseS3URI returns s3 metadata given url
func ParseS3URI(u *url.URL) (*S3URI, error) {
	s3u := &S3URI{
		uri: u,
	}

	switch u.Scheme {
	case "s3", "http", "https":
		s3u.Scheme = String(u.Scheme)
	default:
		return nil, fmt.Errorf("unable to parse schema type: %s", u.Scheme)
	}

	// Handle S3 endpoint URL with the schema s3:// that is neither
	// the host style nor the path style.
	if u.Scheme == "s3" {
		if u.Host == "" {
			return nil, errBucketNotFound
		}
		s3u.Bucket = String(u.Host)

		if u.Path != "" && u.Path != "/" {
			s3u.Key = String(strings.TrimLeft(u.Path, "/"))
		}

		return s3u, nil
	}

	if u.Host == "" {
		return nil, errHostnameNotFound
	}

	matches := s3URLPattern.FindStringSubmatch(u.Host)
	if len(matches) < 1 {
		return nil, errInvalidS3Endpoint
	}

	prefix := matches[1]

	if prefix == "" {
		if u.Path != "" && u.Path != "/" {
			u.Path = u.Path[1:len(u.Path)]

			index := strings.Index(u.Path, "/")
			switch {
			case index == -1:
				s3u.Bucket = String(u.Path)
			case index == len(u.Path)-1:
				s3u.Bucket = String(u.Path[:index])
			default:
				s3u.Bucket = String(u.Path[:index])
				s3u.Key = String(u.Path[index+1:])
			}
		}
	} else {
		s3u.Bucket = String(prefix[:len(prefix)-1])

		if u.Path != "" && u.Path != "/" {
			s3u.Key = String(u.Path[1:len(u.Path)])
		}
	}

	// Query string used when requesting a particular version of a given
	// S3 object (key).
	const versionID = "versionID"
	if s := u.Query().Get(versionID); s != "" {
		s3u.VersionID = String(s)
	}

	return s3u, nil
}
