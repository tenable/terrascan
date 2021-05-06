package accurics

{{.prefix}}backupRetentionRDS[rds.id]{
    rds = input.aws_rds_cluster[_]
	object.get(rds.config, "backup_retention_period", "undefined") == "undefined"
}

{{.prefix}}backupRetention[rds.id]{
    rds = input.aws_rds_cluster[_]
	rds.config.backup_retention_period <= 7
}