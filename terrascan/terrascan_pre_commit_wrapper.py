from __future__ import absolute_import
from __future__ import print_function
from __future__ import unicode_literals

import io
import sys

from terrascan import terrascan


def main(argv=None):
    argv = argv if argv is not None else sys.argv[1:]
    print(argv)
    terrascan.main(['-l', '.'])


if __name__ == '__main__':
    main()