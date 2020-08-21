package accurics

{{.prefix}}{{.name}}[block.id] {
    block := input.aws_ebs_volume[_]
    not checkSnapshotExists(block.id)
}

checkSnapshotExists(ebs_id) = true {
    snap := input.aws_ebs_snapshot[_]
    ebs_id == snap.config.volume_id
}