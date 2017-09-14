import unittest
import os
import terraform_validate
from . import settings


class TestSgsOutsideOfVpc(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_db_security_group_used(self):
        # Assert resource has a KMS with CMK
        self.assertEqual(self.v.resources(
            'aws_db_security_group').resource_list, [])

    def test_aws_redshift_security_group_used(self):
        # Assert resource has a KMS with CMK
        self.assertEqual(self.v.resources(
            'aws_redshift_security_group').resource_list, [])

    def test_aws_elasticache_security_group_used(self):
        # Assert resource has a KMS with CMK
        self.assertEqual(self.v.resources(
            'aws_elasticache_security_group').resource_list, [])


if __name__ == '__main__':
    unittest.main()
