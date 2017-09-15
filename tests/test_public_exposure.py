import unittest
import os
import terraform_validate
from . import settings


class TestPublicExposure(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_alb_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_alb').property('internal').should_not_equal(False)

    def test_aws_db_instance_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_db_instance').property(
            'publicly_accessible').should_not_equal(True)

    def test_aws_dms_replication_instance_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_dms_replication_instance').property(
            'publicly_accessible').should_not_equal(True)

    def test_aws_elb_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_elb').property('internal').should_not_equal(False)

    def test_aws_instance_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_instance').property(
            'associate_public_ip_address').should_not_equal(True)

    def test_aws_launch_configuration_public(self):
        # Assert public ALB
        self.v.resources(
            'aws_launch_configuration').property(
            'associate_public_ip_address').should_not_equal(True)


if __name__ == '__main__':
    unittest.main()
