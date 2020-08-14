package accurics

checkPubliclyAccessible[api.id]
{
     api := input.google_storage_bucket_iam_binding[_]
     data := api.config.members[_]
     contains(data,"allUsers")     
}

checkPubliclyAccessible[api.id]
{
     api := input.google_storage_bucket_iam_binding[_]
     data := api.config.members[_]
     contains(data,"allAuthenticatedUsers")     
}
