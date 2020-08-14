package accurics

{{.prefix}}redshiftAuditLogs[retVal]{
    redshift = input.aws_redshift_cluster[_]
    redshift.config.logging == []
    rc = "ewogICJsb2dnaW5nIjogewogICAgImVuYWJsZSI6IHRydWUsCiAgICAiYnVja2V0X25hbWUiOiAiPGJ1Y2tldF9uYW1lPiIsCiAgICAiczNfa2V5X3ByZWZpeCI6ICI8czNfa2V5X3ByZWZpeD4iCiAgfQp9"
    traverse = ""
    retVal := { "Id": redshift.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}

{{.prefix}}redshiftAuditLogs[retVal]{
    redshift = input.aws_redshift_cluster[_]
    redshift.config.logging[_].enable == false
    rc = "ewogICJlbmFibGUiOiB0cnVlLAogICJidWNrZXRfbmFtZSI6ICI8YnVja2V0X25hbWU+IiwKICAiczNfa2V5X3ByZWZpeCI6ICI8czNfa2V5X3ByZWZpeD4iCn0="
    traverse = "logging"
    retVal := { "Id": redshift.id, "ReplaceType": "edit", "CodeType": "block", "Traverse": traverse, "Attribute": "logging", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}