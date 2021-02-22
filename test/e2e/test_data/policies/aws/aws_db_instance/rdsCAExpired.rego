package accurics

{{.prefix}}rdsCAExpired[retVal] {
    rds = input.aws_db_instance[_]
    rds.config.ca_cert_identifier != "rds-ca-2019"
    traverse = "ca_cert_identifier"
    retVal := { "Id": rds.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ca_cert_identifier", "AttributeDataType": "string", "Expected": "rds-ca-2019", "Actual": rds.config.ca_cert_identifier}
}