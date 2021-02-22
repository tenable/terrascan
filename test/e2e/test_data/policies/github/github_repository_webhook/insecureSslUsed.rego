package accurics

insecureSslUsed[api.id] {
    api := input.github_repository_webhook[_]
    api.config.configuration[_].insecure_ssl == true
}
