package accurics

{{.prefix}}{{.name}}[retVal]{
     redshift_cluster = input.aws_redshift_cluster[_]
     redshift_cluster.config.encrypted == true
     not redshift_cluster.config.kms_key_id
     traverse = "kms_key_id"
     retVal := { "Id": redshift_cluster.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_key_id", "AttributeDataType": "string", "Expected": "<kms_key_id>", "Actual": null }
}