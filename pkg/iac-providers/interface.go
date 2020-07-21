package iacProvider

// IacProvider defines the interface which every IaC provider needs to implement
// to claim support in terrascan
type IacProvider interface {
	LoadIacFile(string) (interface{}, error)
}
