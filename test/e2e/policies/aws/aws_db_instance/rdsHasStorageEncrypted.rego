package accurics

rdsHasStorageEncrypted[rds.id] {
    rds := input.aws_db_instance[_]
    rds.config.storage_encrypted == null
}

rdsHasStorageEncrypted[rds.id] {
    rds := input.aws_db_instance[_]
    rds.config.storage_encrypted == false
}

rdsHasStorageEncrypted[rds.id] {
    rds := input.aws_db_instance[_]
    not rds.config.kms_key_id
}

rdsHasStorageEncrypted[rds.id] {
    rds := input.aws_db_instance[_]
    rds.config.kms_key_id == null
}