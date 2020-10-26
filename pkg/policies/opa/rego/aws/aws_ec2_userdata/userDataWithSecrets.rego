package accurics


EC2withSecrets[retVal] {
    pattern := ["[A-Za-z0-9/+=]{40}","(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}","(\"|')?(AWS|aws|Aws)?_?(SECRET|secret|Secret)?_?(ACCESS|access|Access)?_?(KEY|key|Key)(\"|')?\\s*(:|=>|=)\\s*(\"|')?[A-Za-z0-9/\\+=]{40}(\"|')?"]
	some i
    instance := input.aws_instance[_]
    user_data := instance.config.user_data
    regex.match(pattern[i],user_data)
    retVal := { "Id": instance.id, "ReplaceType": "edit", "CodeType": "block", "Traverse": "", "Attribute": "", "AttributeDataType": "", "Expected": "No AWS Secrets in user data", "Actual": instance.config.user_data }
  
}