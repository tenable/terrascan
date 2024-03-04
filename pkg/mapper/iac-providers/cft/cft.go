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

package cft

import (
	"errors"

	"github.com/awslabs/goformation/v7/cloudformation/applicationautoscaling"
	"github.com/awslabs/goformation/v7/cloudformation/appmesh"
	"github.com/awslabs/goformation/v7/cloudformation/athena"
	"github.com/awslabs/goformation/v7/cloudformation/autoscaling"
	"github.com/awslabs/goformation/v7/cloudformation/backup"
	"github.com/awslabs/goformation/v7/cloudformation/certificatemanager"
	"github.com/awslabs/goformation/v7/cloudformation/cloudfront"
	"github.com/awslabs/goformation/v7/cloudformation/cloudtrail"
	"github.com/awslabs/goformation/v7/cloudformation/codebuild"
	"github.com/awslabs/goformation/v7/cloudformation/cognito"
	"github.com/awslabs/goformation/v7/cloudformation/dms"
	"github.com/awslabs/goformation/v7/cloudformation/eks"
	"github.com/awslabs/goformation/v7/cloudformation/emr"
	"github.com/awslabs/goformation/v7/cloudformation/globalaccelerator"
	"github.com/awslabs/goformation/v7/cloudformation/lambda"
	"github.com/awslabs/goformation/v7/cloudformation/msk"
	"github.com/awslabs/goformation/v7/cloudformation/qldb"
	"github.com/awslabs/goformation/v7/cloudformation/ram"
	"github.com/awslabs/goformation/v7/cloudformation/sagemaker"
	"github.com/awslabs/goformation/v7/cloudformation/serverless"
	"github.com/awslabs/goformation/v7/cloudformation/sns"
	"github.com/awslabs/goformation/v7/cloudformation/sqs"
	"github.com/awslabs/goformation/v7/cloudformation/waf"

	cf "github.com/awslabs/goformation/v7/cloudformation/cloudformation"
	cnf "github.com/awslabs/goformation/v7/cloudformation/config"
	"github.com/awslabs/goformation/v7/cloudformation/ecr"
	"github.com/awslabs/goformation/v7/cloudformation/neptune"
	"github.com/awslabs/goformation/v7/cloudformation/secretsmanager"
	"github.com/awslabs/goformation/v7/cloudformation/workspaces"

	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation/efs"
	"github.com/awslabs/goformation/v7/cloudformation/elasticache"

	"github.com/awslabs/goformation/v7/cloudformation/dax"
	"github.com/awslabs/goformation/v7/cloudformation/dynamodb"
	"github.com/awslabs/goformation/v7/cloudformation/rds"

	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/awslabs/goformation/v7/cloudformation/logs"

	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/amazonmq"
	"github.com/awslabs/goformation/v7/cloudformation/apigateway"
	"github.com/awslabs/goformation/v7/cloudformation/apigatewayv2"
	"github.com/awslabs/goformation/v7/cloudformation/docdb"
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancing"
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancingv2"
	"github.com/awslabs/goformation/v7/cloudformation/elasticsearch"
	"github.com/awslabs/goformation/v7/cloudformation/guardduty"
	"github.com/awslabs/goformation/v7/cloudformation/iam"
	"github.com/awslabs/goformation/v7/cloudformation/kinesis"
	"github.com/awslabs/goformation/v7/cloudformation/kinesisfirehose"
	"github.com/awslabs/goformation/v7/cloudformation/kms"
	"github.com/awslabs/goformation/v7/cloudformation/redshift"
	"github.com/awslabs/goformation/v7/cloudformation/route53"
	"github.com/awslabs/goformation/v7/cloudformation/s3"
	"github.com/awslabs/goformation/v7/cloudformation/ssm"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/mapper/core"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/config"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/store"
	"github.com/tenable/terrascan/pkg/utils"
)

const errUnsupportedDoc = "unsupported document type"

type cftMapper struct {
}

// Mapper returns an CFT mapper for given template schema
func Mapper() core.Mapper {
	return cftMapper{}
}

