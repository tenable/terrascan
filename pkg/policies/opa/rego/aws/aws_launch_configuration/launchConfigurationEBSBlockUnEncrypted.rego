package accurics

{{.prefix}}launchConfigurationEBSBlockUnEncrypted[result.retVal] {
    block := input.aws_launch_configuration[_]
    result := checkEncryption(block.id, block.config)
    result != null
}

checkEncryption(id, c) = { "retVal": retVal } {
	some i
    ebsBlock := c.ebs_block_device[i]
    ebsBlock.encrypted != null
    ebsBlock.encrypted == false
    traverse := sprintf("ebs_block_device[%d].encrypted", [i])
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ebs_block_device.encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": ebsBlock.encrypted }
}

checkEncryption(id, c) = { "retVal": retVal } {
	some i
    ebsBlock := c.ebs_block_device[i]
	ebsBlock.encrypted == null
	traverse := sprintf("ebs_block_device[%d].encrypted", [i])
	retVal := { "Id": id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ebs_block_device.encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": null }
}