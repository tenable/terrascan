provider "aws" {
  region = "us-west-2"
}

resource "aws_kms_key" "PtShGgAkk1" {
  description             = "ptshggakk1"
  deletion_window_in_days = 10
}

resource "aws_db_instance" "PtShGgAdi1" {
  allocated_storage                   = 20
  storage_type                        = "gp2"
  engine                              = "mysql"
  engine_version                      = "5.7"
  instance_class                      = "db.t3.micro"
  name                                = "ptshggadi1"
  backup_retention_period             = 30
  iam_database_authentication_enabled = true
  auto_minor_version_upgrade          = true
  username                            = "slaflheafllaflaehf"
  password                            = "something"
  skip_final_snapshot                 = true
}

resource "aws_db_instance" "PtShGgAdi2" {
  allocated_storage       = 20
  storage_type            = "gp2"
  engine                  = "mysql"
  engine_version          = "5.7"
  instance_class          = "db.t2.micro"
  name                    = "ptshggadi2"
  backup_retention_period = 0
  username                = "slaflheafllaflaehf"
  password                = "something"
  skip_final_snapshot     = true
  multi_az = false
}

resource "aws_db_instance" "PtShGgAdi3" {
  allocated_storage                   = 20
  storage_type                        = "gp2"
  engine                              = "mysql"
  engine_version                      = "5.7"
  instance_class                      = "db.t3.micro"
  name                                = "ptshggadi3"
  backup_retention_period             = 30
  iam_database_authentication_enabled = false
  auto_minor_version_upgrade          = false
  publicly_accessible                 = true
  username                            = "slaflheafllaflaehf"
  password                            = "something"
  skip_final_snapshot                 = true
}

resource "aws_db_instance" "PtShGgAdi4" {
  allocated_storage       = 20
  storage_type            = "gp2"
  engine                  = "mysql"
  engine_version          = "5.7"
  instance_class          = "db.t2.micro"
  name                    = "ptshggadi4"
  backup_retention_period = 0
  ca_cert_identifier      = "rds-ca-2019"
  username                = "slaflheafllaflaehf"
  password                = "something"
  skip_final_snapshot     = true
}

resource "aws_db_instance" "PtShGgAdi5" {
  allocated_storage     = 20
  max_allocated_storage = 50
  storage_type          = "gp2"
  engine                = "mysql"
  engine_version        = "5.7"
  instance_class        = "db.t2.micro"
  name                  = "ptshggadi5"
  username              = "slaflheafllaflaehf"
  password              = "something"
  skip_final_snapshot   = true
}

resource "aws_db_instance" "PtShGgAdi6" {
  allocated_storage   = 20
  storage_type        = "gp2"
  engine              = "mysql"
  engine_version      = "5.7"
  instance_class      = "db.t3.micro"
  name                = "ptshggadi6"
  storage_encrypted   = false
  username            = "slaflheafllaflaehf"
  password            = "something"
  skip_final_snapshot = true
}

resource "aws_db_instance" "PtShGgAdi7" {
  allocated_storage                   = 120
  storage_type                        = "io1"
  iops                                = 2400
  engine                              = "mysql"
  engine_version                      = "5.7"
  instance_class                      = "db.t3.micro"
  name                                = "ptshggadi7"
  storage_encrypted                   = true
  kms_key_id                          = aws_kms_key.PtShGgAkk1.arn
  auto_minor_version_upgrade          = true
  backup_retention_period             = 30
  ca_cert_identifier                  = "rds-ca-2019"
  iam_database_authentication_enabled = true
  publicly_accessible                 = false
  username                            = "slaflheafllaflaehf"
  password                            = "something"
  skip_final_snapshot                 = true
}

resource "aws_db_instance" "PtShGgAdi8" {
  allocated_storage                   = 20
  storage_type                        = "gp2"
  engine                              = "mysql"
  engine_version                      = "5.7"
  instance_class                      = "db.t3.micro"
  name                                = "ptshggadi8"
  storage_encrypted                   = true
  kms_key_id                          = aws_kms_key.PtShGgAkk1.arn
  auto_minor_version_upgrade          = true
  backup_retention_period             = 30
  ca_cert_identifier                  = "rds-ca-2019"
  iam_database_authentication_enabled = true
  publicly_accessible                 = false
  username                            = "awsuser"
  password                            = "something"
  skip_final_snapshot                 = true
}

resource "aws_db_instance" "PtShGgAdi9" {
  allocated_storage                   = 20
  storage_type                        = "gp2"
  engine                              = "mysql"
  engine_version                      = "5.7"
  instance_class                      = "db.t3.micro"
  name                                = "ptshggadi9"
  storage_encrypted                   = true
  kms_key_id                          = aws_kms_key.PtShGgAkk1.arn
  auto_minor_version_upgrade          = true
  backup_retention_period             = 30
  ca_cert_identifier                  = "rds-ca-2019"
  iam_database_authentication_enabled = true
  publicly_accessible                 = false
  username                            = "somelahfaflahf"
  password                            = "something"
  skip_final_snapshot                 = true
}
