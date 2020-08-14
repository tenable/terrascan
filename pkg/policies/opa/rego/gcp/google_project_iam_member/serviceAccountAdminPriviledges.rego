package accurics

serviceAccountAdminPriviledges[api.id]
{
    api := input.google_project_iam_member[_]
    api.config.role == "roles/editor"
    endswith(api.config.member, ".gserviceaccount.com")
}