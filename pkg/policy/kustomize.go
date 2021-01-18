package policy

const (
	defaultKustomizeIacType    supportedIacType    = "kustomize"
	defaultKustomizeIacVersion supportedIacVersion = version3
)

func init() {
	// Register helm as a provider with terrascan
	RegisterCloudProvider(kubernetes, defaultKustomizeIacType, defaultKustomizeIacVersion)
}
