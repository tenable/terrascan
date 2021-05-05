package accurics

{{.prefix}}docDbEncrypted[doc_cluster.id] {
    doc_cluster := input.aws_docdb_cluster[_]
    object.get(doc_cluster.config, "storage_encrypted", "undefined") == [false, "undefined"][_]
}

{{.prefix}}docDbEncrypted[doc_cluster.id] {
    doc_cluster := input.aws_docdb_cluster[_]
    object.get(doc_cluster.config, "kms_key_id", "undefined") == "undefined"
}