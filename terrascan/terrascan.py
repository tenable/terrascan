# -*- coding: utf-8 -*-
"""
    terrascan: A collection of security and best practice tests for static code analysis of terraform templates using terraform_validate.

    Copyright (C) 2020 Accurics, Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
"""

import argparse
import unittest
import os
import re
import sys
import subprocess
import json
import time
from terrascan.embedded import terraform_validate
import logging

jsonOutput = {
    "dateTimeStamp": "",
    "terrascan-version": "",
    "failures": [],
    "errors": [],
    "files": [],
    "rules": []
}


###############################################################################################################################################################################
# Rules:  these are the rules used to verify that the terraform files are set up correctly.
#
# available methods for resources:
#   property(property_name): returns a list of the requested property_name; if self.v.error_if_property_missing() is called before the rule, will fail if property is missing.
#   with_property(property_name, regex_value): returns a list of the requested property_name with the requested regex_value.
#   should_not_exist(): fails if the resource doesn't exist.
#   should_have_properties(properties_list): fails if any of the properties in the given properties_list doesn't exist.
#   should_not_have_properties(properties_list): fails if any of the properties in the given properties_list exists.
#   find_property(property_name_regex): returns a list of the requested property_name_regex.
# available methods for properties:
#   property(property_name):  returns a list of the requested property_name; if self.v.error_if_property_missing() is called before the rule, will fail if property is missing.
#   should_equal(expected_value): fails if property value doesn't equal given expected_value.
#   should_equal_case_insensitive(expected_value): fails if property value doesn't equal given expected_value ignoring case.
#   should_not_equal(expected_value): fails if property value equals given expected_value.
#   should_not_equal_case_insensitive(expected_value): fails if property value equals given expected_value ignoring case.
#   list_should_contain_any(values_list): fails if the value of the property doesn't contain any of the values in values_list.
#   list_should_contain(values_list): fails if the value of the property doesn't contain all of the values in values_list.
#   list_should_not_contain(values_list): fails if the value of the property contains any of the values in values_list.
#   should_have_properties(properties_list): fails if a property doesn't contain any of the properties in properties_list.
#   should_not_have_properties(properties_list): fails if a property contains any of the properties in properties_list.
#   find_property(property_name_regex): returns a list of the requested property_name_regex.
#   should_match_regex(property_value_regex): fails if the value of the property doesn't match the given property_value_regex.
#   should_contain_valid_json(): fails if the value of the property doesn't contain valid json.
###############################################################################################################################################################################
class Rules(unittest.TestCase):

    rules = []

    def setUp(self):
        self.v = terraform_validate.Validator()
        self.v.preprocessor = self.preprocessor
        self.v.overrides = self.overrides

    #################################################################################################
    # examples of good and bad (marked with ***error***) are given before each rule
    #################################################################################################

    #################################################################################################
    # This resource block creates an S3 bucket with encryption.
    # resource "aws_s3_bucket" "encryptedBucket" {
    #   bucket = "good-bucket-name"
    #   server_side_encryption_configuration {
    #     rule {
    #       apply_server_side_encryption_by_default {
    #         kms_master_key_id = "${data.aws_kms_key.bucket.arn}"
    #         sse_algorithm = "aws:kms"
    #       }
    #     }
    #   }
    # }
    #
    # This resource block creates an S3 bucket with no encryption.  ***error***
    # resource "aws_s3_bucket" "noEncryption" {
    #   bucket = "bad-bucket-name"
    # }
    def test_aws_s3_bucket_server_side_encryption_configuration(self):
        # get name of rule from function name
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that server side encryption is turned on for s3 buckets
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            # to change severity, override it here (default is high)
            validator.severity = "high"
            validator.resources('aws_s3_bucket').should_have_properties(['server_side_encryption_configuration'])

    #################################################################################################
    # This resource block creates a dynamodb table with no encryption.  ***error***
    # resource "aws_dynamodb_table" "noEncryption" {
    #   name = "${local.env}"
    # }
    #
    # This resource block creates a dynamodb table with no encryption.  ***error***
    # resource "aws_dynamodb_table" "encryptionEnabledFalse" {
    #   name = "${local.env}"
    #   server_side_encryption {
    #     enabled = false
    #   }
    # }
    #
    # This resource block creates a dynamodb table with encryption.
    # resource "aws_dynamodb_table" "encryptionEnabledTrue" {
    #   name = "${local.env}"
    #   server_side_encryption {
    #     enabled = true
    #   }
    # }
    def test_aws_dynamodb_table_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that encryption is turned on for dynamodb tables
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dynamodb_table').property('server_side_encryption').property('enabled').should_equal(True)

    #################################################################################################
    # This resource block creates an ebs volue with no encryption.  ***error***
    # resource "aws_ebs_volume" "noEncryption" {
    #   availability_zone = "us-east-1a"
    #   size = 10
    #   type = "gp2"
    #   tags {
    #     Name = "Encryption Test"
    #   }
    #   kms_key_id = "${data.aws_kms_key.volume.arn}"
    # }
    #
    # This resource block creates an ebs volue with no encryption.  ***error***
    # resource "aws_ebs_volume" "encryptionEnabledFalse" {
    #   availability_zone = "us-east-1a"
    #   size = 10
    #   type = "gp2"
    #   tags {
    #     Name = "Encryption Test"
    #   }
    #   encrypted = false
    #   kms_key_id = "${data.aws_kms_key.volume.arn}"
    # }
    #
    # resource "aws_ebs_volume" "encryptionEnabledTrue" {
    #   availability_zone = "us-east-1a"
    #   size = 10
    #   type = "gp2"
    #   tags {
    #     Name = "Encryption Test"
    #   }
    #   # The attributes below mark the volume for encryption.
    #   encrypted = true
    #   kms_key_id = "${data.aws_kms_key.volume.arn}"
    # }
    def test_aws_ebs_volume_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that all resources of type 'aws_ebs_volume' are encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ebs_volume').property('encrypted').should_equal(True)

    #################################################################################################
    # This resource block creates a kms key without key rotation.  ***error***
    # resource "aws_kms_key" "no_enable_key_rotation" {
    #   description             = "KMS key 1"
    #   deletion_window_in_days = 10
    # }
    #
    # This resource block creates a kms key without key rotation.  ***error***
    # resource "aws_kms_key" "enable_key_rotation_false" {
    #   description             = "KMS key 1"
    #   deletion_window_in_days = 10
    #   enable_key_rotation     = false
    # }
    #
    # This resource block creates a kms key with key rotation.
    # resource "aws_kms_key" "enable_key_rotation_true" {
    #   description             = "KMS key 1"
    #   deletion_window_in_days = 10
    #   enable_key_rotation     = true
    # }
    def test_aws_kms_key_rotation(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that all aws_kms_key resources have key rotation enabled
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kms_key').property('enable_key_rotation').should_equal(True)

    #################################################################################################
    # This resource block creates an iam user login profile.  ***error***
    # resource "aws_iam_user_login_profile" "badExample" {
    #   user    = "some_user_name"
    #   pgp_key = "keybase:some_person_that_exists"
    # }
    def test_aws_iam_user_login_profile(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # resource aws_iam_user_login_profile should not exist
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_iam_user_login_profile').should_not_exist()

    #################################################################################################
    # This resource block creates a security group rule.  ***error***
    # resource "aws_security_group_rule" "badExample" {
    #   type              = "ingress"
    #   from_port         = 0
    #   to_port           = 65535
    #   protocol          = "tcp"
    #   cidr_blocks       = ["0.0.0.0/0"]
    #   self              = true
    #   security_group_id = "${aws_security_group.emr-master.id}"
    # }
    def test_aws_security_group_rule_ingress_open(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that security group rule ingress is not open to 0.0.0.0/0
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_security_group_rule').with_property('type', 'ingress').property('cidr_blocks').list_should_not_contain('0.0.0.0/0')

    #################################################################################################
    # This resource block creates a security group.  ***error***
    # resource "aws_security_group" "badExample" {
    #   name        = "generic-emr-master"
    #   description = "Manage traffic for EMR masters"
    #   vpc_id      = "${local.emr_vpc_id}"
    #   ingress {
    #     from_port = 443
    #     to_port = 443
    #     protocol = "tcp"
    #     cidr_blocks = ["0.0.0.0/0"]
    #   }
    # }
    def test_aws_security_group_ingress_open(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # verify that security group ingress is not open to 0.0.0.0/0
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_security_group').property('ingress').property('cidr_blocks').list_should_not_contain('0.0.0.0/0')

    #################################################################################################
    # This resource block creates an aws_db_instance.
    # resource "aws_db_instance" "default" {
    #   allocated_storage    = 20
    #   engine               = "mysql"
    #   instance_class       = "db.t2.micro"
    #   name                 = "mydb"
    #   storage_encrypted    = true
    # }
    #
    # This resource block creates an aws_db_instance.  ***error***
    # resource "aws_db_instance" "default" {
    #   allocated_storage    = 20
    #   engine               = "mysql"
    #   instance_class       = "db.t2.micro"
    #   name                 = "mydb"
    #   storage_encrypted    = false
    # }
    def test_aws_db_instance_encrypted(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that DB is encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_db_instance').property('storage_encrypted').should_equal(True)

    #################################################################################################
    # This resource block creates an aws_rds_cluster.
    # resource "aws_rds_cluster" "default" {
    #   database_name           = "mydb"
    #   master_username         = "foo"
    #   master_password         = "bar"
    #   storage_encrypted       = true
    # }
    #
    # This resource block creates an aws_rds_cluster.  ***error***
    # resource "aws_rds_cluster" "default" {
    #   database_name           = "mydb"
    #   master_username         = "foo"
    #   master_password         = "bar"
    # }
    def test_aws_rds_cluster_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource is encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_rds_cluster').property('storage_encrypted').should_equal(True)

    #################################################################################################
    # public exposure - these were part of the original terrascan
    #################################################################################################

    def test_aws_alb_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb', 'aws_alb']).property('internal').should_not_equal(False)

    def test_aws_db_instance_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_db_instance').property('publicly_accessible').should_not_equal(True)

    def test_aws_dms_replication_instance_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dms_replication_instance').property('publicly_accessible').should_not_equal(True)

    def test_aws_elb_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').property('internal').should_not_equal(False)

    def test_aws_instance_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_instance').property('associate_public_ip_address').should_not_equal(True)

    def test_aws_launch_configuration_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_launch_configuration').property('associate_public_ip_address').should_not_equal(True)

    def test_aws_rds_cluster_instance_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_rds_cluster_instance').property('publicly_accessible').should_not_equal(True)

    def test_aws_redshift_cluster_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_redshift_cluster').property('publicly_accessible').should_not_equal(True)

    def test_aws_s3_bucket_public(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_s3_bucket').property('acl').should_not_equal('public-read')
            validator.resources('aws_s3_bucket').property('acl').should_not_equal('public-read-write')
            validator.resources('aws_s3_bucket').property('acl').should_not_equal('authenticated-read')
            validator.resources('aws_s3_bucket').should_not_have_properties(['website'])

    #################################################################################################
    # other terrascan original rules - prefix test with X to disable
    #################################################################################################
    # encryption
    #################################################################################################

    def test_aws_alb_listener_port(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that listener port is 443
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb_listener', 'aws_alb_listener']).property('port').should_equal('443')

    def test_aws_alb_listener_protocol(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that protocol is not http
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb_listener', 'aws_alb_listener']).property('protocol').should_not_equal_case_insensitive('http')

    def test_aws_alb_listener_ssl_policy(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that old ssl policies are not used
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb_listener', 'aws_alb_listener']).property('ssl_policy').should_not_equal('ELBSecurityPolicy-2015-05')
            validator.resources(['aws_lb_listener', 'aws_alb_listener']).property('ssl_policy').should_not_equal('ELBSecurityPolicy-TLS-1-0-2015-04')

    def test_aws_alb_listener_certificate(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that certificate_arn is set
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb_listener', 'aws_alb_listener']).should_have_properties(['certificate_arn'])

    def test_aws_ami_ebs_block_device_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ami 'ebs_block_device' blocks are encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ami').property('ebs_block_device').property('encrypted').should_equal(True)

    def test_aws_ami_ebs_block_device_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ami 'ebs_block_device' blocks has KMS
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ami').property('ebs_block_device').should_have_properties(['kms_key_id'])

    def test_aws_ami_copy_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resources are encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ami_copy').property('encrypted').should_equal(True)

    def test_aws_ami_copy_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ami_copy').should_have_properties(['kms_key_id'])

    def test_aws_api_gateway_domain_name_certificate(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that certificate settings have been configured
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_api_gateway_domain_name').should_have_properties(
                [
                    'certificate_name',
                    'certificate_body',
                    'certificate_chain',
                    'certificate_private_key'
                ])

    def test_aws_instance_ebs_block_device_encrypted(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ec2 instance 'ebs_block_device' is encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_instance').property('ebs_block_device').property('encrypted').should_equal(True)

    def test_aws_cloudfront_distribution_origin_protocol_policy(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that origin receives https only traffic
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudfront_distribution').property('origin').property('custom_origin_config').property('origin_protocol_policy').should_equal("https-only")

    def test_aws_cloudfront_distribution_def_cache_viewer_prot_policy(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that cache protocol doesn't allow all
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudfront_distribution').property('default_cache_behavior').property('viewer_protocol_policy').should_not_equal("allow-all")

    def test_aws_cloudfront_distribution_cache_beh_viewer_proto_policy(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that cache protocol doesn't allow all
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudfront_distribution').property('cache_behavior').property('viewer_protocol_policy').should_not_equal("allow-all")

    def test_aws_cloudtrail_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudtrail').should_have_properties(['kms_key_id'])

    def test_aws_codebuild_project_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_codebuild_project').should_have_properties(['encryption_key'])

    def test_aws_codepipeline_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_codepipeline').should_have_properties(['encryption_key'])

    def test_aws_db_instance_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_db_instance').should_have_properties(['kms_key_id'])

    def test_aws_dms_endpoint_ssl_mode(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that SSL is verified
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dms_endpoint').property('ssl_mode').should_equal('verify-full')

    def test_aws_dms_endpoint_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dms_endpoint').should_have_properties(
                [
                    'kms_key_arn'
                ])

    def test_aws_dms_endpoint_certificate(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that SSL cert has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dms_endpoint').should_have_properties(
                [
                    'certificate_arn'
                ])

    def test_aws_dms_replication_instance_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_dms_replication_instance').should_have_properties(['kms_key_arn'])

    def test_aws_ebs_volume_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ebs_volume').should_have_properties(['kms_key_id'])

    def test_aws_efs_file_system_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that all resources of type 'aws_efs_file_system' are encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_efs_file_system').property('encrypted').should_equal(True)

    def test_aws_efs_file_system_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_efs_file_system').should_have_properties(['kms_key_id'])

    def test_aws_elastictranscoder_pipeline_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elastictranscoder_pipeline').should_have_properties(['aws_kms_key_arn'])

    def test_aws_elb_listener_port_80(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 80 (http)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').property('listener').property('lb_port').should_not_equal(80)

    def test_aws_elb_listener_port_21(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 21 ftp
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').property('listener').property('lb_port').should_not_equal(21)

    def test_aws_elb_listener_port_23(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 23 telnet
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').property('listener').property('lb_port').should_not_equal(23)

    def test_aws_elb_listener_port_5900(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 5900 VNC
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').property('listener').property('lb_port').should_not_equal(5900)

    def test_aws_kinesis_firehose_delivery_stream_s3_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 80 (http)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kinesis_firehose_delivery_stream').property('s3_configuration').should_have_properties(['kms_key_arn'])

    def test_aws_kinesis_firehose_delivery_stream_extended_s3_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert ELB listener port is not 80 (http)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kinesis_firehose_delivery_stream').property('extended_s3_configuration').should_have_properties(['kms_key_arn'])

    def test_aws_lambda_function_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_lambda_function').should_have_properties(['kms_key_arn'])

    def test_aws_opsworks_application_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource is encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_opsworks_application').property('enable_ssl').should_equal(True)

    def test_aws_rds_cluster_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource has a KMS with CMKs
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_rds_cluster').should_have_properties(['kms_key_id'])

    def test_aws_redshift_cluster_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource is encrypted
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_redshift_cluster').property('encrypted').should_equal(True)

    def test_aws_redshift_cluster_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource has a KMS with CMKs
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_redshift_cluster').should_have_properties(['kms_key_id'])

    def test_aws_s3_bucket_object_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource has a KMS with CMKs
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_s3_bucket_object').should_have_properties(['kms_key_id'])

    def test_aws_sqs_queue_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource has a KMS with CMK
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_sqs_queue').should_have_properties(['kms_master_key_id', 'kms_data_key_reuse_period_seconds'])

    def test_aws_ssm_parameter_encryption(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource is encrypted with KMS
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ssm_parameter').property('type').should_equal("SecureString")

    def test_aws_ssm_parameter_kms(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # Assert resource has a KMS with CMK
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ssm_parameter').should_have_properties(['key_id'])

    #################################################################################################
    # logging and monitoring
    #################################################################################################

    def test_aws_alb_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources(['aws_lb', 'aws_alb']).should_have_properties(['access_logs'])

    def test_aws_cloudfront_distribution_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudfront_distribution').should_have_properties(['logging_config'])

    def test_aws_cloudtrail_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_cloudtrail').property('enable_logging').should_not_equal(False)

    def test_aws_elb_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elb').should_have_properties(['access_logs'])

    def test_aws_emr_cluster_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_emr_cluster').should_have_properties(['log_uri'])

    def test_aws_kinesis_firehose_delivery_stream__s3_config_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kinesis_firehose_delivery_stream').property('s3_configuration').should_have_properties(['cloudwatch_logging_options'])
            validator.resources('aws_kinesis_firehose_delivery_stream').property('s3_configuration').property('cloudwatch_logging_options').property('enabled').should_equal(True)

    def test_aws_kinesis_firehose_delivery_stream_redshift_conf_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kinesis_firehose_delivery_stream').property('redshift_configuration').should_have_properties(['cloudwatch_logging_options'])
            validator.resources('aws_kinesis_firehose_delivery_stream').property('redshift_configuration').property('cloudwatch_logging_options').property('enabled').should_equal(True)

    def test_aws_kinesis_firehose_delivery_stream__es_config_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_kinesis_firehose_delivery_stream').property('elasticsearch_configuration').should_have_properties(['cloudwatch_logging_options'])
            validator.resources('aws_kinesis_firehose_delivery_stream').property('elasticsearch_configuration').property('cloudwatch_logging_options').property('enabled').should_equal(True)

    def test_aws_redshift_cluster_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        self.v.error_if_property_missing()
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_redshift_cluster').property('enable_logging').should_not_equal(False)

    def test_aws_s3_bucket_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_s3_bucket').should_have_properties(['logging'])

    def test_aws_ssm_maintenance_window_task_logging(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_ssm_maintenance_window_task').should_have_properties(['logging_info'])

    #################################################################################################
    # security group
    #################################################################################################

    def test_aws_db_security_group_used(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # This security group type exists outside of VPC (e.g. ec2 classic)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_db_security_group').should_not_exist()

    def test_aws_redshift_security_group_used(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # This security group type exists outside of VPC (e.g. ec2 classic)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_redshift_security_group').should_not_exist()

    def test_aws_elasticache_security_group_used(self):
        ruleName = sys._getframe().f_code.co_name[5:]
        self.rules.append(ruleName)
        # This security group type exists outside of VPC (e.g. ec2 classic)
        validator_generator = self.v.get_terraform_files(self.isRuleOverridden(ruleName))
        for validator in validator_generator:
            validator.resources('aws_elasticache_security_group').should_not_exist()


    def isRuleOverridden(self, ruleName):
        for override in self.overrides:
            if ruleName == override[0]:
                return True
        return False


#################################################################################################
# run the tests
#################################################################################################
def terrascan(args):
    start = time.time()

    try:
        result = subprocess.run(['pip', 'show', 'terrascan-sf'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout = result.stdout.decode("utf-8")
        versionStr = "Version: "
        startIndex = stdout.find(versionStr)
    except:
        startIndex = -1
    if startIndex == -1:
        version = "?"
    else:
        startIndex += len(versionStr)
        endIndex = stdout.find("\r", startIndex)
        version = stdout[startIndex:endIndex]

    # process the arguments
    terraformLocation = args.location[0]
    if not os.path.isabs(terraformLocation):
        terraformLocation = os.path.join(os.sep, os.path.abspath("."),  terraformLocation)
    if args.vars:
        variablesJsonFilename = []
        for fileName in args.vars:
            if not os.path.isabs(fileName):
                fileName = os.path.join(os.sep, os.path.abspath("."),  fileName)
            variablesJsonFilename.append(fileName)
    else:
        variablesJsonFilename = None
    if args.overrides:
        Rules.overrides = []
        overridesFileName = args.overrides[0]
        if not os.path.isabs(overridesFileName):
            overridesFileName = os.path.join(os.sep, os.path.abspath("."),  overridesFileName)
        try:
            with open(overridesFileName, "r", encoding="utf-8") as fp:
                overridesFileString = fp.read()
            overrides = json.loads(overridesFileString)
            overrides = overrides["overrides"]
            # validate overrides
            for override in overrides:
                if len(override) < 2 or len(override) > 3 or len(override) == 3 and not re.match(r"RR-\d{1,10}$|RAR-\d{1,10}$", override[2]):
                    print("***Invalid entry in overrides file:  " + override)
                    print("Needs to be in the following format:  rule_name:resource_name or rule_name:resource_name:RR-xxx or rule_name:resource_name:RAR-xxx where xxx is 1-10 digits")
                    sys.exit(99)
                Rules.overrides.append(override)
        except Exception as e:
            print("***Error loading overrides file " + overridesFileName)
            print(e)
            sys.exit(99)
    else:
        Rules.overrides = ""
    if args.results:
        outputJsonFileName = args.results[0]
        if not os.path.isabs(outputJsonFileName):
            outputJsonFileName = os.path.join(os.sep, os.path.abspath("."),  outputJsonFileName)
    else:
        outputJsonFileName = None
    if args.config:
        config = args.config[0]
    else:
        config = None

    # set logging based on logging.config if present (default is error)
    if config == "none":
        logging.basicConfig(level=logging.CRITICAL)
    elif config == "warning":
        logging.basicConfig(level=logging.WARNING)
    elif config == "info":
        logging.basicConfig(level=logging.INFO)
    elif config == "debug":
        logging.basicConfig(level=logging.DEBUG)
    else:
        config = "error"
        logging.basicConfig(level=logging.ERROR)

    print("terrascan version {0}".format(version))
    print("Logging level set to {0}.".format(config))

    Rules.preprocessor = terraform_validate.PreProcessor(jsonOutput)
    Rules.preprocessor.process(terraformLocation, variablesJsonFilename)

    runner = unittest.TextTestRunner()
    itersuite = unittest.TestLoader().loadTestsFromTestCase(Rules)
    runner.run(itersuite)

    end = time.time()
    elapsedTime = end - start
    processedMessage = "Processed on " + time.strftime("%m/%d/%Y") + " at " + time.strftime("%H:%M")

    if outputJsonFileName:
        jsonOutput["dateTimeStamp"] = processedMessage
        jsonOutput["terrascan-version"] = version
        for fileName in Rules.preprocessor.fileNames:
            jsonOutput["files"].append(fileName)
        for rule in Rules.rules:
            jsonOutput["rules"].append(rule)
        with open(outputJsonFileName, 'w') as jsonOutFile:
            json.dump(jsonOutput, jsonOutFile)

    print("\nProcessed " + str(len(Rules.preprocessor.fileNames)) + " files in " + terraformLocation + "\n")
    for fileName in Rules.preprocessor.fileNames:
        logging.debug("  Processed " + fileName)
    print("")

    print(processedMessage)
    print("Results (took %.2f seconds):" % elapsedTime)
    rc = 0
    print("\nFailures: (" + str(len(jsonOutput["failures"])) + ")")
    for failure in jsonOutput["failures"]:
        m, f = getMF(failure)
        waived = ""
        if len(failure["waived"]) > 0:
            waived =  failure["waived"] + " "
        print("[" + failure["severity"] + "] " + waived + failure["message"] + m + f)
        if failure["waived"] == "":
            rc = 4
    print("\nErrors: (" + str(len(jsonOutput["errors"])) + ")")
    for error in jsonOutput["errors"]:
        m, f = getMF(error)
        print("[" + error["severity"] + "] " + error["message"] + m + f)
        rc = 4
    if args.displayRules:
        print("\nRules used:")
        for rule in Rules.rules:
            print(rule)

    sys.exit(rc)


# Returns command line parser for terrascan
def create_parser():
    parser = argparse.ArgumentParser(description="A collection of security and best practice tests for static code analysis of terraform templates using terraform_validate.")

    parser.add_argument(
        '-l',
        '--location',
        help='location of terraform templates to scan',
        nargs=1,
        required=True
    )
    parser.add_argument(
        '-v',
        '--vars',
        help='variables json or .tf file name',
        nargs='*',
    )
    parser.add_argument(
        '-o',
        '--overrides',
        help='override rules file name',
        nargs=1
    )
    parser.add_argument(
        '-r',
        '--results',
        help='output results file name',
        nargs=1,    )
    parser.add_argument(
        '-d',
        '--displayRules',
        help='display the rules used',
        nargs='?',
        const=True, default=False
    )
    parser.add_argument(
        '-c',
        '--config',
        help='logging configuration:  error, warning, info, debug, or none; default is error',
        nargs=1,    )
    parser.set_defaults(func=terrascan)

    return parser


def getMF(json):
    if json["moduleName"] == "---":
        m = ""
    else:
        m = " in module " + json["moduleName"]
    if json["fileName"] == "---":
        f = ""
    else:
        if m == "":
            f = " in file " + json["fileName"]
        else:
            f = ", file " + json["fileName"]
    return m, f

def main(args=None):
    """
    Terrascan console script. Parses user input to determine location of
    terraform templates and which tests to execute
    """
    parser = create_parser()
    args = parser.parse_args(args)

    #tests = args.tests[0]
    #location = path.abspath(args.location[0])

    #if not path.exists(location):
        #print("ERROR: The specified location doesn't exists")
        #exit(1)

    #exit(run_test(location, tests))
    terrascan(args)
