package accurics

checkVM_NoFullCloudAccess[log_object.id] {
  log_object := input.google_compute_instance[_]
  service_account := log_object.config.service_account[_]
  scope := service_account.scopes[_]
  contains(scope, "cloud-platform")
}
