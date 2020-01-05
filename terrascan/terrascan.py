# -*- coding: utf-8 -*-

"""Console script for terrascan."""

import argparse
import unittest
from os import path

from .checks.security_group import TestSecurityGroups
from .checks.encryption import TestEncryption
from .checks.logging_and_monitoring import TestLoggingAndMonitoring
from .checks.public_exposure import TestPublicExposure


test_to_class = {
    'encryption': TestEncryption,
    'logging_and_monitoring': TestLoggingAndMonitoring,
    'public_exposure': TestPublicExposure,
    'security_group': TestSecurityGroups
}


def run_test(location, tests):
    """
    Executes template checks based on cli options
    """
    # Generating list of tests to run
    if tests == 'all':
        tests_to_run = [
            'encryption',
            'logging_and_monitoring',
            'public_exposure',
            'security_group']
    else:
        tests_to_run = tests.split(',')

    # Executing tests
    exit_status = True
    for test_type in tests_to_run:
        print('\n\nRunning {} Tests'.format(test_type))
        test = test_to_class[test_type]
        test.TERRAFORM_LOCATION = location
        runner = unittest.TextTestRunner()
        itersuite = unittest.TestLoader().loadTestsFromTestCase(test)
        result = runner.run(itersuite)
        exit_status = exit_status and result.wasSuccessful()
    if exit_status:
        return 0
    else:
        return 1


def create_parser():
    """
    Returns command line parser for terrascan
    """
    parser = argparse.ArgumentParser()

    parser.add_argument(
        '-l',
        '--location',
        help='Location of terraform templates to scan',
        nargs=1
    )
    parser.add_argument(
        '-t',
        '--tests',
        help='''Comma separated list of test to run or "all" for all tests
(e.g. encryption,security_group) Valid values include:encryption,
logging_and_monitoring, public_exposure, security_group''',
        nargs=1,
        default=['all']
    )
    parser.set_defaults(func=run_test)

    return parser


def main(args=None):
    """
    Terrascan console script. Parses user input to determine location of
    terraform templates and which tests to execute
    """
    parser = create_parser()
    args = parser.parse_args(args)
    
    tests = args.tests[0]
    location = path.abspath(args.location[0])
    
    if not path.exists(location):
        print("ERROR: The specified location doesn't exists")
        exit(1)
        
    exit(run_test(location, tests))
    