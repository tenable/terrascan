package accurics

{{.prefix}}configEnabledForAllRegions[retVal]{
    con = input.aws_config_configuration_aggregator[_]
    some i
    ag_source = con.config.account_aggregation_source[i]
    # need some logic to guess ReplaceType as add / edit, we get this value in both cases
    ag_source.all_regions == false
    traverse = sprintf("account_aggregation_source[%d].all_regions", [i])
    retVal := { "Id": con.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "account_aggregation_source.all_regions", "AttributeDataType": "boolean", "Expected": true, "Actual": ag_source.all_regions }
}

{{.prefix}}configEnabledForAllRegions[retVal]{
    con = input.aws_config_configuration_aggregator[_]
    some i
    ag_source = con.config.organization_aggregation_source[i]
    # need some logic to guess ReplaceType as add / edit, we get this value in both cases
    ag_source.all_regions == false
    traverse = sprintf("organization_aggregation_source[%d].all_regions", [i])
    retVal := { "Id": con.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "organization_aggregation_source.all_regions", "AttributeDataType": "boolean", "Expected": true, "Actual": ag_source.all_regions }
}