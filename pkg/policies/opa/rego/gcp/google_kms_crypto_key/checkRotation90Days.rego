package accurics

checkRotation90Days[crypto_key.id]{
    crypto_key := input.google_kms_crypto_key[_]
    to_number(trim(crypto_key.config.rotation_period, "s")) > 7776000
}