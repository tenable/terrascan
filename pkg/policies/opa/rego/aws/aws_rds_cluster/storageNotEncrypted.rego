package accurics

{{.prefix}}storageNotEncrypted[retVal]{
    rds = input.aws_rds_cluster[_]
	not Encrypted(rds.config)
    traverse = "storage_encrypted"
	retVal := { "Id": rds.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "storage_encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": rds.config.storage_encrypted }
}

{{.prefix}}storageNotEncrypted[retVal]{
    rds = input.aws_db_instance[_]
	not Encrypted(rds.config)
    traverse = "storage_encrypted"
	retVal := { "Id": rds.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "storage_encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": rds.config.storage_encrypted }
}

Encrypted(config) = true  {
	config.storage_encrypted != null
	config.storage_encrypted == true
}

Encrypted(config) = true  {
	config.storage_encrypted == null
	config.kms_key_id != ""
}