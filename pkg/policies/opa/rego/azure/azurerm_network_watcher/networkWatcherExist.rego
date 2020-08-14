package accurics

{{.prefix}}networkWatcherExist[retVal] {
  not input.azurerm_network_watcher
  rc = "cmVzb3VyY2UgImF6dXJlcm1fbmV0d29ya193YXRjaGVyIiAibmV0d29ya193YXRjaGVyIiB7CiAgbmFtZSAgICAgICAgICAgICAgICA9ICJuZXR3b3JrX3dhdGNoZXIiCiAgbG9jYXRpb24gICAgICAgICAgICA9ICMjcmVzb3VyY2VfZ3JvdXBfbG9jYXRpb24jIwogIHJlc291cmNlX2dyb3VwX25hbWUgPSAjI3Jlc291cmNlX2dyb3VwX25hbWUjIwp9"
  traverse = ""
  retVal := { "Id": "network_watcher_does_not_exist", "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}

{{.prefix}}networkWatcherExist[retVal] {
  count(input.azurerm_network_watcher) <= 0
  rc = "cmVzb3VyY2UgImF6dXJlcm1fbmV0d29ya193YXRjaGVyIiAibmV0d29ya193YXRjaGVyIiB7CiAgbmFtZSAgICAgICAgICAgICAgICA9ICJuZXR3b3JrX3dhdGNoZXIiCiAgbG9jYXRpb24gICAgICAgICAgICA9ICMjcmVzb3VyY2VfZ3JvdXBfbG9jYXRpb24jIwogIHJlc291cmNlX2dyb3VwX25hbWUgPSAjI3Jlc291cmNlX2dyb3VwX25hbWUjIwp9"
  traverse = ""
  retVal := { "Id": "network_watcher_does_not_exist", "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}