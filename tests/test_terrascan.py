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

def test_success():
    """
    Test successful terraform templates
    """
    with pytest.raises(SystemExit) as pytest_wrapped_e:
        terrascan.main([
            '-l',
            'tests/infrastructure/success',
        ])
    assert pytest_wrapped_e.value.code == 0
    
def test_fail():
    """
    Test successful terraform templates
    """
    with pytest.raises(SystemExit) as pytest_wrapped_e:
        terrascan.main([
            '-l',
            'tests/infrastructure/fail',
        ])
    assert pytest_wrapped_e.value.code == 4