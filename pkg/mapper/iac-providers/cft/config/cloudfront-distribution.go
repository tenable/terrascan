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

package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/cloudfront"
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
func GetCloudFrontDistributionConfig(d *cloudfront.Distribution) []AWSResourceConfig {
	cf := CloudFrontDistributionConfig{
		Config: Config{
			Tags: d.Tags,
		},
	}
	if d.DistributionConfig != nil &&
		d.DistributionConfig.Restrictions != nil &&
		d.DistributionConfig.Restrictions.GeoRestriction != nil &&
		len(d.DistributionConfig.Restrictions.GeoRestriction.RestrictionType) > 0 {
		restrictions := make([]map[string]interface{}, 0)
		restriction := make(map[string]interface{})
		geoRestrictions := make([]map[string]interface{}, 0)
		geoRestriction := make(map[string]interface{})
		geoRestriction["restriction_type"] = d.DistributionConfig.Restrictions.GeoRestriction.RestrictionType
		if len(d.DistributionConfig.Restrictions.GeoRestriction.Locations) > 0 {
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
		for i := range d.DistributionConfig.CacheBehaviors {
			orderedCacheBehavior := make(map[string]interface{})
			orderedCacheBehavior["viewer_protocol_policy"] = d.DistributionConfig.CacheBehaviors[i].ViewerProtocolPolicy
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
	cf.WebACLId = d.DistributionConfig.WebACLId

	return []AWSResourceConfig{{Resource: cf}}
}
