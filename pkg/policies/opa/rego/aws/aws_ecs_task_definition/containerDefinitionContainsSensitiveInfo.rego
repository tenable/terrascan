package accurics

{{.prefix}}{{.name}}[instance.id]{
	instance := input.aws_ecs_task_definition[_]
    taskDef := instance.config.container_definitions
    taskDefJson := json_unmarshal(taskDef)
    envEntry := taskDefJson[_].environment[_]
    contains(upper(envEntry.name), upper("{{.keyword}}"))
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}