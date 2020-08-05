package file

type FileInfo struct {
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
	Directories      []*FileInfo
	Files            []*FileInfo
}

// Metadata File metadata
type Metadata struct {
	Version string
	Groups  []*Group
}
