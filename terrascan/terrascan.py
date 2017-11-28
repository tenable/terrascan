# -*- coding: utf-8 -*-

"""Console script for terrascan."""

import argparse
import unittest
import sys
from os import path

from terrascan.checks.security_group import TestSecurityGroups
from terrascan.checks.encryption import TestEncryption
from terrascan.checks.logging_and_monitoring import TestLoggingAndMonitoring
from terrascan.checks.public_exposure import TestPublicExposure


test_to_class = {
    'encryption': TestEncryption,
    'logging_and_monitoring': TestLoggingAndMonitoring,
    'public_exposure': TestPublicExposure,
    'security_group': TestSecurityGroups
}


def run_test(args):
    """
    Executes template checks based on cli options
    """
    # Gets absolute location path
    location = path.abspath(args.location[0])
    if not path.exists(location):
        raise Exception("The specified location doesn't exists")

    # Generating list of tests to run
    if args.tests[0] == 'all':
        tests_to_run = [
            'encryption',
            'logging_and_monitoring',
            'public_exposure',
            'security_group']
    else:
        tests_to_run = args.tests[0].split(',')

    # Executing tests
    exit_status = True
    for test_type in tests_to_run:
        print('\n\nRunning {} Tests'.format(test_type))
        test = test_to_class[test_type]
        test.TERRAFORM_LOCATION = location
        runner = unittest.TextTestRunner()
        itersuite = unittest.TestLoader().loadTestsFromTestCase(test)
        result = runner.run(itersuite)
        exit_status = exit_status and not result.wasSuccessful()
    sys.exit(exit_status)


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
(e.g. encryption,security_group) Valid values include:encription,
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
    try:
        args.func(args)
    except Exception:
        print("ERROR: The specified location doesn't exists")
        sys.exit(1)