// Map transforms the provider specific template to terrascan native format.
func (m cftMapper) Map(resource interface{}, params ...map[string]interface{}) ([]output.ResourceConfig, error) {
	// transform each resource and generate config
	var configs []output.ResourceConfig
	template, ok := resource.(*cloudformation.Template)
	if !ok {
		return nil, errors.New(errUnsupportedDoc)
	}
	for name, untypedRes := range template.Resources {
		for _, resourceConfig := range m.mapConfigForResource(untypedRes, name) {
			if resourceConfig.Resource != nil {
				config := output.ResourceConfig{
					Name:      name,
					SkipRules: make([]output.SkipRule, 0),
					Config:    resourceConfig.Resource,
				}

				// fill config
				if resourceConfig.Name != "" {
					config.Name = resourceConfig.Name
				}

				// determine resource type
				cfnType := untypedRes.AWSCloudFormationType()
				if resourceConfig.Type != "" {
					cfnType = cfnType + "." + resourceConfig.Type
				}
				if terraType, ok := store.ResourceTypes[cfnType]; ok {
					config.Type = terraType
					config.ID = config.Type + "." + config.Name
				} else {
					continue
				}

				// add skipRules if available
				if resourceConfig.Metadata != nil {
					skipRules := utils.ReadSkipRulesFromMap(resourceConfig.Metadata, config.ID)
					if skipRules != nil {
						config.SkipRules = append(config.SkipRules, skipRules...)
					}
				}

				configs = append(configs, config)
			}
		}
	}
	return configs, nil
}

