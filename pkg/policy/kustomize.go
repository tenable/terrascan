package policy

const (
	defaultKustomizeIacType    supportedIacType    = "kustomize"
	defaultKustomizeIacVersion supportedIacVersion = "v3"
)

func init() {
	// Register helm as a provider with terrascan
	RegisterCloudProvider(kubernetes, defaultKustomizeIacType, defaultKustomizeIacVersion)
}
