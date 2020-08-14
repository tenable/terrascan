package accurics

{{.prefix}}noRoute53RecordSet[retVal] {
    route := input.aws_route53_record[_]
    check_empty_records(route.config.records)
    traverse = "records"
    retVal := { "Id": route.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "records", "AttributeDataType": "list", "Expected": ["<record>"], "Actual": null }
}

check_empty_records(records) = true {
	records == null
}

check_empty_records(records) = true {
	records != null
	count(records) <= 0
}