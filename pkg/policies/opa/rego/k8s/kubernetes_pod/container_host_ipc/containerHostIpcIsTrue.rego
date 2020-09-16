package accurics

{{.prefix}}{{.name}}{{.suffix}}[api.id] {
    {{- template "spec" .}}
    spec.hostIPC == true
}

{{.prefix}}{{.name}}{{.suffix}}[api.id] {
    {{- template "specTF" .}}
    specTF.host_ipc == true
}


##################################
### Template definitions below ###
##################################
{{- define "api" }}
    api = input.{{.resource_type}}[_]
{{- end}}

# resolves path to the spec key
{{- define "spec" }}
    {{- template "api" . }}
    {{- if eq .resource_type "kubernetes_pod" }}
    spec = api.config.spec
    {{- else if eq .resource_type "kubernetes_pod_security_policy" }}
    spec = api.config.spec
    {{- else if eq .resource_type "kubernetes_cron_job" }}
    spec = api.config.spec.jobTemplate.spec.template.spec
    {{- else }}
    spec = api.config.spec.template.spec
    {{- end }}
{{- end }}

# resolves path to the spec key for terraform-defined k8s resources
{{- define "specTF" }}
    {{- template "api" . }}
    {{- if eq .resource_type "kubernetes_pod" }}
    specTF = api.config.spec
    {{- else if eq .resource_type "kubernetes_pod_security_policy" }}
    specTF = api.config.spec
    {{- else if eq .resource_type "kubernetes_cron_job" }}
    specTF = api.config.spec.job_template.spec.template.spec
    {{- else }}
    specTF = api.config.spec.template.spec
    {{- end }}
{{- end }}

# resolves path to the containers list
{{- define "containers" }}
    {{- template "spec" . }}
    containers = spec.containers[_]
{{- end }}

# resolves path to the containers' security context
{{- define "containersSecurityContext" }}
    {{- template "containers" . }}
    containersSecurityContext = containers.securityContext
{{- end }}

# resolves path to the containers list for terraform-defined k8s resources
{{- define "containersTF" }}
    {{- template "specTF" . }}
    containersTF = specTF.containers[_]
{{- end }}

# resolves path to the containers' security context for terraform-defined k8s resources
{{- define "containersSecurityContextTF" }}
    {{- template "containersTF" . }}
    containersSecurityContextTF = containersTF.security_context
{{- end }}

# resolves path to the initContainers list
{{- define "initContainers" }}
    {{- template "spec" . }}
    initContainers = spec.initContainers[_]
{{- end }}

# resolves path to the initContainers' security context
{{- define "initContainersSecurityContext" }}
    {{- template "initContainers" . }}
    initContainersSecurityContext = initContainers.securityContext
{{- end }}

# resolves path to the initContainers list for terraform-defined k8s resources
{{- define "initContainersTF" }}
    {{- template "specTF" . }}
    initContainersTF = specTF.init_containers[_]
{{- end }}

# resolves path to the initContainers' security context for terraform-defined k8s resources
{{- define "initContainersSecurityContextTF" }}
    {{- template "initContainersTF" . }}
    initContainersSecurityContextTF = initContainersTF.security_context
{{- end }}
