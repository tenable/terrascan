package accurics

uniformBucketEnabled[api.id]
{
     api := input.google_storage_bucket[_]
     not api.config.uniform_bucket_level_access == true  
}