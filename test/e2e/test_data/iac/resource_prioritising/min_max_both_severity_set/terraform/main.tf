resource "aws_db_instance" "PtShGgAdi4" {
  #ts:minseverity=High
  #ts:maxseverity=Low
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