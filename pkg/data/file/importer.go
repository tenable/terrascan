package file

// Info File info
type Info struct {
	Path       string
	Hash       string
	HashType   string
	Attributes string
}

// Group Group metadata
type Group struct {
	Name             string
	IsReadOnly       bool
	VerifySignatures bool
	Directories      []*Info
	Files            []*Info
}

// Metadata File metadata
type Metadata struct {
	Version string
	Groups  []*Group
}
