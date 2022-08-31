package accurics

serviceAttachmentAcceptAutomatic[saconf.id] {
  saconf := input.google_compute_service_attachment[_]
  saconf.config.connection_preference == "ACCEPT_AUTOMATIC"
}
