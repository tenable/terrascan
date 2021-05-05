package accurics

{{.prefix}}{{.name}}[workspace.id] {
    workspace := input.aws_workspaces_workspace[_]
    object.get(workspace.config, "{{.attribute_name}}", "undefined") == [false, "undefined"][_]
}