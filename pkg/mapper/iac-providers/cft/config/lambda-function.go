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
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"

	"github.com/awslabs/goformation/v5/cloudformation/lambda"
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
	FileName                     string               `json:"filename"`
	FunctionName                 string               `json:"function_name"`
	Role                         string               `json:"role"`
	Handler                      string               `json:"handler"`
	SourceCodeHash               string               `json:"source_code_hash"`
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
func GetLambdaFunctionConfig(f *lambda.Function) []AWSResourceConfig {
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
		FileName:                     f.Code.ZipFile,
		FunctionName:                 f.FunctionName,
		Role:                         f.Role,
		Handler:                      f.Handler,
		SourceCodeHash:               gethash(f.Code.ZipFile),
		MemorySize:                   f.MemorySize,
		ReservedConcurrentExecutions: f.ReservedConcurrentExecutions,
		Runtime:                      f.Runtime,
		Timeout:                      f.Timeout,
		TracingConfig:                tracingConfig,
		VPCConfig:                    vpcConfig,
		Environment:                  environment,
		KMSKeyARN:                    f.KmsKeyArn,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: f.AWSCloudFormationMetadata,
	}}
}

func gethash(codefile string) string {
	data, _ := ioutil.ReadFile(codefile)
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}
