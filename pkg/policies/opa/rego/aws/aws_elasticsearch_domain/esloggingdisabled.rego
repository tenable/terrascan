package accurics

{{.prefix}}esloggingdisabled[retVal] {
  esin := input.aws_elasticsearch_domain[_]
  esin.config.log_publishing_options == []
  esin.type == "aws_elasticsearch_domain"
  rc = "ewogICJsb2dfcHVibGlzaGluZ19vcHRpb25zIjogewogICAgImNsb3Vkd2F0Y2hfbG9nX2dyb3VwX2FybiI6ICI8Y2xvdWR3YXRjaF9sb2dfZ3JvdXBfYXJuPiIsCiAgICAibG9nX3R5cGUiOiAiPGxvZ190eXBlPiIsCiAgICAiZW5hYmxlZCI6IHRydWUKICB9Cn0="
  traverse = ""
  retVal := { "Id": esin.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}

{{.prefix}}esloggingdisabled[retVal] {
  esin := input.aws_elasticsearch_domain[_]
  esin.type == "aws_elasticsearch_domain"
  not esin.config.log_publishing_options
  rc = "ewogICJsb2dfcHVibGlzaGluZ19vcHRpb25zIjogewogICAgImNsb3Vkd2F0Y2hfbG9nX2dyb3VwX2FybiI6ICI8Y2xvdWR3YXRjaF9sb2dfZ3JvdXBfYXJuPiIsCiAgICAibG9nX3R5cGUiOiAiPGxvZ190eXBlPiIsCiAgICAiZW5hYmxlZCI6IHRydWUKICB9Cn0="
  traverse = ""
  retVal := { "Id": esin.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}

{{.prefix}}esloggingdisabled[retVal] {
  esin := input.aws_elasticsearch_domain[_]
  esin.type == "aws_elasticsearch_domain"
  some i
  esin.config.log_publishing_options[i].log_type != "INDEX_SLOW_LOGS"
  traverse := sprintf("log_publishing_options[%d].log_type", [i])
  retVal := { "Id": esin.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "log_publishing_options.log_type", "AttributeDataType": "string", "Expected": "INDEX_SLOW_LOGS", "Actual": esin.config.log_publishing_options[_].log_type }
}

{{.prefix}}esloggingdisabled[retVal] {
  esin := input.aws_elasticsearch_domain[_]
  some i
  esin.config.log_publishing_options[i].enabled == false
  traverse := sprintf("log_publishing_options[%d].enabled", [i])
  retVal := { "Id": esin.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "log_publishing_options.enabled", "AttributeDataType": "bool", "Expected": true, "Actual": esin.config.log_publishing_options[_].enabled }
}