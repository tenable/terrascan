package accurics

checkAuditLoggingConfig[iam_audit.id] {
  iam_audit := input.google_project_iam_audit_config[_]
  iam_audit.config.service != "allServices"
} {
  iam_audit := input.google_project_iam_audit_config[_]
  count(iam_audit.config.audit_log_config) < 3
  audit_log_config := iam_audit.config.audit_log_config[_]
  not check_log_type_value(audit_log_config)
  count(audit_log_config.exempted_members) != 0
}

check_log_type_value(item) {
    item.log_type == "ADMIN_READ"
}

check_log_type_value(item) {
    item.log_type == "DATA_READ"
}

check_log_type_value(item) {
    item.log_type == "DATA_WRITE"
}