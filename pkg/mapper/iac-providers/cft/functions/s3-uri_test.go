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
	"net/url"
	"reflect"
	"testing"
)

const (
	bucket = "bucket"
	key    = "key"
)

func TestParseS3URI(t *testing.T) {
	table := []struct {
		wantErr     error
		name        string
		templateURL string
	}{
		{
			name:        "https host style",
			templateURL: "http://bucket.s3-region.amazonaws.com/key",
			wantErr:     nil,
		}, {
			name:        "https path style",
			templateURL: "http://s3-region.amazonaws.com/bucket/key",
			wantErr:     nil,
		}, {
			name:        "dualstack 1",
			templateURL: "https://s3.dualstack.region.amazonaws.com/bucket/key",
			wantErr:     nil,
		}, {
			name:        "dualstack 2",
			templateURL: "http://bucket.s3.dualstack.region.amazonaws.com/key",
			wantErr:     nil,
		}, {
			name:        "static 1",
			templateURL: "http://bucket.s3-website.region.amazonaws.com/key",
			wantErr:     nil,
		}, {
			name:        "static 2",
			templateURL: "http://bucket.s3-website-region.amazonaws.com/key",
			wantErr:     nil,
		}, {
			name:        "s3 1",
			templateURL: "https://s3.region.amazonaws.com/bucket/key",
			wantErr:     nil,
		}, {
			name:        "s3 2",
			templateURL: "http://s3-region.amazonaws.com/bucket/key",
			wantErr:     nil,
		}, {
			name:        "s3 3",
			templateURL: "https://s3.amazonaws.com/bucket/key",
			wantErr:     nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.templateURL)
			s3u, err := ParseS3URI(u)
			if err != nil {
				if err != tt.wantErr {
					t.Errorf("unexpected error; got: '%+v'", reflect.TypeOf(err))
				}
			} else {
				if *s3u.Bucket != bucket || *s3u.Key != key {
					t.Errorf("unexpected metadata; got '%+v'", s3u)
				}
			}
		})
	}
}
