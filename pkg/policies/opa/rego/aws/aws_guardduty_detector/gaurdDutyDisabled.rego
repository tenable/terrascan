package accurics

{{.prefix}}gaurdDutyDisabled[retVal] {
  duty := input.aws_guardduty_detector[_]
  duty.config.enable == false
  traverse = "enable"
  retVal := { "Id": duty.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "enable", "AttributeDataType": "bool", "Expected": true, "Actual": duty.config.enable }
}