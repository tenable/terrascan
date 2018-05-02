/**
 * Template to validate encryption test
 */

variable "encryption" {
  description = "Set to false to fail boolen based tests"
  default     = "false"
}

resource "aws_alb_listener" "front_end" {
  port       = "80"
  protocol   = "http"
  ssl_policy = "ELBSecurityPolicy-2015-05"

  # certificate_arn   = "arn:aws:iam::187416307283:server-certificate/test_cert_rab3wuqwgja25ct3n4jdj2tzu4"
}

resource "aws_ami" "example" {
  ebs_block_device {
    encrypted = "${var.encryption}"

    # Comment the line below to fail KMS test
    # kms_key_id = "1234"
  }
}

resource "aws_ami_copy" "example" {
  # Comment the line below to fail KMS test
  # kms_key_id = "1234"
  encrypted = "${var.encryption}"
}

resource "aws_api_gateway_domain_name" "example" {
  # Comment the lines below to fail certificate test  # certificate_name        = "example-api"  # certificate_body        = "${file("${path.module}/example.com/example.crt")}"  # certificate_chain       = "${file("${path.module}/example.com/ca.crt")}"  # certificate_private_key = "${file("${path.module}/example.com/example.key")}"
}

resource "aws_instance" "foo" {
  associate_public_ip_address = true

  ebs_block_device {
    encrypted = "${var.encryption}"
  }
}

resource "aws_cloudfront_distribution" "distribution" {
  origin {
    custom_origin_config {
      origin_protocol_policy = "http-only"
    }
  }

  default_cache_behavior {
    viewer_protocol_policy = "allow-all"
  }

  cache_behavior {
    viewer_protocol_policy = "allow-all"
  }
}

resource "aws_cloudtrail" "foo" {
  # Comment the line below to fail KMS test  # kms_key_id = "1234"
  enable_logging = false
}

resource "aws_codebuild_project" "foo" {
  # Comment the line below to fail KMS test  # encryption_key = "1234"
}

resource "aws_codepipeline" "foo" {
  # Comment the line below to fail KMS test  # encryption_key = "1234"
}

resource "aws_db_instance" "default" {
  # Comment the line below to fail KMS test
  # kms_key_id = "1234"
  storage_encrypted = "${var.encryption}"
}

resource "aws_dms_endpoint" "test" {
  # certificate_arn = "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"
  # kms_key_arn     = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
  ssl_mode = "none"
}

resource "aws_dms_replication_instance" "test" {
  # kms_key_arn = "arn:aws:kms:us-east-1:123456789012:key/12345 78-1234-1234-1234-123456789012"
  publicly_accessible = true
}

resource "aws_ebs_volume" "foo" {
  # Comment the line below to fail KMS test
  # kms_key_id = "1234"
  encrypted = "${var.encryption}"
}

resource "aws_efs_file_system" "foo" {
  # Comment the line below to fail KMS test
  # kms_key_id = "1234"
  encrypted = "${var.encryption}"
}

resource "aws_elastictranscoder_pipeline" "bar" {
  # aws_kms_key_arn = "${var.kms_key_arn}"
}

resource "aws_elb" "foo" {
  internal = false

  listener {
    lb_port     = 80
    lb_protocol = "http"
  }
}

resource "aws_elb" "bar" {
  listener {
    lb_port     = 21
    lb_protocol = "tcp"
  }
}

resource "aws_elb" "baz" {
  listener {
    lb_port     = 23
    lb_protocol = "tcp"
  }
}

resource "aws_elb" "foobar" {
  listener {
    lb_port     = 5900
    lb_protocol = "tcp"
  }
}

resource "aws_kinesis_firehose_delivery_stream" "foo" {
  s3_configuration {
    # kms_key_arn = "${var.kms_key_arn}"
  }

  extended_s3_configuration {
    # kms_key_arn = "${var.kms_key_arn}"
  }

  redshift_configuration {
    cloudwatch_logging_options {
      enabled = false
    }
  }

  elasticsearch_configuration {}
}

resource "aws_lambda_function" "foo" {
  # kms_key_arn = "${var.kms_key_arn}"
}

resource "aws_opsworks_application" "foo-app" {
  enable_ssl = false
}

resource "aws_rds_cluster" "default" {
  storage_encrypted = false

  # kms_key_id      = "${var.kms_key_arn}"
}

resource "aws_redshift_cluster" "default" {
  publicly_accessible = true
  encrypted           = "${var.encryption}"

  # kms_key_id = "${var.kms_key_arn}"
}

resource "aws_sqs_queue" "terraform_queue" {
  # kms_master_key_id                 = "alias/aws/sqs"  # kms_data_key_reuse_period_seconds = 300
}

resource "aws_ssm_parameter" "secret" {
  type = "String"

  # key_id = "${var.kms_key_arn}"
}

resource "aws_elasticache_security_group" "bar" {}

resource "aws_db_security_group" "default" {}

resource "aws_redshift_security_group" "default" {}

resource "aws_security_group_rule" "allow_all" {
  type            = "ingress"
  from_port       = 0
  to_port         = 65535
  protocol        = "tcp"
  cidr_blocks     = ["0.0.0.0/0"]
  prefix_list_ids = ["pl-12c4e678"]

  security_group_id = "sg-123456"
}

resource "aws_security_group" "allow_all" {
  name        = "allow_all"
  description = "Allow all inbound traffic"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
    prefix_list_ids = ["pl-12c4e678"]
  }
}

resource "aws_alb" "test" {
  internal = false
}

resource "aws_db_instance" "default" {
  publicly_accessible = true
}

resource "aws_launch_configuration" "as_conf" {
  associate_public_ip_address = true
}

resource "aws_rds_cluster_instance" "cluster_instances" {
  publicly_accessible = true
}

resource "aws_ssm_maintenance_window_task" "task" {}
