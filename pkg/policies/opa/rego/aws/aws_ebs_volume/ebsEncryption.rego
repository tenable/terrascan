package accurics

{{.prefix}}{{.name}}[retVal] {
    block := input.aws_ebs_volume[_]
    checkEncryption(block.config) == true

    traverse = "encrypted"
    retVal := { "Id": block.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": false }
}

checkEncryption(c) = true {
    not c.encrypted
}

checkEncryption(c) =true {
    c.encrypted == false
}

checkEncryption(c) =true {
    c.encrypted == true
    not c.kms_key_id
}