package accurics

checkOSLoginEnabled[metadata.id] {
  metadata := input.google_compute_project_metadata[_]
  metadata.config.metadata == null
} {
  metadata := input.google_compute_project_metadata[_]
  metadata.config.metadata != null
  not metadata.config.metadata["enable-oslogin-2fa"]
} {
  metadata := input.google_compute_project_metadata[_]
  metadata.config.metadata != null
  metadata.config.metadata["enable-oslogin-2fa"] != "TRUE"
}
