package accurics

{{.prefix}}launchConfigurationRootBlockUnEncrypted[result.retVal] {
    block := input.aws_launch_configuration[_]
    result := checkEncryption(block.id, block.config)
    result != null
}

checkEncryption(id, c) = { "retVal": retVal } {
	some i
    rootBlock := c.root_block_device[i]
	rootBlock.encrypted != null
    rootBlock.encrypted == false
    traverse := sprintf("root_block_device[%d].encrypted", [i])
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "root_block_device.encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": rootBlock.encrypted }
}

checkEncryption(id, c) = { "retVal": retVal } {
	some i
    rootBlock := c.root_block_device[i]
	rootBlock.encrypted == null
	traverse := sprintf("root_block_device[%d].encrypted", [i])
	retVal := { "Id": id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "root_block_device.encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": null }
}