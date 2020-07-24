#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""The setup script."""

from setuptools import setup, find_packages

readme = '''=========
Terrascan
=========

.. image:: https://img.shields.io/pypi/v/terrascan.svg
        :target: https://pypi.python.org/pypi/terrascan
        :alt: pypi

.. image:: https://img.shields.io/travis/cesar-rodriguez/terrascan.svg
        :target: https://travis-ci.org/cesar-rodriguez/terrascan
        :alt: build

.. image:: https://readthedocs.org/projects/terrascan/badge/?version=latest
        :target: https://terrascan.readthedocs.io/en/latest/?badge=latest
        :alt: Documentation Status

.. image:: https://pyup.io/repos/github/cesar-rodriguez/terrascan/shield.svg
     :target: https://pyup.io/repos/github/cesar-rodriguez/terrascan/
     :alt: Updates


A linter for security best practices testing of Terraform_ templates.

.. _Terraform: https://www.terraform.io

* GitHub Repo: https://github.com/cesar-rodriguez/terrascan
* Documentation: https://terrascan.readthedocs.io.
* Free software: Apache-2.0
'''

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'pyhcl>=0.4.4',
]

setup(
    name='Terrascan',
    version='0.2.3',
    description="Security best practice static code analysis for terraform",
    long_description=readme,
    author="Accurics",
    author_email='support@accurics.com',
    url='https://github.com/accurics/terrascan',
    download_url='https://github.com/accurics/terrascan' +
    '/archive/v0.2.3.tar.gz',
    packages=find_packages(where='.'),
    entry_points={
        'console_scripts': [
            'terrascan = terrascan.terrascan:main'
        ]
    },
    include_package_data=True,
    license="Apache-2.0",
    zip_safe=False,
    keywords='terrascan',
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License',
        'Natural Language :: English',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
    ],
    test_suite='tests',
    tests_require=requirements,
    setup_requires=requirements,
    install_requires=requirements,
)
