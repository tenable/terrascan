package accurics

checkVersioningEnabled[log_object.id] {
  log_object := input.google_storage_bucket[_]
  versioning := log_object.config.versioning[_]
  versioning.enabled == false
}

checkLoggingEnabled[log_object.id] {
  log_object := input.google_storage_bucket[_]
  count(log_object.config.logging) <= 0
}