=========
Terrascan
=========


.. image:: https://img.shields.io/pypi/v/terrascan.svg
        :target: https://pypi.python.org/pypi/terrascan

.. image:: https://img.shields.io/travis/cesar-rodriguez/terrascan.svg
        :target: https://travis-ci.org/cesar-rodriguez/terrascan

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

--------------------
Using as pre-commit
--------------------
Terrascan can be used on pre-commit hooks to prevent accidental introduction of security weaknesses into your repository. 
This requires having pre-commit_ installed. An example configuration is provided in the comments of the here_ file in this repository.

.. _pre-commit: https://pre-commit.com/
.. _here: .pre-commit-config.yaml

--------------
Feature Status
--------------
Legend:
    - `:heavy_minus_sign:` = test needs to be implemented
    - `:heavy_check_mark:` = test implemented
    - **blank** - N/A

========================================  ======================  ======================  ======================  ======================
 Terraform resources                       Encryption              Security Groups         Public exposure         Logging & Monitoring
========================================  ======================  ======================  ======================  ======================
 aws_alb                                                                                   `:heavy_check_mark:`    `:heavy_check_mark:`
 aws_alb_listener                          `:heavy_check_mark:`
 aws_ami                                   `:heavy_check_mark:`
 aws_ami_copy                              `:heavy_check_mark:`
 aws_api_gateway_domain_name               `:heavy_check_mark:`
 aws_cloudfront_distribution               `:heavy_check_mark:`                                                    `:heavy_check_mark:`
 aws_cloudtrail                            `:heavy_check_mark:`                                                    `:heavy_check_mark:`
 aws_codebuild_project                     `:heavy_check_mark:`
 aws_codepipeline                          `:heavy_check_mark:`
 aws_db_instance                           `:heavy_check_mark:`                            `:heavy_check_mark:`
 aws_db_security_group                                             `:heavy_check_mark:`
 aws_dms_endpoint                          `:heavy_check_mark:`
 aws_dms_replication_instance              `:heavy_check_mark:`                            `:heavy_check_mark:`
 aws_ebs_volume                            `:heavy_check_mark:`
 aws_efs_file_system                       `:heavy_check_mark:`
 aws_elasticache_security_group                                    `:heavy_check_mark:`
 aws_efs_file_system                       `:heavy_check_mark:`
 aws_elasticache_security_group                                    `:heavy_check_mark:`
 aws_elastictranscoder_pipeline            `:heavy_check_mark:`
 aws_elb                                   `:heavy_check_mark:`                            `:heavy_check_mark:`    `:heavy_check_mark:`
 aws_emr_cluster                                                                                                   `:heavy_check_mark:`
 aws_instance                              `:heavy_check_mark:`                            `:heavy_check_mark:`
 aws_kinesis_firehose_delivery_stream      `:heavy_check_mark:`                                                    `:heavy_check_mark:`
 aws_lambda_function                       `:heavy_check_mark:`
 aws_launch_configuration                                                                                          `:heavy_check_mark:`
 aws_lb_ssl_negotiation_policy             `:heavy_minus_sign:`
 aws_load_balancer_backend_server_policy   `:heavy_minus_sign:`
 aws_load_balancer_listener_policy         `:heavy_minus_sign:`
 aws_load_balancer_policy                  `:heavy_minus_sign:`
 aws_opsworks_application                  `:heavy_check_mark:`                            `:heavy_minus_sign:`
 aws_opsworks_custom_layer                                                                 `:heavy_minus_sign:`
 aws_opsworks_ganglia_layer                                                                `:heavy_minus_sign:`
 aws_opsworks_haproxy_layer                                                                `:heavy_minus_sign:`
 aws_opsworks_instance                                                                     `:heavy_minus_sign:`
 aws_opsworks_java_app_layer                                                               `:heavy_minus_sign:`
 aws_opsworks_memcached_layer                                                              `:heavy_minus_sign:`
 aws_opsworks_mysql_layer                                                                  `:heavy_minus_sign:`
 aws_opsworks_nodejs_app_layer                                                             `:heavy_minus_sign:`
 aws_opsworks_php_app_layer                                                                `:heavy_minus_sign:`
 aws_opsworks_rails_app_layer                                                              `:heavy_minus_sign:`
 aws_opsworks_static_web_layer                                                             `:heavy_minus_sign:`
 aws_rds_cluster                           `:heavy_check_mark:`
 aws_rds_cluster_instance                                                                  `:heavy_check_mark:`
 aws_redshift_cluster                      `:heavy_check_mark:`                            `:heavy_check_mark:`    `:heavy_check_mark:`
 aws_redshift_parameter_group              `:heavy_minus_sign:`                                                    `:heavy_minus_sign:`
 aws_redshift_security_group                                        `:heavy_check_mark:`
 aws_s3_bucket                                                                             `:heavy_check_mark:`    `:heavy_check_mark:`
 aws_s3_bucket_object                      `:heavy_check_mark:`
 aws_security_group                                                 `:heavy_check_mark:`
 aws_security_group_rule                                            `:heavy_check_mark:`
 aws_ses_receipt_rule                      `:heavy_minus_sign:`
 aws_sqs_queue                             `:heavy_check_mark:`
 aws_ssm_maintenance_window_task                                                                                   `:heavy_check_mark:`
 aws_ssm_parameter                         `:heavy_check_mark:`
========================================  ======================  ======================  ======================  ======================


