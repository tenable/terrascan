# -*- coding: utf-8 -*-

"""Tests for security group configuration in terraform templates"""

import unittest
import os
import terraform_validate


class TestSecurityGroups(unittest.TestCase):

    # Set this before running the Test Case
    TERRAFORM_LOCATION = ''

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                # os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
                os.path.realpath(__file__)), self.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_db_security_group_used(self):
        # This SG type exists outside of VPC (e.g. ec2 classic)
        self.assertEqual(self.v.resources(
            'aws_db_security_group').resource_list, [])

    def test_aws_redshift_security_group_used(self):
        # This SG type exists outside of VPC (e.g. ec2 classic)
        self.assertEqual(self.v.resources(
            'aws_redshift_security_group').resource_list, [])

    def test_aws_elasticache_security_group_used(self):
        # This SG type exists outside of VPC (e.g. ec2 classic)
        self.assertEqual(self.v.resources(
            'aws_elasticache_security_group').resource_list, [])

    def test_aws_security_group_rule_open(self):
        # Assert that ingress rule is open to 0.0.0.0/0
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_security_group_rule').with_property(
            'type', 'ingress').property(
            'cidr_blocks').list_should_not_contain('0.0.0.0/0')

    def test_aws_security_group_inline_rule_open(self):
        # Assert that SG has ingress rule open to 0.0.0.0/0
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_security_group').property(
            'ingress').property(
            'cidr_blocks').list_should_not_contain('0.0.0.0/0')
