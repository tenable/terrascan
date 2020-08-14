package accurics

privateRepoEnabled[api.id]
{
    api := input.github_repository[_]
    not api.config.private == true
}