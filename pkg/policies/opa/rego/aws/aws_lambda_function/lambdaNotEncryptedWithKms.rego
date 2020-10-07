package accurics

lambdaNotEncryptedWithKms[lambda.id] {
  lambda := input.aws_lambda_function[_]
  not lambda.config.kms_key_arn
}
