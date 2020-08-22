package accurics

checkPubliclyAccessible[api.id] {
     api := input.google_storage_bucket_iam_binding[_]
     member := api.config.members[_]
     contains(member, ["allUsers", "allAuthenticatedUsers"][_])     
}