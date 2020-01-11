#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""The setup script."""

from setuptools import setup, find_packages

with open('README.rst') as readme_file:
    readme = readme_file.read()

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'pyhcl==0.4.0',
]

setup(
    name='terrascan',
    version='0.1.2',
    description="Best practices tests for terraform",
    long_description=readme,
    author="Cesar Rodriguez",
    author_email='therasec@gmail.com',
    url='https://github.com/cesar-rodriguez/terrascan',
    download_url='https://github.com/cesar-rodriguez/terrascan' +
    '/archive/v0.1.2.tar.gz',
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
    dependency_links=[
        "git+git://github.com/cesar-rodriguez/terraform_validate.git@master#egg=terraform-validate"
    ],
    tests_require=requirements,
    setup_requires=requirements,
    install_requires=requirements,
)
