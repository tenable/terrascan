package accurics

{{.prefix}}neptuneClusterNotEncrypted[np.id] {
    np = input.aws_neptune_cluster[_]
    object.get(np.config, "storage_encrypted", "undefined") == [false, "undefined"][_]
}