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
	"github.com/awslabs/goformation/v7/cloudformation/cloudfront"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// CloudFrontDistributionConfig holds config for aws_cloudfront_distribution
type CloudFrontDistributionConfig struct {
	Config
	Restrictions         interface{} `json:"restrictions,omitempty"`
	OrderedCacheBehavior interface{} `json:"ordered_cache_behavior,omitempty"`
	LoggingConfig        interface{} `json:"logging_config,omitempty"`
	ViewerCertificate    interface{} `json:"viewer_certificate,omitempty"`
	WebACLId             string      `json:"web_acl_id,omitempty"`
}

// GetCloudFrontDistributionConfig returns config for aws_cloudfront_distribution
// aws_cloudfront_distribution
func GetCloudFrontDistributionConfig(d *cloudfront.Distribution) []AWSResourceConfig {
	cf := CloudFrontDistributionConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(d.Tags),
		},
	}
	if checkDistributionConfig(d.DistributionConfig) {
		restrictions := make([]map[string]interface{}, 0)
		restriction := make(map[string]interface{})
		geoRestrictions := make([]map[string]interface{}, 0)
		geoRestriction := make(map[string]interface{})

		geoRestriction["restriction_type"] = d.DistributionConfig.Restrictions.GeoRestriction.RestrictionType
		if (d.DistributionConfig.Restrictions.GeoRestriction.Locations) != nil {
			geoRestriction["locations"] = d.DistributionConfig.Restrictions.GeoRestriction.Locations
		}
		geoRestrictions = append(geoRestrictions, geoRestriction)
		restriction["geo_restriction"] = geoRestrictions

		restrictions = append(restrictions, restriction)
		if len(restrictions) > 0 {
			cf.Restrictions = restrictions
		}
	}
	if d.DistributionConfig.CacheBehaviors != nil {
		orderedCacheBehaviors := make([]map[string]interface{}, 0)
		for _, cacheBehaviour := range d.DistributionConfig.CacheBehaviors {
			orderedCacheBehavior := make(map[string]interface{})
			orderedCacheBehavior["viewer_protocol_policy"] = cacheBehaviour.ViewerProtocolPolicy
			orderedCacheBehaviors = append(orderedCacheBehaviors, orderedCacheBehavior)
		}
		if len(orderedCacheBehaviors) > 0 {
			cf.OrderedCacheBehavior = orderedCacheBehaviors
		}
	}
	if d.DistributionConfig.Logging != nil {
		loggingConfigs := make([]interface{}, 0)
		loggingConfigs = append(loggingConfigs, d.DistributionConfig.Logging)
		if len(loggingConfigs) > 0 {
			cf.LoggingConfig = loggingConfigs
		}
	}
	if d.DistributionConfig.ViewerCertificate != nil {
		viewerCertificates := make([]map[string]interface{}, 0)
		viewerCertificate := make(map[string]interface{})
		viewerCertificate["cloudfront_default_certificate"] = d.DistributionConfig.ViewerCertificate.CloudFrontDefaultCertificate
		viewerCertificate["minimum_protocol_version"] = d.DistributionConfig.ViewerCertificate.MinimumProtocolVersion
		viewerCertificates = append(viewerCertificates, viewerCertificate)
		if len(viewerCertificates) > 0 {
			cf.ViewerCertificate = viewerCertificates
		}
	}
	cf.WebACLId = functions.GetVal(d.DistributionConfig.WebACLId)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: d.AWSCloudFormationMetadata,
	}}
}

func checkDistributionConfig(distributionConfig *cloudfront.Distribution_DistributionConfig) bool {
	return distributionConfig != nil &&
		distributionConfig.Restrictions != nil &&
		distributionConfig.Restrictions.GeoRestriction != nil &&
		len(distributionConfig.Restrictions.GeoRestriction.RestrictionType) > 0
}
