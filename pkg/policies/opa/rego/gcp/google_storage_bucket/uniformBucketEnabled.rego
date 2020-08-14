package accurics

uniformBucketEnabled[api.id]
{
     api := input.google_storage_bucket[_]
     not api.config.bucket_policy_only == true  
}