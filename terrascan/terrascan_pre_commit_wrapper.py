from __future__ import absolute_import
from __future__ import print_function
from __future__ import unicode_literals

import argparse
from os import path

from terrascan import terrascan


def main(argv=None):
    parser = argparse.ArgumentParser()
    parser.add_argument(
        'filenames', nargs='*',
        help='Filenames pre-commit believes are changed.',
    )
    parser.add_argument(
        '-l',
        '--location',
        help='Location of terraform templates to scan',
        nargs=1,
        default='.'
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
    
    args = parser.parse_args(argv)
    
    tests = args.tests[0]
    location = path.abspath(args.location[0])
    
    if not path.exists(location):
        print("ERROR: The specified location doesn't exists")
        return 1
    
    return terrascan.run_test(location, tests)


if __name__ == '__main__':
    exit(main())