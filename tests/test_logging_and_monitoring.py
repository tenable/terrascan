import unittest
import os
import terraform_validate
from . import settings


class TestLoggingAndMonitoring(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
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


if __name__ == '__main__':
    unittest.main()
