package accurics

rdsHasStorageEncrypted[data.id]{
    data := input.aws_db_instance[_]
    data.config.storage_encrypted == null
}

rdsHasStorageEncrypted[data.id]{
    data := input.aws_db_instance[_]
    data.config.storage_encrypted == false
}

rdsHasStorageEncrypted[data.id]{
    data := input.aws_db_instance[_]
    not data.config.kms_key_id
}

rdsHasStorageEncrypted[data.id]{
    data := input.aws_db_instance[_]
    data.config.kms_key_id == null
}