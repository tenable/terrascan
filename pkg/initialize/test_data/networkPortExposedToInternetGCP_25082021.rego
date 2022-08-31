package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
    rule := input.google_compute_firewall[_]
    config := rule.config
    config.direction == "INGRESS"
    config.source_ranges[_] == "0.0.0.0/0"
    fire_rule := config.allow[_]
    fire_rule.protocol == "{{.protocol}}"
    fire_rule.ports[_] == "{{.portNumber}}"

    expected := [ item | item := validate_source(config.source_ranges[_]) ]
	traverse := "source_ranges"

     retVal := {
       "Id": rule.id,
       "ReplaceType": "edit",
       "CodeType": "attribute",
       "Traverse": traverse,
       "Attribute": traverse,
       "AttributeDataType": "list",
       "Expected": expected,
       "Actual": config.source_ranges
     }
}

validate_source(source) = value {
	source == "0.0.0.0/0"
    value := "<cidr>"
}
validate_source(source) = value {
	source != "0.0.0.0/0"
    value := source
}