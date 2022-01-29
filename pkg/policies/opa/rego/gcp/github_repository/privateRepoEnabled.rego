package accurics

privateRepoEnabled[api.id] {
    api := input.github_repository[_]
    not api.config.private == true
    not api.config.visibility == "private"
    not api.config.visibility == "internal"
}