func (m cftMapper) mapConfigForResource(r cloudformation.Resource, resourceName string) []config.AWSResourceConfig {
	switch resource := r.(type) {
	case *docdb.DBCluster:
		return config.GetDocDBConfig(resource)
	case *apigateway.RestApi:
		return config.GetAPIGatewayRestAPIConfig(resource)
	case *apigateway.Stage:
		return config.GetAPIGatewayStageConfig(resource)
	case *apigatewayv2.Stage:
		return config.GetAPIGatewayV2StageConfig(resource)
	case *apigatewayv2.Api:
		return config.GetAPIGatewayV2ApiConfig(resource)
	case *athena.WorkGroup:
		return config.GetAthenaWorkGroupConfig(resource)
	case *logs.LogGroup:
		return config.GetLogCloudWatchGroupConfig(resource)
	case *ecs.Service:
		return config.GetEcsServiceConfig(resource)
	case *dynamodb.Table:
		return config.GetDynamoDBTableConfig(resource)
	case *dax.Cluster:
		return config.GetDaxClusterConfig(resource)
	case *rds.DBInstance:
		return config.GetDBInstanceConfig(resource)
	case *rds.EventSubscription:
		return config.GetDBEventSubscriptionConfig(resource)
	case *qldb.Ledger:
		return config.GetQldbLedgerConfig(resource)
	case *ecs.Cluster:
		return config.GetEcsClusterConfig(resource)
	case *iam.Role:
		return config.GetIamRoleConfig(resource)
	case *iam.Policy:
		return config.GetIamPolicyConfig(resource)
	case *iam.AccessKey:
		return config.GetIamAccessKeyConfig(resource)
	case *iam.User:
		return config.GetIamUserConfig(resource)
	case *iam.Group:
		return config.GetIamGroupConfig(resource)
	case *rds.DBSecurityGroup:
		return config.GetDBSecurityGroupConfig(resource)
	case *ec2.SecurityGroup:
		return config.GetSecurityGroupConfig(resource)
	case *ec2.Volume:
		return config.GetEbsVolumeConfig(resource)
	case *ec2.VPC:
		return config.GetEc2VpcConfig(resource)
	case *ec2.SubnetRouteTableAssociation:
		return config.GetRouteTableAssociationConfig(resource)
	case *ec2.RouteTable:
		return config.GetRouteTableConfig(resource)
	case *ec2.NatGateway:
		return config.GetNatGatewayConfig(resource)
	case *ec2.Subnet:
		return config.GetSubnetConfig(resource)
	case *ec2.Route:
		return config.GetRouteConfig(resource)
	case *efs.FileSystem:
		return config.GetEfsFileSystemConfig(resource)
	case *elasticache.CacheCluster:
		return config.GetElastiCacheClusterConfig(resource)
	case *elasticache.ReplicationGroup:
		return config.GetElastiCacheReplicationGroupConfig(resource)
	case *amazonmq.Broker:
		return config.GetMqBorkerConfig(resource)
	case *guardduty.Detector:
		return config.GetGuardDutyDetectorConfig(resource)
	case *redshift.Cluster:
		return config.GetRedshiftClusterConfig(resource)
	case *redshift.ClusterParameterGroup:
		return config.GetRedshiftParameterGroupConfig(resource, resourceName)
	case *rds.DBCluster:
		return config.GetRDSClusterConfig(resource)
	case *route53.RecordSet:
		return config.GetRoute53RecordConfig(resource)
	case *workspaces.Workspace:
		return config.GetWorkspacesWorkspaceConfig(resource)
	case *neptune.DBCluster:
		return config.GetNeptuneClusterConfig(resource)
	case *neptune.DBInstance:
		return config.GetNeptuneClusterInstanceConfig(resource)
	case *globalaccelerator.Accelerator:
		return config.GetGlobalAcceleratorConfig(resource)
	case *waf.SizeConstraintSet:
		return config.GetWafSizeConstraintSetConfig(resource)
	case *secretsmanager.Secret:
		return config.GetSecretsManagerSecretConfig(resource)
	case *ecr.Repository:
		return config.GetEcrRepositoryConfig(resource)
	case *kms.Key:
		return config.GetKmsKeyConfig(resource)
	case *kinesis.Stream:
		return config.GetKinesisStreamConfig(resource)
	case *kinesisfirehose.DeliveryStream:
		return config.GetKinesisFirehoseDeliveryStreamConfig(resource)
	case *cf.Stack:
		return config.GetCloudFormationStackConfig(resource)
	case *cloudfront.Distribution:
		return config.GetCloudFrontDistributionConfig(resource)
	case *cloudtrail.Trail:
		return config.GetCloudTrailConfig(resource)
	case *cnf.ConfigRule:
		return config.GetConfigConfigRuleConfig(resource)
	case *cnf.ConfigurationAggregator:
		return config.GetConfigConfigurationAggregatorConfig(resource)
	case *elasticloadbalancingv2.Listener:
		return config.GetElasticLoadBalancingV2ListenerConfig(resource)
	case *elasticloadbalancingv2.TargetGroup:
		return config.GetElasticLoadBalancingV2TargetGroupConfig(resource)
	case *elasticloadbalancing.LoadBalancer:
		return config.GetElasticLoadBalancingLoadBalancerConfig(resource, resourceName)
	case *elasticsearch.Domain:
		return config.GetElasticsearchDomainConfig(resource)
	case *secretsmanager.ResourcePolicy:
		return config.GetSecretsManagerSecretPolicyConfig(resource)
	case *ecs.TaskDefinition:
		return config.GetEcsTaskDefinitionConfig(resource)
	case *s3.Bucket:
		return config.GetS3BucketConfig(resource, resourceName)
	case *s3.BucketPolicy:
		return config.GetS3BucketPolicyConfig(resource)
	case *sqs.Queue:
		return config.GetSqsQueueConfig(resource)
	case *sqs.QueuePolicy:
		return config.GetSqsQueuePolicyConfig(resource)
	case *sns.Topic:
		return config.GetSnsTopicConfig(resource)
	case *sns.TopicPolicy:
		return config.GetSnsTopicPolicyConfig(resource)
	case *autoscaling.LaunchConfiguration:
		return config.GetAutoScalingLaunchConfigurationConfig(resource)
	case *ec2.Instance:
		return config.GetEC2InstanceConfig(resource, resourceName)
	case *cognito.UserPool:
		return config.GetCognitoUserPoolConfig(resource)
	case *lambda.Function:
		return config.GetLambdaFunctionConfig(resource)
	case *serverless.Function:
		return config.GetLambdaFunctionConfig(resource)
	case *certificatemanager.Certificate:
		return config.GetCertificateManagerCertificateConfig(resource)
	case *sagemaker.NotebookInstance:
		return config.GetSagemakerNotebookInstanceConfig(resource)
	case *sagemaker.Model:
		return config.GetSagemakerModelConfig(resource)
	case *dms.ReplicationInstance:
		return config.GetDmsReplicationInstanceConfig(resource)
	case *codebuild.Project:
		return config.GetCodebuildProjectConfig(resource)
	case *emr.Cluster:
		return config.GetEmrClusterConfig(resource)
	case *msk.Cluster:
		return config.GetMskClusterConfig(resource)
	case *eks.Cluster:
		return config.GetEksClusterConfig(resource)
	case *eks.Nodegroup:
		return config.GetEksNodeGroupConfig(resource)
	case *backup.BackupVault:
		return config.GetBackupVaultConfig(resource)
	case *appmesh.Mesh:
		return config.GetAppMeshMeshConfig(resource)
	case *ram.ResourceShare:
		return config.GetRAMResourceShareConfig(resource)
	case *applicationautoscaling.ScalingPolicy:
		return config.GetAppAutoScalingPolicyConfig(resource)
	case *secretsmanager.RotationSchedule:
		return config.GetSecretsManagerSecretRotationConfig(resource)
	case *ssm.Parameter:
		return config.GetSSMParameterConfig(resource)
	case *elasticloadbalancingv2.LoadBalancer:
		return config.GetElasticLoadBalancingV2LoadBalancerConfig(resource, resourceName)
	default:
	}
	return []config.AWSResourceConfig{}
}
