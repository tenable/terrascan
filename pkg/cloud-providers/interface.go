package cloudProvider

// CloudProvider defines the interface which every cloud provider needs to implement
// to claim support in terrascan
type CloudProvider interface {
	CreateNormalizedJson()
}
