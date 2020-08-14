package accurics

{{.prefix}}scanOnPushDisabled[retVal] {
  imageScan := input.aws_ecr_repository[_]
  imageScan.config.image_scanning_configuration == []
  traverse = "image_scanning_configuration[0].scan_on_push"
  retVal := { "Id": imageScan.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "image_scanning_configuration.scan_on_push", "AttributeDataType": "bool", "Expected": true, "Actual": imageScan.config.image_scanning_configuration[_].scan_on_push }
}

{{.prefix}}scanOnPushDisabled[retVal] {
  imageScan := input.aws_ecr_repository[_]
  some i
  imageScan.config.image_scanning_configuration[i].scan_on_push == false
  traverse := sprintf("image_scanning_configuration[%d].scan_on_push", [i])
  retVal := { "Id": imageScan.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "image_scanning_configuration.scan_on_push", "AttributeDataType": "bool", "Expected": true, "Actual": imageScan.config.image_scanning_configuration[_].scan_on_push }
}