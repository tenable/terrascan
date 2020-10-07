package accurics

#{{.prefix}}lambdaNotEncryptedWithKms[retVal] {
  #lambda := input.aws_lambda_function[_]
  #lambda.config.kms_key_arn == null
  #traverse = ""
#  retVal := { "Id": lambda.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_key_arn", "AttributeDataType": "string", "Expected": "<kms_key_arn>", "Actual": null }
#}

lambdaNotEncryptedWithKms[lambda.id] {
  lambda := input.aws_lambda_function[_]
  not lambda.config.kms_key_arn
}
