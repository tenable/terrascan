# -*- coding: utf-8 -*-
"""Allow terrascan to be executable through `python -m terrascan`."""
from __future__ import absolute_import

from .terrascan import main

if __name__ == "__main__":  # pragma: no cover
    main()
