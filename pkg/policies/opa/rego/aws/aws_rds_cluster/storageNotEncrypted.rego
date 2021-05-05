package accurics

{{.prefix}}storageNotEncrypted[rds.id]{
    rds = input.aws_rds_cluster[_]
	object.get(rds.config, "storage_encrypted", "undefined") == [false, "undefined"][_]
}