package accurics

{{.prefix}}amiSharedToMultipleAccounts[retVal] {
	all_permissions := input.aws_ami_launch_permission
	launch_permission := input.aws_ami_launch_permission[_]

	image_id = launch_permission.config.image_id
    account_id = launch_permission.config.account_id

    accounts := [ account | all_permissions[i].config.image_id == image_id; account := all_permissions[i].config.account_id ]
	count(accounts) > 1
    account_id != accounts[0]

    retVal := { "Id": launch_permission.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}
