package accurics

{{.prefix}}vpcFlowLogsNotEnabled[retVal] {
    vpc := input.aws_vpc[_]
    vpc_input := input
    vpc.type == "aws_vpc"

	not flowLogExist(vpc, vpc_input)

    rc = "cmVzb3VyY2UgImF3c19mbG93X2xvZyIgIiMjcmVzb3VyY2VfbmFtZSMjIiB7CiAgdnBjX2lkICAgICAgICAgID0gIiR7YXdzX3ZwYy4jI3Jlc291cmNlX25hbWUjIy5pZH0iCiAgaWFtX3JvbGVfYXJuICAgID0gIiMjYXJuOmF3czppYW06OjExMTExMTExMTExMTpyb2xlL3NhbXBsZV9yb2xlIyMiCiAgbG9nX2Rlc3RpbmF0aW9uID0gIiR7YXdzX3MzX2J1Y2tldC4jI3Jlc291cmNlX25hbWUjIy5hcm59IgogIHRyYWZmaWNfdHlwZSAgICA9ICJBTEwiCgogIHRhZ3MgPSB7CiAgICBHZW5lcmF0ZWRCeSA9ICJBY2N1cmljcyIKICAgIFBhcmVudFJlc291cmNlSWQgPSAiYXdzX3ZwYy4jI3Jlc291cmNlX25hbWUjIyIKICB9Cn0KCnJlc291cmNlICJhd3NfczNfYnVja2V0IiAiIyNyZXNvdXJjZV9uYW1lIyMiIHsKICBidWNrZXQgPSAiIyNyZXNvdXJjZV9uYW1lIyNfZmxvd19sb2dfczNfYnVja2V0IgogIGFjbCAgICA9ICJwcml2YXRlIgogIGZvcmNlX2Rlc3Ryb3kgPSB0cnVlCgogIHZlcnNpb25pbmcgewogICAgZW5hYmxlZCA9IHRydWUKICAgIG1mYV9kZWxldGUgPSB0cnVlCiAgfQoKICBzZXJ2ZXJfc2lkZV9lbmNyeXB0aW9uX2NvbmZpZ3VyYXRpb24gewogICAgcnVsZSB7CiAgICAgIGFwcGx5X3NlcnZlcl9zaWRlX2VuY3J5cHRpb25fYnlfZGVmYXVsdCB7CiAgICAgICAgc3NlX2FsZ29yaXRobSA9ICJBRVMyNTYiCiAgICAgIH0KICAgIH0KICB9Cn0KCnJlc291cmNlICJhd3NfczNfYnVja2V0X3BvbGljeSIgIiMjcmVzb3VyY2VfbmFtZSMjIiB7CiAgYnVja2V0ID0gIiR7YXdzX3MzX2J1Y2tldC4jI3Jlc291cmNlX25hbWUjIy5pZH0iCgogIHBvbGljeSA9IDw8UE9MSUNZCnsKICAiVmVyc2lvbiI6ICIyMDEyLTEwLTE3IiwKICAiU3RhdGVtZW50IjogWwogICAgewogICAgICAiU2lkIjogIiMjcmVzb3VyY2VfbmFtZSMjLXJlc3RyaWN0LWFjY2Vzcy10by11c2Vycy1vci1yb2xlcyIsCiAgICAgICJFZmZlY3QiOiAiQWxsb3ciLAogICAgICAiUHJpbmNpcGFsIjogWwogICAgICAgIHsKICAgICAgICAgICJBV1MiOiBbCiAgICAgICAgICAgICJhcm46YXdzOmlhbTo6IyNhY291bnRfaWQjIzpyb2xlLyMjcm9sZV9uYW1lIyMiLAogICAgICAgICAgICAiYXJuOmF3czppYW06OiMjYWNvdW50X2lkIyM6dXNlci8jI3VzZXJfbmFtZSMjIgogICAgICAgICAgXQogICAgICAgIH0KICAgICAgXSwKICAgICAgIkFjdGlvbiI6ICJzMzpHZXRPYmplY3QiLAogICAgICAiUmVzb3VyY2UiOiAiYXJuOmF3czpzMzo6OiR7YXdzX3MzX2J1Y2tldC4jI3Jlc291cmNlX25hbWUjIy5pZH0vKiIKICAgIH0KICBdCn0KUE9MSUNZCn0="
    decode_rc = base64.decode(rc)
    replaced_vpc_id := replace(decode_rc, "##resource_name##", vpc.name)

    traverse = ""
    retVal := { "Id": vpc.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": base64.encode(replaced_vpc_id), "Actual": null }
}

flowLogExist(vpc, vpc_input) = exists {
	flow_log_vpcs_set := { vpc_id | input.aws_flow_log[i].type == "aws_flow_log"; vpc_id := input.aws_flow_log[i].config.vpc_id }
	flow_log_vpcs_set[vpc.id]
    exists = true
} else = exists {
	flow_log_tags_set := { resource_id | input.aws_flow_log[i].type == "aws_flow_log"; resource_id := input.aws_flow_log[i].config.tags.ParentResourceId }
    vpc_name := sprintf("aws_vpc.%s", [vpc.name])
	flow_log_tags_set[vpc_name]
    exists = true
}