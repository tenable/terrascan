package accurics

{{.prefix}}redshiftPublicAccess[retVal]{
    redshift = input.aws_redshift_cluster[_]
    redshift.config.publicly_accessible == true
    traverse = "publicly_accessible"
    retVal := { "Id": redshift.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "publicly_accessible", "AttributeDataType": "bool", "Expected": false, "Actual": redshift.config.publicly_accessible }
}

{{.prefix}}redshiftPublicAccess[retVal] {
    redshift = input.aws_redshift_cluster[_]
    not redshift.config.publicly_accessible
    traverse = "publicly_accessible"
    retVal := { "Id": redshift.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "publicly_accessible", "AttributeDataType": "bool", "Expected": false, "Actual": true }
}