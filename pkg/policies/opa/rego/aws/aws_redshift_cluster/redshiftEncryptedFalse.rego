package accurics

{{.prefix}}redshiftEncryptedFalse[retVal]{
     redshift_cluster = input.aws_redshift_cluster[_]
     redshift_cluster.config.encrypted == false
     traverse = "encrypted"
     retVal := { "Id": redshift_cluster.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": redshift_cluster.config.encrypted }
}

{{.prefix}}redshiftEncryptedFalse[retVal]{
     redshift_cluster = input.aws_redshift_cluster[_]
     not redshift_cluster.config.encrypted
     traverse = "encrypted"
     retVal := { "Id": redshift_cluster.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": false }
}