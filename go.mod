module github.com/accurics/terrascan

go 1.15

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-getter v1.5.1
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform v0.14.4
	github.com/iancoleman/strcase v0.1.3
	github.com/mattn/go-isatty v0.0.12
	github.com/open-policy-agent/opa v0.22.0
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.1
	github.com/zclconf/go-cty v1.7.1
	go.uber.org/zap v1.16.0
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd
	golang.org/x/tools v0.0.0-20210115202250-e0d201561e39 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	helm.sh/helm/v3 v3.5.1
	honnef.co/go/tools v0.1.0 // indirect
	sigs.k8s.io/kustomize/api v0.7.2
)
