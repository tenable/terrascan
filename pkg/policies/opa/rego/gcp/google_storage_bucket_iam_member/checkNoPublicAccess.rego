package accurics

checkNoPublicAccess[bucket_iam.id] {
  bucket_iam := input.google_storage_bucket_iam_member[_]
  bucket_iam_members := bucket_iam.config.members
  public_access_users := ["allUsers", "allAuthenticatedUsers"]
  some i, j
  bucket_iam_members[i] == public_access_users[j]
}