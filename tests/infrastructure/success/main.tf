/**
 * Template to validate encryption test
 */

variable "encryption" {
  description = "Set to false to fail boolen based tests"
  default     = "true"
}

variable "kms_key_arn" {
  description = "Dummy KMS key ARN"
  default     = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
}

variable "certificate_arn" {
  description = "Dummy certificate ARN"
  default     = "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"
}

resource "aws_alb_listener" "front_end" {
  port            = "443"
  protocol        = "HTTPS"
  ssl_policy      = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn = "${var.certificate_arn}"
}

resource "aws_ami" "example" {
  ebs_block_device {
    encrypted  = "${var.encryption}"
    kms_key_id = "${var.kms_key_arn}"
  }
}

resource "aws_ami_copy" "example" {
  kms_key_id = "${var.kms_key_arn}"
  encrypted  = "${var.encryption}"
}

resource "aws_api_gateway_domain_name" "example" {
  certificate_name        = "example-api"
  certificate_body        = "${file("${path.module}/example.com/example.crt")}"
  certificate_chain       = "${file("${path.module}/example.com/ca.crt")}"
  certificate_private_key = "${file("${path.module}/example.com/example.key")}"
}

resource "aws_instance" "foo" {
  associate_public_ip_address = false

  ebs_block_device {
    encrypted = "${var.encryption}"
  }
}

resource "aws_cloudfront_distribution" "distribution" {
  origin {
    custom_origin_config {
      origin_protocol_policy = "https-only"
    }
  }

  default_cache_behavior {
    viewer_protocol_policy = "redirect-to-https"
  }

  cache_behavior {
    viewer_protocol_policy = "redirect-to-https"
  }

  logging_config {
    include_cookies = false
    bucket          = "mylogs.s3.amazonaws.com"
    prefix          = "myprefix"
  }
}

resource "aws_cloudtrail" "foo" {
  kms_key_id     = "${var.kms_key_arn}"
  enable_logging = true
}

resource "aws_codebuild_project" "foo" {
  encryption_key = "${var.kms_key_arn}"
}

resource "aws_codepipeline" "foo" {
  encryption_key = "${var.kms_key_arn}"
}

resource "aws_db_instance" "default" {
  kms_key_id          = "${var.kms_key_arn}"
  storage_encrypted   = "${var.encryption}"
  publicly_accessible = false
}

resource "aws_dms_endpoint" "test" {
  certificate_arn = "${var.certificate_arn}"
  kms_key_arn     = "${var.kms_key_arn}"
  ssl_mode        = "verify-full"
}

resource "aws_dms_replication_instance" "test" {
  kms_key_arn         = "${var.kms_key_arn}"
  publicly_accessible = false
}

resource "aws_ebs_volume" "foo" {
  kms_key_id = "${var.kms_key_arn}"
  encrypted  = "${var.encryption}"
}

resource "aws_efs_file_system" "foo" {
  kms_key_id = "${var.kms_key_arn}"
  encrypted  = "${var.encryption}"
}

resource "aws_elastictranscoder_pipeline" "bar" {
  aws_kms_key_arn = "${var.kms_key_arn}"
}

resource "aws_elb" "foo" {
  internal = true

  listener {
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "${var.certificate_arn}"
  }

  access_logs {
    bucket        = "foo"
    bucket_prefix = "bar"
    interval      = 60
  }
}

resource "aws_kinesis_firehose_delivery_stream" "foo" {
  s3_configuration {
    kms_key_arn = "${var.kms_key_arn}"

    cloudwatch_logging_options {
      enabled = true
    }
  }

  redshift_configuration {
    cloudwatch_logging_options {
      enabled = true
    }
  }

  elasticsearch_configuration {
    cloudwatch_logging_options {
      enabled = true
    }
  }

  extended_s3_configuration {
    kms_key_arn = "${var.kms_key_arn}"
  }
}

resource "aws_lambda_function" "foo" {
  kms_key_arn = "${var.kms_key_arn}"
}

resource "aws_opsworks_application" "foo-app" {
  enable_ssl = true

  ssl_configuration = {
    private_key = "${file("./foobar.key")}"
    certificate = "${file("./foobar.crt")}"
  }
}

resource "aws_rds_cluster" "default" {
  storage_encrypted = "${var.encryption}"
  kms_key_id        = "${var.kms_key_arn}"
}

resource "aws_redshift_cluster" "default" {
  encrypted           = "${var.encryption}"
  kms_key_id          = "${var.kms_key_arn}"
  publicly_accessible = false
  enable_logging      = true
}

resource "aws_s3_bucket_object" "examplebucket_object" {
  server_side_encryption = "aws:kms"
  kms_key_id             = "${var.kms_key_arn}"
}

resource "aws_sqs_queue" "terraform_queue" {
  kms_master_key_id                 = "alias/aws/sqs"
  kms_data_key_reuse_period_seconds = 300
}

resource "aws_ssm_parameter" "secret" {
  type   = "SecureString"
  key_id = "${var.kms_key_arn}"
}

resource "aws_security_group_rule" "allow_all" {
  type            = "ingress"
  from_port       = 0
  to_port         = 65535
  protocol        = "tcp"
  cidr_blocks     = ["192.168.1.1/32"]
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
    cidr_blocks = ["192.168.1.1/32"]
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
  internal = true

  access_logs {
    bucket = "${aws_s3_bucket.alb_logs.bucket}"
    prefix = "test-alb"
  }
}

resource "aws_launch_configuration" "as_conf" {
  associate_public_ip_address = false
}

resource "aws_rds_cluster_instance" "cluster_instances" {
  publicly_accessible = false
}

resource "aws_s3_bucket" "b" {
  acl = "private"

  logging {
    target_bucket = "${aws_s3_bucket.log_bucket.id}"
    target_prefix = "log/"
  }
}

resource "aws_emr_cluster" "emr-test-cluster" {
  log_uri = "s3bucket/test"
}

resource "aws_ssm_maintenance_window_task" "task" {
  logging_info {
    s3_bucket_name = "${aws_s3_bucket.log_bucket.id}"
    s3_region      = "us-east-1"
  }
}
