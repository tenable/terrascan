module github.com/accurics/terrascan

go 1.15

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/awslabs/goformation/v4 v4.19.1
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-getter v1.5.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform v0.14.4
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734
	github.com/iancoleman/strcase v0.1.3
	github.com/itchyny/gojq v0.12.1
	github.com/mattn/go-isatty v0.0.12
	github.com/mattn/go-sqlite3 v1.12.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.15.1
	github.com/onsi/gomega v1.11.0
	github.com/open-policy-agent/opa v0.22.0
	github.com/owenrumney/go-sarif v1.0.4
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.8.2
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210601080250-7ecdf8ef093b
	golang.org/x/tools v0.1.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.4.0
	honnef.co/go/tools v0.2.0 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v10.0.0+incompatible
	sigs.k8s.io/kustomize/api v0.8.5
)
