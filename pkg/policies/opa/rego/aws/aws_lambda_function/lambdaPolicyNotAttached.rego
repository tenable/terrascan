package accurics

{{.prefix}}lambdaPolicyNotAttached[lambda_function.id] {
    lambda_function :=  input.aws_lambda_function[_]
    object.get(input, "aws_iam_role_policy_attachment", "undefined") == "undefined"
}