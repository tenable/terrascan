package accurics

lambdaNotInVpc[lambda.id] {
  lambda := input.aws_lambda_function[_]
  not lambda.config.vpc_config
}
