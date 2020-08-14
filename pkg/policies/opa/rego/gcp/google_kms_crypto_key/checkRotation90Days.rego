package accurics

checkRotation90Days[api.id]
{
    api := input.google_kms_crypto_key[_]
    api.config.rotation_period <= "7776000s"
}