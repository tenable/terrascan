import unittest
import os
import terraform_validate
from . import settings


class TestAlbListener(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_alb_listener_port(self):
        # Assert that listener port is 443
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_alb_listener').property('port').should_equal('443')

    def test_alb_listener_protocol(self):
        # Assert that protocol is not http
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_alb_listener').property('protocol').should_not_equal('http')
        self.v.resources(
            'aws_alb_listener').property('protocol').should_not_equal('HTTP')

    def test_alb_listener_ssl_policy(self):
        # Assert that old ssl policies are not used
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_alb_listener').property(
            'ssl_policy').should_not_equal('ELBSecurityPolicy-2015-05')
        self.v.resources(
            'aws_alb_listener').property(
            'ssl_policy').should_not_equal('ELBSecurityPolicy-TLS-1-0-2015-04')

    def test_alb_listener_certificate(self):
        # Assert that certificate_arn is set
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_alb_listener').should_have_properties(['certificate_arn'])


class TestAMI(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_ami_ebs_block_device(self):
        # Assert that all resources of type 'ebs_block_device' that are
        # inside a 'aws_ami' are encrypted
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ami').property(
            'ebs_block_device').property('encrypted').should_equal(True)

    def test_ami_ebs_block_device_kms(self):
        # Assert that all resources of type 'ebs_block_device' that are
        # inside a 'aws_ami' have a KMS key
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ami').property(
            'ebs_block_device').should_have_properties(['kms_key_id'])


class TestAMICopy(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_ami_copy(self):
        # Assert that all resources of type 'aws_ami_copy' are encrypted
        # Fail tests if the property does not exist
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ami_copy').property('encrypted').should_equal(True)

    def test_aws_ami_copy_kms(self):
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ami_copy').should_have_properties(['kms_key_id'])


class TestAPIGatewayDomainName(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_ami_copy_kms(self):
        # Assert that certificate settings have been configured
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_api_gateway_domain_name').should_have_properties(
            [
                'certificate_name',
                'certificate_body',
                'certificate_chain',
                'certificate_private_key'
            ])


class TestEC2Instance(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_ec2_instance_ebs_block_device(self):
        # Assert that all resources of type 'ebs_block_device' that are
        # inside a 'aws_instance' are encrypted
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_instance').property(
            'ebs_block_device').property('encrypted').should_equal(True)


class TestCloudTrail(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_cloudtrail_kms(self):
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_cloudtrail').should_have_properties(['kms_key_id'])


class TestCodeBuild(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_cloudtrail_kms(self):
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_codebuild_project').should_have_properties(['encryption_key'])


class TestCodePipeline(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_cloudtrail_kms(self):
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_codepipeline').should_have_properties(['encryption_key'])


class TestEBS(unittest.TestCase):

    def setUp(self):
        # Tell the module where to find your terraform configuration folder
        self.path = os.path.join(
            os.path.dirname(
                os.path.realpath(__file__)), settings.TERRAFORM_LOCATION)
        self.v = terraform_validate.Validator(self.path)

    def test_aws_ebs_volume(self):
        # Assert that all resources of type 'aws_ebs_volume' are encrypted
        # Fail tests if the property does not exist
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ebs_volume').property('encrypted').should_equal(True)

    def test_aws_ebs_volume_kms(self):
        # Assert that a KMS key has been provided
        self.v.error_if_property_missing()
        self.v.enable_variable_expansion()
        self.v.resources(
            'aws_ebs_volume').should_have_properties(['kms_key_id'])


if __name__ == '__main__':
    unittest.main()
