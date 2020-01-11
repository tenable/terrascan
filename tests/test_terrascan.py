#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Tests for `terrascan` package."""

import pytest
from terrascan import terrascan


def test_with_empty_args():
    """
    User passes no args, should fail with SystemExit
    """
    with pytest.raises(SystemExit):
        terrascan.main(['-h'])
