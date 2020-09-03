package utils

// IacDocument contains raw IaC file data and other metadata for a given file
type IacDocument struct {
	Type      string
	StartLine int
	EndLine   int
	FilePath  string
	Data      []byte
}
