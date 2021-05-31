package accurics

lambdaNotEncryptedWithKms[lambda.id] {
	lambda := input.aws_lambda_function[_]

	object.get(lambda.config, "environment", "undefined") != "undefined"
	object.get(lambda.config.environment[_], "variables", "undefined") != "undefined"
	lambda.config.environment[_].variables != {}

	object.get(lambda.config, "kms_key_arn", "undefined") == "undefined"
}
