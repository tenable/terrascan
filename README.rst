=========
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
Terrascan uses Python and depends on pyhcl and terraform-validate (a fork has 
been included as part of terrascan that supports terraform 0.12+). 
After installing python in your system you can follow these steps:

    $ pip install terrascan

-----------------
Running the tests
-----------------
To run, execute terrascan.py as follows replacing with the location of your terraform templates:

    $ terrascan --location tests/infrastructure/success --vars tests/infrastructure/vars.json

- **Returns 0 if no failures or errors; 4 otherwise**
	- helps with use in a delivery pipeline
	
- **Parameters**::

	-h, --help            show this help message and exit
	-l LOCATION, --location LOCATION
	                      location of terraform templates to scan
	-v [VARS [VARS ...]], --vars [VARS [VARS ...]]
	                      variables json or .tf file name
	-o OVERRIDES, --overrides OVERRIDES
	                      override rules file name
	-r RESULTS, --results RESULTS
	                      output results file name
	-d [DISPLAYRULES], --displayRules [DISPLAYRULES]
	                      display the rules used
	-w [WARRANTY], --warranty [WARRANTY]
	                      displays the warranty
	-g [GPL], --gpl [GPL]
	                      displays license information
	-c CONFIG, --config CONFIG
	                      logging configuration: error, warning, info, debug, or
	                      none; default is error
- **Override file example**

1. The first attribute is the name of the rule to be overridden.
2. The second attribute is the name of the resource to be overridden.
3. The third atttribute is the RR or RAR number that waives the failure.  
This is required for high severity rules; can be an empty string for medium and low severity rules.

.. code:: json

	{
		"overrides": [
			[
				"aws_s3_bucket_server_side_encryption_configuration",
				"noEncryptionWaived",
				"RR-1234"
			],
			[
				"aws_rds_cluster_encryption",
				"rds_cluster_bad",
				"RAR-98765"
			]
		]
	}

- **Example output**::

	Logging level set to error.
	................
	----------------------------------------------------------------------
	Ran 16 tests in 0.015s

	OK

	Processed 19 files in C:\DEV\terraforms\backends\10-network-analytics


	Results (took 1.08 seconds):

	Failures: (2)
	[high] [aws_dynamodb_table.encryption.server_side_encryption.enabled] should be 'True'. Is: 'False' in module 10-network-analytics, file C:\DEV\terraforms\backends\10-network-analytics\main.tf
	[high] [aws_s3_bucket.noEncryption] should have property: 'server_side_encryption_configuration' in module 10-network-analytics, file C:\DEV\terraforms\backends\10-network-analytics\main.tf

	Errors: (0)

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
 aws_dynamodb_table                        `:heavy_check_mark:`                            
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
 aws_s3_bucket                             `:heavy_check_mark:`                            `:heavy_check_mark:`    `:heavy_check_mark:`
 aws_s3_bucket_object                      `:heavy_check_mark:`
 aws_security_group                                                 `:heavy_check_mark:`   `:heavy_check_mark:`
 aws_security_group_rule                                            `:heavy_check_mark:`   `:heavy_check_mark:`
 aws_ses_receipt_rule                      `:heavy_minus_sign:`
 aws_sqs_queue                             `:heavy_check_mark:`
 aws_ssm_maintenance_window_task                                                                                   `:heavy_check_mark:`
 aws_ssm_parameter                         `:heavy_check_mark:`
========================================  ======================  ======================  ======================  ======================


