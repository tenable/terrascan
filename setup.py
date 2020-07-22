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


A collection of security and best practice tests for static code analysis of terraform_ templates using terraform_validate_.

.. _terraform: https://www.terraform.io
.. _terraform_validate: https://github.com/elmundio87/terraform_validate

* GitHub Repo: https://github.com/cesar-rodriguez/terrascan
* Documentation: https://terrascan.readthedocs.io.
* Free software: GNU General Public License v3
'''

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'pyhcl>=0.4.4',
]

setup(
    name='terrascan',
    version='0.2.2',
    description="Best practices tests for terraform",
    long_description=readme,
    author="Cesar Rodriguez",
    author_email='therasec@gmail.com',
    url='https://github.com/cesar-rodriguez/terrascan',
    download_url='https://github.com/cesar-rodriguez/terrascan' +
    '/archive/v0.2.2.tar.gz',
    packages=find_packages(where='.'),
    entry_points={
        'console_scripts': [
            'terrascan = terrascan.terrascan:main'
        ]
    },
    include_package_data=True,
    license="GNU General Public License v3",
    zip_safe=False,
    keywords='terrascan',
    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: GNU General Public License v3 (GPLv3)',
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
