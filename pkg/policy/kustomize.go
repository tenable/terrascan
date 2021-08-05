package policy

const (
	defaultKustomizeIacType    supportedIacType    = "kustomize"
	defaultKustomizeIacVersion supportedIacVersion = version4
)

func init() {
	// Register helm as a provider with terrascan
	RegisterCloudProvider(kubernetes, defaultKustomizeIacType, defaultKustomizeIacVersion)
}
