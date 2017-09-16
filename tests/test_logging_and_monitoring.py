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


if __name__ == '__main__':
    unittest.main()
