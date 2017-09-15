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


if __name__ == '__main__':
    unittest.main()
