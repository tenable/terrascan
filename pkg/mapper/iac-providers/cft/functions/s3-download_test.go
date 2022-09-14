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
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func TestDownloadBucketObj(t *testing.T) {
	table := []struct {
		errorOperation string
		name           string
		templateURL    string
	}{
		{
			name:        "public template",
			templateURL: "https://s3.amazonaws.com/cloudformation-templates-us-east-1/S3_Bucket.template",
		}, {
			errorOperation: "HeadObject",
			name:           "private template head object error",
			templateURL:    "https://s3.amazonaws.com/cloudformation-templates-us-east-1/S3_Bucket_Not_There.template",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DownloadBucketObj(tt.templateURL)
			if tt.errorOperation != "" {
				var oe *smithy.OperationError
				if errors.As(err, &oe) {
					if oe.Operation() != tt.errorOperation {
						t.Errorf("unexpected operation; got: '%+v'", oe.Operation())
					}
				} else {
					t.Errorf("unexpected error; got: '%+v'", reflect.TypeOf(err))
				}
			}
		})
	}
}

type mockHeadObjectAPI func(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error)

func (m mockHeadObjectAPI) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m(ctx, params, optFns...)
}

type mockS3DownloadManager func(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)

func (m mockS3DownloadManager) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error) {
	return m(ctx, w, input, options...)
}

func TestDownloadPrivateBucketObj(t *testing.T) {

	table := []struct {
		client         func(t *testing.T) HeadBucketAPIClient
		manager        func(t *testing.T) S3DownloadManager
		errorOperation string
		name           string
		templateURL    string
	}{
		{
			client: func(t *testing.T) HeadBucketAPIClient {
				return mockHeadObjectAPI(func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
					t.Helper()
					return &s3.HeadObjectOutput{ContentLength: 64}, nil
				})
			},
			manager: func(t *testing.T) S3DownloadManager {
				return mockS3DownloadManager(func(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error) {
					t.Helper()
					return 0, fmt.Errorf("error in download operation")
				})
			},
			errorOperation: "error in download operation",
			name:           "read obj error",
			templateURL:    "https://s3.amazonaws.com/cloudformation-templates-us-east-1/S3_Bucket_Not_There.template",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.templateURL)
			s3c := S3Client{
				client:     tt.client(t),
				downloader: tt.manager(t),
			}
			_, err := downloadPrivateTemplate(u, &s3c)
			if tt.errorOperation != "" {
				if err.Error() != tt.errorOperation {
					t.Errorf("unexpected error; got: '%+v'", reflect.TypeOf(err))
				}
			}
		})
	}
}
