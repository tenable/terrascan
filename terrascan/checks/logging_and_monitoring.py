# -*- coding: utf-8 -*-

"""Tests for logging and monitoring configuration in terraform templates"""

import unittest
import os
import terraform_validate


class TestLoggingAndMonitoring(unittest.TestCase):

    # Set this before running the Test Case
    TERRAFORM_LOCATION = ''

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), self.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_alb_logging(self):
        self.v.resources(
            'aws_alb').should_have_properties(['access_logs'])

    def test_aws_cloudfront_distribution_logging(self):
        self.v.resources(
            'aws_cloudfront_distribution').should_have_properties(
            ['logging_config'])

    def test_aws_cloudtrail_logging(self):
        self.v.resources(
            'aws_cloudtrail').property(
            'enable_logging').should_not_equal(False)

    def test_aws_elb_logging(self):
        self.v.resources(
            'aws_elb').should_have_properties(
            ['access_logs'])

    def test_aws_emr_cluster_logging(self):
        self.v.resources(
            'aws_emr_cluster').should_have_properties(
            ['log_uri'])

    def test_aws_kinesis_firehose_delivery_stream__s3_config_logging(self):
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            's3_configuration').should_have_properties(
            ['cloudwatch_logging_options'])
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            's3_configuration').property(
            'cloudwatch_logging_options').property(
            'enabled').should_equal(True)

    def test_aws_kinesis_firehose_delivery_stream_redshift_conf_logging(self):
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            'redshift_configuration').should_have_properties(
            ['cloudwatch_logging_options'])
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            'redshift_configuration').property(
            'cloudwatch_logging_options').property(
            'enabled').should_equal(True)

    def test_aws_kinesis_firehose_delivery_stream__es_config_logging(self):
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            'elasticsearch_configuration').should_have_properties(
            ['cloudwatch_logging_options'])
        self.v.resources(
            'aws_kinesis_firehose_delivery_stream').property(
            'elasticsearch_configuration').property(
            'cloudwatch_logging_options').property(
            'enabled').should_equal(True)

    def test_aws_redshift_cluster_logging(self):
        self.v.enable_variable_expansion()
        self.v.error_if_property_missing()
        self.v.resources(
            'aws_redshift_cluster').property(
            'enable_logging').should_not_equal(False)

    def test_aws_s3_bucket_logging(self):
        self.v.resources(
            'aws_s3_bucket').should_have_properties(
            ['logging'])

    def test_aws_ssm_maintenance_window_task_logging(self):
        self.v.resources(
            'aws_ssm_maintenance_window_task').should_have_properties(
            ['logging_info'])
