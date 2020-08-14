package accurics

checkRotation365Days[kms.id] {
  kms := input.google_kms_crypto_key[_]
  kms.config.rotation_period <= "31536000s"
}
