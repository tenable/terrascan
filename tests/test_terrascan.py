#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Tests for `terrascan` package."""


from unittest import TestCase
from terrascan import terrascan


class CliTestCase(TestCase):

    def test_with_empty_args(self):
        """
        User passes no args, should fail with SystemExit
        """
        with self.assertRaises(SystemExit):
            terrascan.main(['-h'])

    def test_success(self):
        """
        Test successful terraform templates
        """
        with self.assertRaises(SystemExit):
            terrascan.main([
                '-l',
                'tests/infrastructure/success',
            ])

    def test_fail(self):
        """
        Test successful terraform templates
        """
        with self.assertRaises(SystemExit):
            terrascan.main([
                '-l',
                'tests/infrastructure/fail',
            ])
