package accurics

athenaQueryEncryption[api.id]{
    api := input.aws_athena[_]
    data := api.config.configuration[_]
    resConfig := data.result_configuration[_]
    encOpt := resConfig.encryption_configuration[_]
    not encOpt.encryption_option == "SSE_KMS"
    not encOpt.encryption_option == "CSE_KMS"
    not encOpt.encryption_option == "SSE_S3"
}