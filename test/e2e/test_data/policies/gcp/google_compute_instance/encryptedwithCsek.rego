package accurics

encryptedwithCsek[retVal]
{
    api := input.google_compute_disk[_]
    not api.config.disk_encryption_key

    association := input.google_compute_attached_disk[_]
    diskName := split(association.config.disk, ".")[1]

    api.name == diskName
    instanceName := split(association.config.instance, ".")[1]

    instance := input.google_compute_instance[_]
    instanceName == instance.name
    retVal := instance.id
}

encryptedwithCsek[retVal]
{
    api := input.google_compute_disk[_]
    api.config.disk_encryption_key == null

    association := input.google_compute_attached_disk[_]
    diskName := split(association.config.disk, ".")[1]

    api.name == diskName
    instanceName := split(association.config.instance, ".")[1]

    instance := input.google_compute_instance[_]
    instanceName == instance.name
    retVal := instance.id
}