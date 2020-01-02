=========
Terrascan
=========

A collection of security and best practice tests for static code analysis of terraform_ templates.

.. _terraform: https://www.terraform.io

* GitHub Repo: https://github.com/cesar-rodriguez/terrascan
* Documentation: https://terrascan.readthedocs.io.
* Free software: GNU General Public License v3

--------
Features
--------
Terrascan will perform tests on your terraform templates to ensure:

- **Encryption**
    - Server Side Encryption (SSE) enabled
    - Use of AWS Key Management Service (KMS) with Customer Managed Keys (CMK)
    - Use of SSL/TLS and proper configuration
- **Security Groups**
    - Provisioning SGs in EC2-classic
    - Ingress open to 0.0.0.0/0
- **Public Exposure**
    - Services with public exposure other than Gateways (NAT, VGW, IGW)
- **Logging & Monitoring**
    - Access logs enabled to resources that support it

----------
Installing
----------
Terrascan uses Python and depends on terraform-validate and pyhcl. After installing python in your system you can follow these steps:

    $ pip install terrascan


-----------------
Running the tests
-----------------
To run execute terrascan.py as follows replacing with the location of your terraform templates:

    $ terrascan --location tests/infrastructure/success --tests all

To run a specific test run the following command replacing encryption with the name of the test to run:

    $ terrascan --location tests/infrastructure/success --tests encryption

To learn more about the options to the cli execute the following:

    $ terrascan -h
