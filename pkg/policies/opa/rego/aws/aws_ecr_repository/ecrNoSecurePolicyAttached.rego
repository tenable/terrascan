package accurics

{{.prefix}}ecrNoSecurePolicyAttached[ecr_repo.id] {
	ecr_repo := input.aws_ecr_repository[_]
    object.get(input, "aws_ecr_repository_policy", "undefined") == "undefined"
}