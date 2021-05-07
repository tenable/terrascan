package accurics

{{.prefix}}configEnabledForAllRegions[con.id]{
    con = input.aws_config_configuration_aggregator[_]
    ag_source = con.config.account_aggregation_source[_]
    object.get(ag_source, "all_regions", "undefined") == ["undefined", false][_]
}


{{.prefix}}configEnabledForAllRegions[con.id]{
    con = input.aws_config_configuration_aggregator[_]
    ag_source = con.config.organization_aggregation_source[_]
    object.get(ag_source, "all_regions", "undefined") == ["undefined", false][_]
}