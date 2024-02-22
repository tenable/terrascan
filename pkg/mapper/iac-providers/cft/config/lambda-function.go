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

package config

import (
	"strconv"
	"strings"

	"github.com/awslabs/goformation/v7/cloudformation/lambda"
	"github.com/awslabs/goformation/v7/cloudformation/serverless"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
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
	Tags                         map[string]string    `json:"tags"`
}

// GetLambdaFunctionConfig returns config for LambdaFunction
// aws_lambda_function
func GetLambdaFunctionConfig(sf interface{}) []AWSResourceConfig {
	if l, ok := sf.(*lambda.Function); ok {
		return getLambdaConfig(l)
	}

	if s, ok := sf.(*serverless.Function); ok {
		return getServerlessConfig(s)
	}

	return nil
}

func getServerlessConfig(sf *serverless.Function) []AWSResourceConfig {
	tracingConfig := make([]TracingConfigBlock, 1)
	tracingConfig[0].Mode = functions.GetVal(sf.Tracing)

	var vpcConfig []VPCConfigBlock
	if sf.VpcConfig != nil {
		vpcConfig = make([]VPCConfigBlock, 1)

		vpcConfig[0].SecurityGroupIDs = sf.VpcConfig.SecurityGroupIds
		vpcConfig[0].SubnetIDs = sf.VpcConfig.SubnetIds
	}

	var environment []EnvironmentBlock
	if sf.Environment != nil {
		environment = make([]EnvironmentBlock, 1)

		environment[0].Variables = sf.Environment.Variables
	}

	cf := LambdaFunctionConfig{
		Config: Config{
			Name: functions.GetVal(sf.FunctionName),
		},
		FunctionName:                 functions.GetVal(sf.FunctionName),
		Role:                         functions.GetVal(sf.Role),
		Handler:                      functions.GetVal(sf.Handler),
		MemorySize:                   functions.GetVal(sf.MemorySize),
		ReservedConcurrentExecutions: functions.GetVal(sf.ReservedConcurrentExecutions),
		Runtime:                      functions.GetVal(sf.Runtime),
		Timeout:                      functions.GetVal(sf.Timeout),
		TracingConfig:                tracingConfig,
		VPCConfig:                    vpcConfig,
		Environment:                  environment,
		KMSKeyARN:                    functions.GetVal(sf.KmsKeyArn),
	}
	if sf.Tags != nil {
		cf.Tags = sf.Tags
	}
	cf = setServerlessCodePackage(cf, sf)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: sf.AWSCloudFormationMetadata,
	}}
}

func setServerlessCodePackage(cf LambdaFunctionConfig, f *serverless.Function) LambdaFunctionConfig {
	if f == nil {
		return cf
	}

	if f.ImageUri != nil {
		cf.ImageURI = functions.GetVal(f.ImageUri)
		return cf
	}

	if f.CodeUri != nil && f.CodeUri.String != nil &&
		functions.GetVal(f.CodeUri.String) != "" && !strings.HasPrefix(functions.GetVal(f.CodeUri.String), "s3") {
		cf.FileName = functions.GetVal(f.CodeUri.String)
		return cf
	}

	if f.CodeUri != nil && f.CodeUri.S3Location != nil {
		cf.S3Bucket = f.CodeUri.S3Location.Bucket
		cf.S3Key = f.CodeUri.S3Location.Key
		cf.S3ObjectVersion = strconv.Itoa(functions.GetVal(f.CodeUri.S3Location.Version))
	}

	return cf
}

func getLambdaConfig(lf *lambda.Function) []AWSResourceConfig {
	var tracingConfig []TracingConfigBlock
	if lf.TracingConfig != nil {
		tracingConfig = make([]TracingConfigBlock, 1)

		tracingConfig[0].Mode = functions.GetVal(lf.TracingConfig.Mode)
	}

	var vpcConfig []VPCConfigBlock
	if lf.VpcConfig != nil {
		vpcConfig = make([]VPCConfigBlock, 1)

		vpcConfig[0].SecurityGroupIDs = lf.VpcConfig.SecurityGroupIds
		vpcConfig[0].SubnetIDs = lf.VpcConfig.SubnetIds
	}

	var environment []EnvironmentBlock
	if lf.Environment != nil {
		environment = make([]EnvironmentBlock, 1)

		environment[0].Variables = lf.Environment.Variables
	}

	cf := LambdaFunctionConfig{
		Config: Config{
			Name: functions.GetVal(lf.FunctionName),
		},
		FunctionName:                 functions.GetVal(lf.FunctionName),
		Role:                         lf.Role,
		Handler:                      functions.GetVal(lf.Handler),
		MemorySize:                   functions.GetVal(lf.MemorySize),
		ReservedConcurrentExecutions: functions.GetVal(lf.ReservedConcurrentExecutions),
		Runtime:                      functions.GetVal(lf.Runtime),
		Timeout:                      functions.GetVal(lf.Timeout),
		TracingConfig:                tracingConfig,
		VPCConfig:                    vpcConfig,
		Environment:                  environment,
		KMSKeyARN:                    functions.GetVal(lf.KmsKeyArn),
	}

	cf = setLambdaCodePackage(cf, lf)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: lf.AWSCloudFormationMetadata,
	}}
}

func setLambdaCodePackage(cf LambdaFunctionConfig, f *lambda.Function) LambdaFunctionConfig {
	if f.Code != nil {
		cf.ImageURI = functions.GetVal(f.Code.ImageUri)
		cf.S3Bucket = functions.GetVal(f.Code.S3Bucket)
		cf.S3Key = functions.GetVal(f.Code.S3Key)
		cf.S3ObjectVersion = functions.GetVal(f.Code.S3ObjectVersion)
		return cf
	}
	return cf
}
