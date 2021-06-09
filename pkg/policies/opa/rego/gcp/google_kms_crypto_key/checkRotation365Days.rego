package accurics

checkRotation365Days[crypto_key.id]{
    crypto_key := input.google_kms_crypto_key[_]
    to_number(trim(crypto_key.config.rotation_period, "s")) > 31536000
}