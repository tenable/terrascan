#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""The setup script."""

from setuptools import setup, find_packages

with open('README.rst') as readme_file:
    readme = readme_file.read()

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'ply==3.10',
    'pyhcl==0.3.9',
    'terraform-validate==2.5.0',
]

setup_requirements = [
    'ply==3.10',
    'pyhcl==0.3.9',
    'terraform-validate==2.5.0',
]

test_requirements = [
    'ply==3.10',
    'pyhcl==0.3.9',
    'terraform-validate==2.5.0',
]

setup(
    name='terrascan',
    version='0.1.0',
    description="Security and best practices tests for terraform",
    long_description=readme + '\n\n' + history,
    author="Cesar Rodriguez",
    author_email='therasec@gmail.com',
    url='https://github.com/cesar-rodriguez/terrascan',
    packages=find_packages(where='.'),
    entry_points={
        'console_scripts': [
            'terrascan = terrascan.__main__:main',
        ]
    },
    include_package_data=True,
    install_requires=requirements,
    license="GNU General Public License v3",
    zip_safe=False,
    keywords='terrascan',
    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: GNU General Public License v3 (GPLv3)',
        'Natural Language :: English',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
    ],
    test_suite='tests',
    tests_require=test_requirements,
    setup_requires=setup_requirements,
)
