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
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/ptr"
	getter "github.com/hashicorp/go-getter"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// HeadBucketAPIClient is an S3 API client that can invoke the HeadBucket operation.
type HeadBucketAPIClient interface {
	HeadObject(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// S3DownloadManager is an S3 manager that can invoke the Download operation.
type S3DownloadManager interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)
}

// S3Client struct is used to hold s3.Client, manager.Downloader corresponding interfaces
type S3Client struct {
	client     HeadBucketAPIClient
	downloader S3DownloadManager
}

// NewS3Client returns S3Client initialized with AWS credentials
func NewS3Client() (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		zap.S().Debug("error loading AWS credentials for bucket", err)
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Client{
		client:     client,
		downloader: manager.NewDownloader(client),
	}, nil
}

// DownloadBucketObj returns the content for S3 bucket object
func DownloadBucketObj(templateURL string) ([]byte, error) {
	s3URL, err := url.Parse(templateURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse given S3 endpoint URL: %w", err)
	}

	// assuming url points to public object
	switch s3URL.Scheme {
	case "http", "https":
		buf, err := downloadPublicTemplate(templateURL)
		if err != nil {
			zap.S().Debug("the s3 url for nested stack is not a public object", zap.String("url", templateURL), err)
		} else {
			return buf, nil
		}
	}

	// if not public get bucket name and key
	s3c, err := NewS3Client()
	if err != nil {
		zap.S().Debug("error loading AWS credentials for bucket", zap.String("url", templateURL))
		return nil, err
	}

	return downloadPrivateTemplate(s3URL, s3c)
}

func downloadPrivateTemplate(url *url.URL, s3c *S3Client) ([]byte, error) {
	s3URI, err := ParseS3URI(url)
	if err != nil {
		zap.S().Debug("error parsing S3 uri", s3URI, err)
		return nil, err
	}

	// get size and check access
	headInput := &s3.HeadObjectInput{
		Bucket: s3URI.Bucket,
		Key:    s3URI.Key,
	}
	if s3URI.VersionID != nil {
		headInput.VersionId = s3URI.VersionID
	}

	headObject, err := s3c.client.HeadObject(context.TODO(), headInput)
	if err != nil {
		zap.S().Debug("error in HEAD operation for bucket object", err)
		return nil, err
	}
	buf := make([]byte, ptr.ToInt64(headObject.ContentLength))
	w := manager.NewWriteAtBuffer(buf)

	// get the object
	downloaderInput := &s3.GetObjectInput{
		Bucket: s3URI.Bucket,
		Key:    s3URI.Key,
	}
	if s3URI.VersionID != nil {
		downloaderInput.VersionId = s3URI.VersionID
	}

	_, err = s3c.downloader.Download(context.TODO(), w, downloaderInput)
	if err != nil {
		zap.S().Debug("error downloading bucket object for uri", s3URI)
		return nil, err
	}

	return buf, nil

}

func downloadPublicTemplate(uri string) ([]byte, error) {
	dst := utils.GenerateTempDir()
	defer os.RemoveAll(dst)
	parts := strings.Split(uri, "/")
	path := filepath.Join(dst, parts[len(parts)-1])

	client := getter.Client{
		Src:  uri,
		Dst:  path,
		Mode: getter.ClientModeFile,
	}
	err := client.Get()
	if err != nil {
		zap.S().Debug("unable to parse linked template", zap.Error(err), zap.String("file", path))
		return nil, err
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileData, nil

}
