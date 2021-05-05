package accurics

{{.prefix}}notEncryptedSecrets[secrets_mgm.id] {
	secrets_mgm := input.aws_secretsmanager_secret[_]
    object.get(secrets_mgm.config, "kms_key_id", "undefined") == [null, "undefined"][_]
}