package accurics

encryptedwithCsek[api.id]
{
     api := input.google_compute_disk[_]
     not api.config.disk_encryption_key == null
}
