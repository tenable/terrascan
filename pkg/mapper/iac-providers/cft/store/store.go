/*
    Copyright (C) 2021 Accurics, Inc.

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

package store

// ResourceTypes holds mapping for CFT resource types to TF types
var ResourceTypes = map[string]string{
	"AWS::DocDB::DBCluster":                    AwsDocDBCluster,
	"AWS::ApiGatewayV2::Stage":                 AwsAPIGatewayV2Stage,
	"AWS::ApiGateway::Stage":                   AwsAPIGatewayStage,
	"AWS::ApiGateway::Stage.MethodSettings":    AwsAPIGatewayStageMethodSettings,
	"AWS::ApiGateway::RestApi":                 AwsAPIGatewayRestAPI,
	"AWS::ECS::Service":                        AwsEcsService,
	"AWS::Logs::LogGroup":                      AwsLogGroup,
	"AWS::DynamoDB::Table":                     AwsDynamoDBTable,
	"AWS::DAX::Cluster":                        AwsDaxCluster,
	"AWS::RDS::DBInstance":                     AwsDBInstance,
	"AWS::IAM::Role":                           AwsIamRole,
	"AWS::IAM::Role.Policy":                    AwsIamRolePolicy,
	"AWS::IAM::Group.Policy":                   AwsIamGroupPolicy,
	"AWS::IAM::Policy":                         AwsIamPolicy,
	"AWS::IAM::AccessKey":                      AwsIamAccessKey,
	"AWS::IAM::User":                           AwsIamUser,
	"AWS::IAM::User.LoginProfile":              AwsIamUserLoginProfile,
	"AWS::IAM::User.Policy":                    AwsIamUserPolicy,
	"AWS::RDS::DBSecurityGroup":                AwsDBSecurityGroup,
	"AWS::EC2::Volume":                         AwsEbsVolume,
	"AWS::EFS::FileSystem":                     AwsEfsFileSystem,
	"AWS::ElastiCache::CacheCluster":           AwsElastiCacheCluster,
	"AWS::ElastiCache::ReplicationGroup":       AwsElastiCacheReplicationGroup,
	"AWS::GuardDuty::Detector":                 AwsGuardDutyDetector,
	"AWS::AmazonMQ::Broker":                    AwsMqBroker,
	"AWS::Redshift::Cluster":                   AwsRedshiftCluster,
	"AWS::RDS::DBCluster":                      AwsRdsCluster,
	"AWS::Route53::RecordSet":                  AwsRoute53Record,
	"AWS::EC2::SecurityGroup":                  AwsSecurityGroup,
	"AWS::WorkSpaces::Workspace":               AwsWorkspacesWorkspace,
	"AWS::Neptune::DBCluster":                  AwsNeptuneCluster,
	"AWS::SecretsManager::Secret":              AwsSecretsManagerSecret,
	"AWS::ECR::Repository":                     AwsEcrRepository,
	"AWS::KMS::Key":                            AwsKmsKey,
	"AWS::Kinesis::Stream":                     AwsKinesisStream,
	"AWS::KinesisFirehose::DeliveryStream":     AwsKinesisFirehoseDeliveryStream,
	"AWS::CloudFormation::Stack":               AwsCloudFormationStack,
	"AWS::CloudFront::Distribution":            AwsCloudFrontDistribution,
	"AWS::CloudTrail::Trail":                   AwsCloudTrail,
	"AWS::Config::ConfigRule":                  AwsConfigConfigRule,
	"AWS::Config::ConfigurationAggregator":     AwsConfigConfigurationAggregator,
	"AWS::ElasticLoadBalancingV2::Listener":    AwsLbListener,
	"AWS::ElasticLoadBalancingV2::TargetGroup": AwsLbTargetGroup,
	"AWS::ElasticLoadBalancing::LoadBalancer":  AwsElb,
	"AWS::Elasticsearch::Domain":               AwsElasticsearchDomain,
	"AWS::Elasticsearch::Domain.Policy":        AwsElasticsearchDomainPolicy,
	"AWS::EFS::FileSystem.FileSystemPolicy":    AwsEfsFileSystemPolicy,
	"AWS::SecretsManager::ResourcePolicy":      AwsSecretsManagerResourcePolicy,
	"AWS::ECS::TaskDefinition":                 AwsEcsTaskDefinition,
	"AWS::S3::Bucket":                          AwsS3Bucket,
	"AWS::S3::Bucket.PublicAccessBlock":        AwsS3BucketPublicAccessBlock,
	"AWS::S3::BucketPolicy":                    AwsS3BucketPolicy,
	"AWS::SQS::Queue":                          AwsSqsQueue,
	"AWS::SQS::QueuePolicy":                    AwsSqsQueuePolicy,
	"AWS::SNS::Topic":                          AwsSnsTopic,
	"AWS::SNS::TopicPolicy":                    AwsSnsTopicPolicy,
	"AWS::AutoScaling::LaunchConfiguration":    AwsLaunchConfiguration,
	"AWS::EC2::Instance":                       AwsEc2Instance,
	"AWS::EC2::Instance.NetworkInterface":      AwsEc2NetworkInterface,
	"AWS::Cognito::UserPool":                   AwsCognitoUserPool,
	"AWS::Lambda::Function":                    AwsLambdaFunction,
	"AWS::CertificateManager::Certificate":     AwsAcmCertificate,
}
