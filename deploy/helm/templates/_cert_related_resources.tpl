{{- define "certificate_related_deployments" }}
{{- if and (eq "" .Values.secrets.tlsCertFilePath) (eq "" .Values.secrets.tlsKeyFilePath) }}
{{- $altNames := list ( printf "%s.%s.svc" .Values.name .Release.Namespace ) ( printf "%s.%s.svc.cluster.local" .Values.name .Release.Namespace ) }}
{{- $ca := genCA ( printf "%s-server-ca" "terrascan" ) 365 }}
{{- $certterrascan := genSignedCert ( printf "%s-server" "terrascan" ) nil $altNames 365 $ca }}
{{- $_ := set . "cert" $certterrascan.Cert -}}
{{- $_ := set . "key" $certterrascan.Key -}}
{{- else }}
{{- $fileCert := .Files.Get (printf "%s" .Values.secrets.tlsCertFilePath) -}}
{{- $fileKey := .Files.Get (printf "%s" .Values.secrets.tlsKeyFilePath) -}}
{{- $_ := set . "cert" $fileCert -}}
{{- $_ := set . "key" $fileKey -}}
{{- end }}
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: {{ .Values.cert_secret_name }}
  namespace: {{ .Release.Namespace }}
data:
  tls.crt: {{ .cert | b64enc }}
  tls.key: {{ .key  | b64enc }}
---
{{- if .Values.webhook.mode }}
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: {{ .Values.name }}
webhooks:
  - name: {{ .Values.webhook.name }}
    admissionReviewVersions:
    {{- range .Values.webhook.admissionReviewVersions }}
    - {{ . | printf "%s" }}
    {{ end }}
    failurePolicy: Ignore
    sideEffects: {{ .Values.webhook.sideEffects }}
    clientConfig:
      service:
        name: {{ .Values.name }}
        namespace: {{ .Release.Namespace }}
        path: {{ .Values.terrascan_webhook_key | printf "/v1/k8s/webhooks/%s/scan/validate" }}
      caBundle: {{ .cert | b64enc }}
    rules:
      - apiGroups:
        {{- range .Values.webhook.apiGroups }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        resources:
        {{- range .Values.webhook.resources }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        apiVersions:
        {{- range .Values.webhook.apiVersions }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        operations:
        {{- range .Values.webhook.operations }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
{{- end }}
---
# Had to create this file just to support validatingwebhookconfiguration failurePolicy to be FAIL.
# It turns out, webhook doesn't allow the terrascan server pod to come up in case failurePolicy is Fail.
# So, as a workaround, we create the webhook w/ Ignore, and then upgrade it to Fail in. post install chart hook. ref: https://helm.sh/docs/topics/charts_hooks/
{{- if and .Values.webhook.mode (eq .Values.webhook.failurePolicy "Fail") }}
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: {{ .Values.name }}
  annotations:
    "helm.sh/hook": "post-install"
webhooks:
  - name: {{ .Values.webhook.name }}
    admissionReviewVersions:
    {{- range .Values.webhook.admissionReviewVersions }}
    - {{ . | printf "%s" }}
    {{ end }}
    failurePolicy: Fail
    sideEffects: {{ .Values.webhook.sideEffects }}
    clientConfig:
      service:
        name: {{ .Values.name }}
        namespace: {{ .Release.Namespace }}
        path: {{ .Values.terrascan_webhook_key | printf "/v1/k8s/webhooks/%s/scan/validate" }}
      caBundle: {{ .cert | b64enc }}
    rules:
      - apiGroups:
        {{- range .Values.webhook.apiGroups }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        resources:
        {{- range .Values.webhook.resources }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        apiVersions:
        {{- range .Values.webhook.apiVersions }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
        operations:
        {{- range .Values.webhook.operations }}
        {{- if eq . ""}}
        - ""
        {{- else if eq . "*" }}
        - "*"
        {{- else }}
        - {{ . -}}
        {{- end }}
        {{- end }}
{{- end }}
{{- end -}}