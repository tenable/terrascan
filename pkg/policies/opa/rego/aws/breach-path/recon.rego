package accurics

recon[ins.id] {
    ins = input.aws_instance[_]
    sec_groups := [ sg | sg := split(ins.config.vpc_security_group_ids[_], ".")[1] ]
    sec_group := sec_groups[_]
    checkSecurityGroupWideOpen(sec_group)

    sub = ins.config.subnet_id
    route_table_association = input.aws_route_table_association[_]
    route_table_association.config.subnet_id == sub
    route_table := split(route_table_association.config.route_table_id, ".")[1]
    checkRouteInternet(ins, route_table)
}


checkSecurityGroupWideOpen(sgName) {
	security_group := input.aws_security_group[_]
    sgName == security_group.name

    some i
    ingress = security_group.config.ingress[i]

    # Checks if the cidr block is not a private IP
    ingress.cidr_blocks[_] == "0.0.0.0/0"

    ports_open = (ingress.to_port - ingress.from_port)
    ports_open > 0
}

checkRouteInternet(instance, arg) {
    rt = input.aws_route_table[_]
    rt.name == arg
    routes = rt.config.route[_]
    routes.gateway_id != ""
}

# checkIMDv1(arg) {
#   not arg.config.metadata_options
# }

# checkIMDv1(arg) {
#   value := arg.config.metadata_options[_]
#   not value.http_endpoint == "disabled"
#   not value.http_tokens == "required"
# }

# checkEC2AllPermissions(arg) {
# 	ins_profile_name := split(arg.config.iam_instance_profile, ".")[1]
#     iam_instance_profile := input.aws_iam_instance_profile[_]
#     ins_profile_name == iam_instance_profile.name

#     role_name := split(iam_instance_profile.config.role, ".")[1]
#     role_policy_attachment := input.aws_iam_role_policy_attachment[_]
#     role_name == split(role_policy_attachment.config.role, ".")[1]
#     policy_name := split(role_policy_attachment.config.policy_arn, ".")[1]

#     iam_policy := input.aws_iam_policy[_]
#     policy_name == iam_policy.name
#     policy := json_unmarshal(iam_policy.config.policy)
#     statement = policy.Statement[_]
#     ac := statement.Action[_]
#     action := split(ac, ":")[0]
#     policyCheck(statement, "ec2:*", "Allow", "*")
# }

# json_unmarshal(s) = result {
# 	s != null
# 	result := json.unmarshal(s)
# }

# # policyCheck(s, a, e ,r) {
# #     split(s.Action[_], ":")[1] == a
# #     s.Effect == e
# #     s.Resource == r
# # }

# policyCheck(s, a, e ,r) {
#     s.Action[_] == a
#     s.Effect == e
#     s.Resource == r
# }

# checkTrustedGroups(instance, security_group) = instances[_].id {
# 	instances := {in | in := input.aws_instance[_]; in.id != instance.id}
#     security_groups := [ sg | sg := split(instances[_].config.vpc_security_group_ids[_], ".")[1]]
#     security_groups[_] == security_group
# }

# checkTrustedGroups(instance, security_group) = instances[_].id {
# 	instances := {in | in := input.aws_instance[_]; in.id != instance.id}
#     security_groups := [ sg | sg := split(instances[_].config.vpc_security_group_ids[_], ".")[1]]

#     sec_group2 := security_groups[_]
#     checkSecurityGroupIdIncluded(sec_group2)
# }

# checkSecurityGroupIdIncluded(argument) {
# 	security_group := input.aws_security_group[_]
#     argument == security_group.name

#     some i
#     ingress = security_group.config.ingress[i]

#     # Checks if the cidr block is not a private IP
#     ingress.cidr_blocks[_] == "0.0.0.0/0"

#     ports_open = (ingress.to_port - ingress.from_port)
#     ports_open > 0
# }
