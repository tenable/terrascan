package accurics

iamServiceAccountUsed[api.id]
{
    api := input.google_project_iam_member[_]
    api.config.role == "roles/iam.serviceAccountUser"
}

iamServiceAccountUsed[api.id]
{
    api := input.google_project_iam_member[_]
    api.config.role == "roles/iam.serviceAccountTokenCreator"
}