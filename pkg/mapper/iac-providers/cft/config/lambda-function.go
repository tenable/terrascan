/*
    Copyright (C) 2022 Accurics, Inc.

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

package config

import (
	"strconv"
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation/lambda"
	"github.com/awslabs/goformation/v5/cloudformation/serverless"
)

// TracingConfigBlock holds config for TracingConfig
type TracingConfigBlock struct {
	Mode string `json:"mode"`
}

// VPCConfigBlock holds config for VPCConfig
type VPCConfigBlock struct {
	SecurityGroupIDs []string `json:"security_group_ids"`
	SubnetIDs        []string `json:"subnet_ids"`
}

// EnvironmentBlock holds config for Environment variables
type EnvironmentBlock struct {
	Variables map[string]string `json:"variables"`
}

// LambdaFunctionConfig holds config for LambdaFunction
type LambdaFunctionConfig struct {
	Config
	ImageURI                     string               `json:"image_uri,omitempty"`
	FileName                     string               `json:"filename,omitempty"`
	S3Bucket                     string               `json:"s3_bucket,omitempty"`
	S3Key                        string               `json:"s3_key,omitempty"`
	S3ObjectVersion              string               `json:"s3_object_version,omitempty"`
	FunctionName                 string               `json:"function_name"`
	Role                         string               `json:"role"`
	Handler                      string               `json:"handler"`
	MemorySize                   int                  `json:"memory_size"`
	ReservedConcurrentExecutions int                  `json:"reserved_concurrent_executions"`
	Runtime                      string               `json:"runtime"`
	Timeout                      int                  `json:"timeout"`
	TracingConfig                []TracingConfigBlock `json:"tracing_config"`
	VPCConfig                    []VPCConfigBlock     `json:"vpc_config"`
	Environment                  []EnvironmentBlock   `json:"environment"`
	KMSKeyARN                    string               `json:"kms_key_arn"`
}

// GetLambdaFunctionConfig returns config for LambdaFunction
func GetLambdaFunctionConfig(sf interface{}) []AWSResourceConfig {
	if l, ok := sf.(*lambda.Function); ok {
		return getLambdaConfig(l)
	}

	if s, ok := sf.(*serverless.Function); ok {
		return getServerlessConfig(s)
	}

	return []AWSResourceConfig{{}}
}

func getServerlessConfig(f *serverless.Function) []AWSResourceConfig {
	tracingConfig := make([]TracingConfigBlock, 1)
	tracingConfig[0].Mode = f.Tracing

	var vpcConfig []VPCConfigBlock
	if f.VpcConfig != nil {
		vpcConfig = make([]VPCConfigBlock, 1)

		vpcConfig[0].SecurityGroupIDs = f.VpcConfig.SecurityGroupIds
		vpcConfig[0].SubnetIDs = f.VpcConfig.SubnetIds
	}

	var environment []EnvironmentBlock
	if f.Environment != nil {
		environment = make([]EnvironmentBlock, 1)

		environment[0].Variables = f.Environment.Variables
	}

	cf := LambdaFunctionConfig{
		Config: Config{
			Name: f.FunctionName,
		},
		FunctionName:                 f.FunctionName,
		Role:                         f.Role,
		Handler:                      f.Handler,
		MemorySize:                   f.MemorySize,
		ReservedConcurrentExecutions: f.ReservedConcurrentExecutions,
		Runtime:                      f.Runtime,
		Timeout:                      f.Timeout,
		TracingConfig:                tracingConfig,
		VPCConfig:                    vpcConfig,
		Environment:                  environment,
		KMSKeyARN:                    f.KmsKeyArn,
	}

	cf = setServerlessCodePackage(cf, f)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: f.AWSCloudFormationMetadata,
	}}
}

func setServerlessCodePackage(cf LambdaFunctionConfig, f *serverless.Function) LambdaFunctionConfig {
	if f.ImageUri != "" {
		cf.ImageURI = f.ImageUri
		return cf
	}

	if *f.CodeUri.String != "" && !strings.HasPrefix(*f.CodeUri.String, "s3") {
		cf.FileName = *f.CodeUri.String
		return cf
	}

	cf.S3Bucket = f.CodeUri.S3Location.Bucket
	cf.S3Key = f.CodeUri.S3Location.Key
	cf.S3ObjectVersion = strconv.Itoa(f.CodeUri.S3Location.Version)
	return cf
}

func getLambdaConfig(f *lambda.Function) []AWSResourceConfig {
	var tracingConfig []TracingConfigBlock
	if f.TracingConfig != nil {
		tracingConfig = make([]TracingConfigBlock, 1)

		tracingConfig[0].Mode = f.TracingConfig.Mode
	}

	var vpcConfig []VPCConfigBlock
	if f.VpcConfig != nil {
		vpcConfig = make([]VPCConfigBlock, 1)

		vpcConfig[0].SecurityGroupIDs = f.VpcConfig.SecurityGroupIds
		vpcConfig[0].SubnetIDs = f.VpcConfig.SubnetIds
	}

	var environment []EnvironmentBlock
	if f.Environment != nil {
		environment = make([]EnvironmentBlock, 1)

		environment[0].Variables = f.Environment.Variables
	}

	cf := LambdaFunctionConfig{
		Config: Config{
			Name: f.FunctionName,
		},
		FunctionName:                 f.FunctionName,
		Role:                         f.Role,
		Handler:                      f.Handler,
		MemorySize:                   f.MemorySize,
		ReservedConcurrentExecutions: f.ReservedConcurrentExecutions,
		Runtime:                      f.Runtime,
		Timeout:                      f.Timeout,
		TracingConfig:                tracingConfig,
		VPCConfig:                    vpcConfig,
		Environment:                  environment,
		KMSKeyARN:                    f.KmsKeyArn,
	}

	cf = setLambdaCodePackage(cf, f)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: f.AWSCloudFormationMetadata,
	}}
}

func setLambdaCodePackage(cf LambdaFunctionConfig, f *lambda.Function) LambdaFunctionConfig {
	if f.Code.ImageUri != "" {
		cf.ImageURI = f.Code.ImageUri
		return cf
	}

	cf.S3Bucket = f.Code.S3Bucket
	cf.S3Key = f.Code.S3Key
	cf.S3ObjectVersion = f.Code.S3ObjectVersion
	return cf
}
