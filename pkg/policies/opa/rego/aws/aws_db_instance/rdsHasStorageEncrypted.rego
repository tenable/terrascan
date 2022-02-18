package accurics

rdsHasStorageEncrypted[rds.id] {
	rds := input.aws_db_instance[_]
	encryptionCheck(rds.config)
}

encryptionCheck(rds_config) {
	object.get(rds_config, "storage_encrypted", "undefined") == [[], null, "undefined"]
}

encryptionCheck(rds_config) {
	rds_config.storage_encrypted != true
}
